package fife

import "labrpc"

type Worker struct {
    workers             []*labrpc.ClientEnd
    kernelFunctions     map[string]KernelFunction
    tables              []Table
}

func StartWorker(workers []*labrpc.ClientEnd, kernelFunctions map[string]KernelFunction,
    initialTables []Table) *Worker {

}

type RunArgs struct {
    Master                  *labrpc.ClientEnd
    KernelNumber            int
    KernelFunctionName      string
    KernelArgs              []interface{}
}

type RunReply struct {
}

func (w *Worker) Run(args *RunArgs, reply *RunReply) {
    // set me to this kernel instance number to use in myInstance()
    me = args.KernelNumber
    
    // run kernel function
    w.kernelFunctions[args.KernelFunctionName](args.KernelArgs, w.tables)

    // tell master we are done
    ok := false
    for !ok {
        ok = w.sendDone(args.Master)
    }
}

func (w *Worker) sendDone(master *labrpc.ClientEnd) {
    ok := master.Call("Fife.Done", args, reply)
    return ok
}

// Worker RPC calls to remote tables