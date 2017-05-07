package fife

import ("sync"
        "labrpc"
 )

type Fife struct {
    rwmu        sync.RWMutex
    workers     []*labrpc.ClientEnd
    barrier     sync.WaitGroup
    tables      map[string]*Table //fife master knows all data for all tables
//TODO could include some "ready" bool that gets switched after Setup called and completed
    //tables
}

//test code provides fife with tables
func (f *Fife) Setup(tables map[string]*Table) {
  f.tables = tables
}

//Config uses this to set up a fife instance on a server.
//Applicaiton using fife will later need to call Setup before fife is runnable
func CreateFife(workers []*labrpc.ClientEnd) *Fife {
  fife := &Fife{}
  fife.workers = workers
  fife.tables = make(map[string]*Table)
  return fife
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
    args []interface{}, loc LocalityConstriant) { //args is the args for the kernel function. tables passed separately
    // assign partitions for every table - means constructing the partition map for that table
    // send config messages to workers
    // wait for workers to reply successfully - if one worker is missing data, no workers can run
    // start a go thread to manage each worker
    // add # of kernelFunctions (numInstances)  to barrier.Add()
    // dispatch kernelFunctions to workers (use Run RPC)
    // when kernelFunction returns/worker is done, call barrier.Done()

    f.partitionTables() //partition tables among workers.
    ok := f.configWorkers() //send all tables to workers
    if !ok{
      //TODO what should we do if some workers fail to configure?
      //presumably, re-partition that data and send it to the remaining workers
    }
    //now, start running. we have some num workers and numInstances to run
    f.barrier.Add(numInstances)
    allInstances := make(chan int, numInstances) //Used in no locality situation
    for i := 0; i < numInstances; i ++ {
      allInstances <- i
    }
    for w := range(f.workers){
      go func(w int){
        worker := f.workers[w]
        rArgs := &RunArgs{}
        reply := &RunReply{}
        rArgs.KernelFunctionName = kernelFunction
        rArgs.KernelArgs = args
        if loc.Loc == LOCALITY_REQ {
          partitions := f.tables[loc.Table].PartitionMap
          myInstances := []int{}
          for partition, worker := range(partitions){
            if worker == w{
              myInstances = append(myInstances, partition)
            }
          }
          //now, run all instances
          for _, instance := range(myInstances){
            rArgs.KernelNumber = instance
            ok := worker.Call("Worker.Run", rArgs, reply)
            if ok {
              f.barrier.Done()
            }else {
              panic("failed worker.Run not impelemnted in fife run")
            }
          }
        }else{ //we can run whichever instances need running
          for {
            instance, more := <- allInstances
            if more { //something to run
              rArgs.KernelNumber = instance
              ok := worker.Call("Worker.Run", rArgs, reply)
              if ok {
                f.barrier.Done() //TODO do we want any other checks?
              }
            } else { //TODO this never happens right now, because chan is never closed
              //could close from Barrier, but we'd need a reference to it
              return //chan has been closed, nothing else to run
            }
          }
        }
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

// only makes sense to call after Run()
// However, it will work no matter when you call it since if you call it before Run()
//     it will automatically return since there is nothing in the wait group
func (f *Fife) Barrier() {
    f.barrier.Wait()
    return
}

func (f *Fife) CollectData(tableName string) map[string]interface{} {
  allData := make(map[string]interface{})
  for _, w := range f.workers {
    args := CollectDataArgs{tableName}
    var reply CollectDataReply

    w.Call("Worker.CollectData", &args, &reply)

    for k, v := range reply.TableData {
      allData[k] = v
    }
  }
  return allData
}
