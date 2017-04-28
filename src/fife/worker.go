package fife

import ("labrpc"
        "log"
)

type Worker struct {
    workers             []*labrpc.ClientEnd
    fife                *labrpc.ClientEnd //workers will also need to communicate with master fife
    kernelFunctions     map[string]KernelFunction
    tables              []Table
    me                  int
}

//the test code or application will provide
//kernel functions and table accumulators and partitioners here, separately from creation
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

func (w *Worker) Get(args *GetArgs, reply *GetReply) {

}

func (w *Worker) Put(args *PutArgs, reply *PutReply) {

}

func (w *Worker) Flush(args *FlushArgs, reply *FlushReply) {

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
}

//TODO: do we want to return right away telling master that we have started running?
//if so, then doneargs and donereply will be their own rpc.
//otherwise, we can just not return to the master till our kernel function has finished,
//and put any reply info in RunReply.
func (w *Worker) Run(args *RunArgs, reply *RunReply) {
    // set me to this kernel instance number to use in myInstance()
    me = args.KernelNumber

    // run kernel function
    w.kernelFunctions[args.KernelFunctionName](args.KernelArgs, w.tables)

    // tell master we are done
    dArgs := &DoneArgs{}
    dReply := &DoneReply{}
    ok := false
    for !ok {
        ok = w.sendDone(dArgs, dReply, args.Master)
    }
}

func (w *Worker) sendDone(args *DoneArgs, reply *DoneReply) bool {
    ok := fife.Call("Fife.Done", args, reply)
    return ok
}

// Worker RPC calls to remote tables
// worker is the index of the worker who has the data we want
func (w *Worker) sendPut(args *PutArgs, reply *PutReply, worker int) bool {
    ok := worker.Call("Worker.Put", args, reply)
    return ok
}

func (w *Worker) sendGet(args *GetArgs, reply *GetReply, worker int) bool {
    ok := worker.Call("Worker.Get", args, reply)
    return ok
}

func (w *Worker) sendFlush(args *FlushArgs, reply *FlushReply, worker int) bool {
    ok := worker.Call("Worker.Flush", args, reply)
    return ok
}
