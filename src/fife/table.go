package fife

import "sync"

type Table struct {
    Name            string
    //We get Store and Partition map filled by master fife when it starts this kernel
    // Store maps partition # to data store
    //     data store maps key (string) to data (interface{})
    Store           map[int]map[string]interface{}
    // PartitionMap maps partition # to worker machine # that stores that partition
    PartitionMap    map[int]int

    //private state
    isMaster        bool
    myWorker        *Worker  //TODO might be better to put this in common.go?
    accumulator     Accumulator
    partitioner     Partitioner
    nPartitions     int //The number of partitions,
                        //and the largest partition that the partitioner is allowed to return

    rwmu              sync.RWMutex

    // updateBuffer maps key (string) to accumulated value to send in remote update
    updateBuffer    map[string]interface{}
}

//Tells us how to treat updates for a table item
type Accumulator struct {
    Init        func(value interface{}) interface{}
    Accumulate  func(originalValue interface{}, newValue interface{}) interface{}
}

//Function mapping key to that key's data partition
//TODO Currently, we break if partitioner returns something larger than nPartitions - 1
//because partitions 0 through npartitions - 1 are allocated in fife.partitionTables
//Is that ok? Should we do a safety %nPartitions whenever someone calls which()?
type Partitioner struct {
    Which func(key string) int
}

//Return a table with initialized but empty data structures
//Intended for use on table setup.
func MakeTable(a Accumulator, p Partitioner, partitions int, name string,
    isMaster bool) *Table {
  t := &Table{}
  t.accumulator = a
  t.partitioner = p
  t.Name = name
  t.nPartitions = partitions
  t.isMaster = isMaster
  //below will be filled in when fife starts using this table
  t.Store = make(map[int]map[string]interface{})
  t.PartitionMap = make(map[int]int)
  t.updateBuffer = make(map[string]interface{})
  return t
}

func (t *Table) Config(partitionMap map[int]int, store map[int]map[string]interface{}) {
    t.rwmu.Lock()
    defer t.rwmu.Unlock()
    t.PartitionMap = partitionMap
    t.Store = store
}

//initData in will be key/value pairs. We need to run partitioner on
//all input data and assign it to our store map.
// This is a lot of Puts, mimic being on master (aka everything is local)
//   --> this should really only be called on the master anyways
//        TODO: maybe restrict this to when isMaster is true?
// TODO do we actually want to call update? or just overwrite anything that's there?
func (t *Table) AddData(initData map[string]interface{}) {
    //iterate through all keys, partitioning initData
    savedIsMaster := t.isMaster
    defer func(){t.isMaster = savedIsMaster}()
    t.isMaster = true
    for key, val := range initData {
        t.Put(key, val)
    }
}

func (t *Table) Contains(key string) bool {
    // check if key is in local partition & proceed normally
    localStore, inLocal := t.getLocal(key)
    if inLocal {   
        t.rwmu.RLock()
        defer t.rwmu.RUnlock()
        _, ok := localStore[key]
        return ok
    }
    // otherwise need to do a remote contains
    reply := t.sendRemoteTableOp(CONTAINS, key, nil)
    return reply.Contains
}

func (t *Table) Get(key string) interface{} {
    // check if key is in local partition & proceed normally
    localStore, inLocal := t.getLocal(key)
    if inLocal {
        t.rwmu.RLock()
        defer t.rwmu.RUnlock()
        value := localStore[key]
        return value
    }
    // otherwise need to do a remote get
    reply := t.sendRemoteTableOp(GET, key, nil)
    return reply.Get
}

func (t *Table) Put(key string, value interface{}) {
    // check if key is in local partition & proceed normally
    localStore, inLocal := t.getLocal(key)
    if inLocal {
        t.rwmu.Lock()
        defer t.rwmu.Unlock()
        localStore[key] = value
        return
    }
    // otherwise need to do a remote put
    t.sendRemoteTableOp(PUT, key, value)
}

func (t *Table) Update(key string, value interface{}) {
    t.rwmu.Lock()
    defer t.rwmu.Unlock()
    // check if key is in local partition & proceed normally
    localStore, inLocal := t.getLocal(key)
    if inLocal {
        originalValue, exists := localStore[key]
        if exists {
            localStore[key] = t.accumulator.Accumulate(originalValue, value)
        } else {
            localStore[key] = t.accumulator.Init(value)
        }
        return
    }
    // otherwise, buffer updates
    currentVal, inBuffer := t.updateBuffer[key]
    if inBuffer {
        t.updateBuffer[key] = t.accumulator.Accumulate(currentVal, value)
    } else {
        t.updateBuffer[key] = t.accumulator.Init(value)
    }
}

// flush all buffered updates
func (t *Table) Flush() {
    t.rwmu.Lock()
    localBuffer := t.updateBuffer
    t.updateBuffer = make(map[string]interface{})
    t.rwmu.Unlock()
    for key, val := range localBuffer {
        t.sendRemoteTableOp(UPDATE, key, val)
    }
}

// returns partition of the table's store that is the partition # of kernelFunction
func (t *Table) GetPartition(partition int) map[string]interface{} {
    owner := t.PartitionMap[partition]
    if owner == t.myWorker.me {
        t.rwmu.RLock()
        defer t.rwmu.RUnlock()
        localStore, ok := t.Store[partition]
        if !ok || localStore == nil {
            t.Store[partition] = make(map[string]interface{})
        }
        return t.Store[partition]
    }
    reply := t.myWorker.sendRemoteTableOp(owner, t.Name, PARTITION,
        "", nil, partition)
    return reply.Partition
}

func (t *Table) getLocal(key string) (map[string]interface{}, bool) {
    partition := t.partitioner.Which(key) % t.nPartitions
    if t.isMaster || t.PartitionMap[partition] == t.myWorker.me {
        localStore, ok := t.Store[partition]
        if !ok || localStore == nil {
            t.Store[partition] = make(map[string]interface{})
        }
        return t.Store[partition], true
    }
    return nil, false
}

func (t *Table) sendRemoteTableOp(op Op, key string, value interface{}) TableOpReply {
    remoteWorker := t.PartitionMap[t.partitioner.Which(key) % t.nPartitions]
    return t.myWorker.sendRemoteTableOp(remoteWorker, t.Name, op, key, value, -1)
}

func (t *Table) collectData() map[string]interface{} {
    t.rwmu.RLock()
    defer t.rwmu.RUnlock()
    allData := make(map[string]interface{})
    for _, partitionStore := range t.Store {
        for k,v := range partitionStore {
            allData[k] = v
        }
    }
    return allData
}