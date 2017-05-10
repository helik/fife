package fife

/*
Worker set-up:
1) CreateWorker is called by config, to create server and connect it to its peers and master
2) Setup is called by application on each server, to initialize tables and deliver kernel func
3) Config is called via RPC by fife master, to deliver data and table partition info
4) Run is called by master to tell worker to start kernel function
*/

import (
  "labrpc"
  "sync"
)

type Worker struct {
    workers             []*labrpc.ClientEnd
    fife                *labrpc.ClientEnd //workers will also need to communicate with master fife
    kernelFunctions     map[string]KernelFunction
    tables              map[string]*Table
    me                  int

    mu                  sync.Mutex

    // map update num (int64) to a partition update struct
    partitionUpdateTable        map[int64]*PartitionUpdate

    // map partition num (int) to chan -- indicate on chan once this worker owns the partition num
    waitingRemoteTableOps       map[int]chan bool

    killChan        chan bool
}

//the test code or application will provide
//kernel functions and table accumulators and partitioners here, separately from creation
//Tables here only need Accumulators
func (w *Worker) Setup(kernelFunctions map[string]KernelFunction,
    initialTables map[string]*Table) {
    w.kernelFunctions = kernelFunctions
    w.tables = initialTables
}

//Called by the config file to create a worker server
func CreateWorker(fife *labrpc.ClientEnd, workers []*labrpc.ClientEnd, me int) *Worker {
  worker := &Worker{}
  worker.fife = fife
  worker.workers = workers
  worker.me = me
  worker.partitionUpdateTable = make(map[int64]*PartitionUpdate)
  return worker
}

//done with this server
func (w *Worker) Kill(){
    close(w.killChan)
}

//Called by RPC from fife master
//Must be called before run
func (w *Worker) Config(args *ConfigArgs, reply *Reply) {
  // if args.PerTableData == nil { //no data was passed to us to configure
  //   return
  // }
  for tableName, item := range(args.PerTableData){
    w.tables[tableName].config(item.Partitions, item.Data)
  }
}

//we return to the master when our kernel function has finished,
//and put any reply info in RunReply.
func (w *Worker) Run(args *RunArgs, reply *Reply) {
    // set me to this kernel instance number to use in myInstance()
    // kernelInstance = args.KernelNumber
    //TODO need to get table data and partition map from kernel before we start.
    //This happens in Config, but should we check that we're good to go?

    // run kernel function
    w.kernelFunctions[args.KernelFunctionName](args.KernelNumber, args.KernelArgs, w.tables)

    reply.Done = true
}

// Partition Update

type PartitionUpdate struct {
    ackedWorkers    []int
    partitionNum    int
    // map of table name (string) to store data map[string]interface{}
    partitionData   map[string]map[string]interface{}
}

func (w *Worker) PartitionUpdate(args *PartitionUpdateArgs, reply *Reply) {
    // if we are the new worker, we don't want to update our table until we are
    //   ready to switch -- add ourself to the updateAckTable & return
    if args.NewWorker == w.me {
        update, inTable := w.partitionUpdateTable[args.UpdateNum]
        if inTable {
            workers := update.ackedWorkers
            for _, worker := range workers {
                if worker == w.me { return }
            }
            update.partitionNum = args.PartitionNum
            update.ackedWorkers = append(workers, w.me)
        } else {
            w.partitionUpdateTable[args.UpdateNum] = &PartitionUpdate {
                ackedWorkers: []int{w.me},
                partitionNum: args.PartitionNum,
            }
        }
        return
    }

    // update partition maps
    for _, table := range w.tables {
        table.partitionUpdate(args.PartitionNum, args.NewWorker, nil)
    }

    var partitionData map[string]map[string]interface{}

    // if we are the old worker, send along the data
    if args.OldWorker == w.me {
        partitionData = make(map[string]map[string]interface{})
        for tableName, table := range w.tables {
            partitionData[tableName] = table.getPartitionAndDelete(args.PartitionNum)
        }
    }

    // tell the new worker we are ready to switch
    go w.sendPartitionUpdateAck(args.NewWorker, args.UpdateNum, args.PartitionNum,
        partitionData)
}

type PartitionUpdateAckArgs struct {
    UpdateNum       int64
    PartitionNum    int
    WorkerNum       int
    // map of table name (string) to store data map[string]interface{}
    PartitionData   map[string]map[string]interface{}
}

func (w *Worker) sendPartitionUpdateAck(newWorker int, updateNum int64, partitionNum int,
    partitionData map[string]map[string]interface{}) {
    args := PartitionUpdateAckArgs{updateNum, partitionNum, w.me, partitionData}
    w.workers[newWorker].Call("Worker.PartitionUpdateAck", &args, nil)
}

// RPC handler for partition update acks from other workers
// this worker must be the target for the partition move & should wait until everyone
// has responded before switching to owning the partition
func (w * Worker) PartitionUpdateAck(args *PartitionUpdateAckArgs, reply *Reply) {
    w.mu.Lock()
    defer w.mu.Unlock()

    update, inTable := w.partitionUpdateTable[args.UpdateNum]
    // make sure this worker is not already in the list (to avoid duplicates)
    var workers []int
    if inTable {
        workers = update.ackedWorkers
        for _, w := range workers {
            if w == args.WorkerNum { return }
        }
    } else {
        w.partitionUpdateTable[args.UpdateNum] = &PartitionUpdate{}
    }
    // if not duplicate, add to table
    w.partitionUpdateTable[args.UpdateNum].ackedWorkers = append(workers, args.WorkerNum)

    // if old worker, save the partition data
    if args.PartitionData != nil {
        w.partitionUpdateTable[args.UpdateNum].partitionData = args.PartitionData
    }

    // check to see if we have all of the updates
    //   (the +1 is for the new worker we just added)
    if len(workers) + 1 == len(w.workers) {
        // spawn a go routine to not block the caller
        go func(){
            // add new partition to tables
            for tableName, table := range w.tables {
                table.partitionUpdate(args.PartitionNum, w.me,
                    w.partitionUpdateTable[args.UpdateNum].partitionData[tableName])
            }
            // unblock remote ops
            w.waitingRemoteTableOps[args.PartitionNum] = make(chan bool)
            close(w.waitingRemoteTableOps[args.PartitionNum])
        }()
    }
}

// Worker RPC calls to remote tables

type TableOpArgs struct {
    TableName   string
    Op          Op
    Key         string
    Value       interface{}
    Partition   int
}

type TableOpReply struct {
    Done        bool
    Contains    bool
    Get         interface{}
    Partition   map[string]interface{}
}

func (w *Worker) sendRemoteTableOp(worker int, tableName string, operation Op,
    key string, value interface{}, partition int) TableOpReply {
    args := TableOpArgs{
        TableName: tableName,
        Op: operation,
        Key: key,
        Value: value,
        Partition: partition,
    }
    var reply TableOpReply
    ok := w.workers[worker].Call("Worker.TableOpRPC", &args, &reply)
    if !ok || !reply.Done {
        // TODO: retry if the rpc failed?
    }
    return reply
}

func (w *Worker) TableOpRPC(args *TableOpArgs, reply *TableOpReply) {
    targetTable := w.tables[args.TableName]

    // if this worker doesn't have the partition asked for, assume we need to block
    if !targetTable.isLocalPartition(args.Partition) {
        // wait until worker gets this partition
        _, exists := w.waitingRemoteTableOps[args.Partition]
        if !exists {
            w.waitingRemoteTableOps[args.Partition] = make(chan bool)
        }
        select {
        case <- w.killChan:
            return
        case <- w.waitingRemoteTableOps[args.Partition]:
        }
    }

    switch args.Op {
    case CONTAINS:
        reply.Contains = targetTable.Contains(args.Key)
    case GET:
        reply.Get = targetTable.Get(args.Key)
    case PUT:
        targetTable.Put(args.Key, args.Value)
    case UPDATE:
        targetTable.Update(args.Key, args.Value)
    case PARTITION:
        reply.Partition = targetTable.GetPartition(args.Partition)
    }

    reply.Done = true
}

func (w *Worker) CollectData(args *CollectDataArgs, reply *CollectDataReply) {
    reply.TableData = w.tables[args.TableName].collectData()
}
