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
  "log"
)

type Worker struct {
    workers             []*labrpc.ClientEnd
    fife                *labrpc.ClientEnd //workers will also need to communicate with master fife
    kernelFunctions     map[string]KernelFunction
    tables              map[string]*Table
    me                  int
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
  log.Printf("worker %v in worker.CreateWorker", me)
  worker := &Worker{}
  worker.fife = fife
  worker.workers = workers
  worker.me = me
  return worker
}

//done with this server
func (w *Worker) Kill(){

}

//Called by RPC from fife master
//Must be called before run
func (w *Worker) Config(args *ConfigArgs, reply *ConfigReply) {
  for tableName, item := range(args.PerTableData){
    w.tables[tableName].Config(item.Partitions, item.Data)
  }
}

//we return to the master when our kernel function has finished,
//and put any reply info in RunReply.
func (w *Worker) Run(args *RunArgs, reply *RunReply) {
    // set me to this kernel instance number to use in myInstance()
    // kernelInstance = args.KernelNumber
    //TODO need to get table data and partition map from kernel before we start.
    //This happens in Config, but should we check that we're good to go? 

    // run kernel function
    w.kernelFunctions[args.KernelFunctionName](args.KernelNumber, args.KernelArgs, w.tables)

    reply.Done = true
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