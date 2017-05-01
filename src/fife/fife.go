package fife

import ("sync"
        "labrpc"
        "log"
 )

type Fife struct {
    rwmu        sync.RWMutex
    workers     []*labrpc.ClientEnd
    barrier     sync.WaitGroup
    tables      map[string]Table //fife master knows all data for all tables
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
  fife := &Fife{}
  fife.workers = workers
  fife.tables = make(map[string]Table)
  return fife
}

func (f *Fife) CreateTable(partitions int, accumulator Accumulator,
    partitioner Partitioner, name string, initData map[string]interface{}) *Table {
    f.rwmu.Lock()
    tbl := MakeTable(accumulator, partitioner, partitions, name)
    tbl.Store = map[int]map[string]interface{}{1:initData}
    f.tables[name] = *tbl
    f.rwmu.Unlock()
    f.partitionStore(name, *tbl)
    return tbl
}

//done with this server
func (f *Fife) Kill() {

}

//Pass the workers initial data and table partitions
//Called by control function in Run
//Assumes data has already been partitioned (partiton map has been constructed)
//Returns true if all workers successfully configed.
//TODO what args should this have?
func (f *Fife) configWorkers() bool {
  failure := false
  for workerNum, w := range(f.workers){
    args := &ConfigArgs{}
    reply := &ConfigReply{}
    args.PerTableData = make(map[string]TableData)
    //now, partition data and copy into args for this worker
    //TODO this will be very slow to do in order. can run in parallel for each worker.
    f.rwmu.RLock()
    for name, t := range(f.tables){
      data := MakeTableData()
      data.Partitions = t.PartitionMap
      for partition, keyVal := range(t.Store){
        if t.PartitionMap[partition] == workerNum{
          data.Data[partition] = keyVal
        }
      }
      args.PerTableData[name] = data
    }
    f.rwmu.RUnlock()
    ok := w.Call("Worker.Config", args, reply)
    failure = failure || !ok
    if ! ok {
      //TODO do we want to repeat failed configs, or record them in some way?
      //A failed config means some data is missing in workers
    }
  }
  return !failure
}

//Pass worker the string key to the function they should use
func (f *Fife) Run(kernelFunction string, numInstances int, //TODO should numPartitions be numInstances?
    args []interface{}) { //args is the args for the kernel function. tables passed separately
    // assign partitions for every table - means constructing the partition map for that table
    // send config messages to workers
    // wait for workers to reply successfully - if one worker is missing data, no workers can run
    // start a go thread to manage each worker
    // add # of kernelFunctions (numInstances)  to barrier.Add()
    // dispatch kernelFunctions to workers (use Run RPC)
    // when kernelFunction returns/worker is done, call barrier.Done()

    f.partitionTables()
    ok := f.configWorkers()
    if !ok{
      //TODO what should we do if some workers fail to configure?
    }

    //now, start running
    for _, w := range(f.workers){
      go func(w *labrpc.ClientEnd){

        log.Println(w)
      }(w)
    }
}

//For each table, match table partitions with workers
//Table initialization requres partition function and npartitions, so we have that
//For now, which worker gets which partition is arbitrary
func (f *Fife) partitionTables(){
  f.rwmu.Lock() //write lock
  defer f.rwmu.Unlock()
  for _,t := range(f.tables){
    for i := 0; i < t.nPartitions; i ++ {
      worker := i%len(f.workers) //TODO this may change when we colocate kernel funcs and partitions
      t.PartitionMap[i] = worker
    }
  }
}

//Data as initially read in will be key/value pairs. We need to run partitioner on
//all input data and assign it to our store map.
//Assumes nothing about the initial partition values in f.tables.Store
//Should be called once, after data and partitioner are set for all tables
func (f *Fife) partitionStore(name string, t Table){
  f.rwmu.Lock() //writing
  defer f.rwmu.Unlock()
  newStore := make(map[int]map[string]interface{})
  //give new store maps for all partitions
  for i := 0; i < t.nPartitions; i ++ {
    newStore[i] = make(map[string]interface{})
  }
  //iterate through all keys, re-partitioning
  for _, keyVal := range(t.Store){
    for key, val := range(keyVal){
      partition := t.partitioner.Which(key)
      newStore[partition][key] = val
    }
  }
  //replace table store (cannot directly assign to field in map)
  tmp := f.tables[name]
  tmp.Store = newStore
  f.tables[name] = tmp
}

// only makes sense to call after Run()
func (f *Fife) Barrier() {
    f.barrier.Wait()
    return
}
