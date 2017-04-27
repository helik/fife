package fife

import ("labrpc"
        "log"
)

type Worker struct {
    workers             []*labrpc.ClientEnd
    fife                *labrpc.ClientEnd //workers will also need to communicate with master fife
    kernelFunctions     map[string]KernelFunction
    tables              map[string]Table
    me                  int
}

//the test code or application will provide
//kernel functions and tables here, separately from creation
func (w *Worker) Setup(kernelFunctions map[string]KernelFunction,
    initialTables []Table) {

}

//playing around with this alternative that we can call from config without doing all the setup in config...
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

type RunArgs struct {
    Master                  *labrpc.ClientEnd
    KernelNumber            int
    KernelFunctionName      string
    KernelArgs              []interface{}
}

type RunReply struct {
    Done    bool
}


//TODO: do we want to return right away telling master that we have started running?
//if so, then doneargs and donereply will be their own rpc.
//otherwise, we can just not return to the master till our kernel function has finished,
//and put any reply info in RunReply.
func (w *Worker) Run(args *RunArgs, reply *RunReply) {
    // set me to this kernel instance number to use in myInstance()
    kernelInstance = args.KernelNumber

    // run kernel function
    w.kernelFunctions[args.KernelFunctionName](args.KernelArgs, w.tables)

    reply.Done = true
}

// Worker RPC calls to remote tables

type TableOpArgs struct {
    TableName   string
    Op          Op
    Key         string
    Value       interface{}
}

type TableOpReply struct {
    Done        bool
    Contains    bool
    Get         interface{}
}

func (w *Worker) sendRemoteTableOp(worker int, tableName string, operation Op, 
    key string, value interface{}) TableOpReply {
    args := TableOpArgs{
        TableName: tableName,
        Op: operation,
        Key: key,
        Value: value,
    }
    var reply TableOpReply
    ok := w.workers[worker].Call("Worker.TableOpRPC", args, reply)
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
    }

    reply.Done = true
}