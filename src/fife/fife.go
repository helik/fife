package fife

import ("sync"
        "labrpc"
        "log"
 )

type Fife struct {
    workers     []*labrpc.ClientEnd
    barrier     sync.WaitGroup
//TODO could include some "ready" bool that gets switched after Setup called and completed
    //tables
}

//test code provides fife with tables
//TODO if we are providing all tables, what's the job of createtable?
func (f *Fife) Setup(tables []Table) {

}

//Config uses this to set up a fife instance on a server.
//Applicaiton using fife will later need to call Setup before fife is runnable
func CreateFife(workers []*labrpc.ClientEnd) *Fife {
  log.Printf("in fife.CreateFife")
  fife := &Fife{}
  fife.workers = workers
  return fife
}

func (f *Fife) CreateTable(partitions int, accumulator Accumulator,
    partitioner Partitioner) Table {
    return Table{}
}

//done with this server
func (f *Fife) Kill() {

}

//Pass the workers initial data and table partitions
//Called by kernel function
//Returns true if all workers successfully configed.
//TODO what args should this have?
func (f *Fife) ConfigWorkers() bool {
  args := &ConfigArgs{}
  reply := &ConfigReply{}
  failure := false
  //TODO hand follower some data
  for _, w := range(f.workers){
    ok := w.Call("Worker.Config", args, reply)
    failure = failure || !ok
    if ! ok {
      //TODO do we want to repeat failed configs, or record them in some way?
    }
  }
  return !failure
}

//Pass worker the string key to the function they should use
func (f *Fife) Run(kernelFunction string, numPartitions int, //TODO should numPartitions be numInstances?
    args []interface{}) {
    // assign partitions for every table
    // send config messages to workers
    // add # of kernelFunctions (I think this is partitions) to barrier.Add()
    //        Is num partitions always num kernels?
    // dispatch kernelFunctions to workers (use Run RPC)
    // when kernelFunction returns/worker is done, call barrier.Done()

}

// only makes sense to call after Run()
func (f *Fife) Barrier() {
    f.barrier.Wait()
    return
}
