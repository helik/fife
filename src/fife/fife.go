package fife

import "sync"
import "labrpc"

type Fife struct {
    workers     []*labrpc.ClientEnd
    barrier     sync.WaitGroup

    tables      
}

func StartControl(workers []*labrpc.ClientEnd, tables []Table) *Fife {
    f := Fife{workers: workers}
    return &f
}

func (f *Fife) CreateTable(partitions int, accumulator Accumulator, 
    partitioner Partitioner) Table {
    
}

func (f *Fife) Run(kernelFunction KernelFunction, numPartitions int, 
    args []interface{}) {
    // assign partitions for every table
    // send config messages to workers
    // add # of kernelFunctions (I think this is partitions) to barrier.Add()
    // dispatch kernelFunctions to workers (use Run RPC)
    // when kernelFunction returns/worker is done, call barrier.Done()
}

// only makes sense to call after Run()
func (f *Fife) Barrier() {
    f.barrier.Wait()
    return
}

type DoneArgs struct {
    worker      int
}

type DoneReply struct {

}

// Done RPC Handler
func (f *Fife) Done() {
    // this needs to communicate back to the Run method (using a channel?)
    // to tell it that this worker is done
}