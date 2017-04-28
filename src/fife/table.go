package fife

type Table struct {
    Name            string
    //We get Store and Partition map filled by master fife when it starts this kernel
    // Store maps partition # to data store
    //     data store maps key (string) to data (interface{})
    Store           map[int]map[string]interface{}
    // PartitionMap maps partition # to worker machine # that stores that partition
    PartitionMap    map[int]int

    //private state
    myWorker        *Worker  //TODO might be better to put this in common.go?
    accumulator     Accumulator
    partitioner     Partitioner

    // updateBuffer maps key (string) to accumulated value to send in remote update
    updateBuffer    map[string]interface{}
}

//Tells us how to treat updates for a table item
type Accumulator struct {
    init        func(value interface{}) interface{}
    accumulate  func(originalValue interface{}, newValue interface{}) interface{}
}

//Function mapping key to that key's data partition
type Partitioner struct {
    which func(key string) int
}

//Return a table with initialized but empty data structures
//Intended for use on table setup.
func MakeTable(a Accumulator, p Partitioner, w *Worker, name string) *Table {
  t := &Table{}
  t.accumulator = a
  t.partitioner = p
  t.myWorker = w
  t.Name = name
  //below will be filled in when fife starts using this table
  t.Store = make(map[int]map[string]interface{})
  t.PartitionMap = make(map[int]int)
  t.updateBuffer = make(map[string]interface{})
  return t
}

func (t *Table) Contains(key string) bool {
    // check if key is in local partition & proceed normally
    localStore, inLocal := t.getLocal(key)
    if inLocal {
        _, ok := localStore[key]
        return ok
    }
    // otherwise need to do a remote contains
    var empty interface{}
    reply := t.sendRemoteTableOp(CONTAINS, key, empty)
    return reply.Contains
}

func (t *Table) Get(key string) interface{} {
    // check if key is in local partition & proceed normally
    localStore, inLocal := t.getLocal(key)
    if inLocal {
        value := localStore[key]
        return value
    }
    // otherwise need to do a remote get
    var empty interface{}
    reply := t.sendRemoteTableOp(GET, key, empty)
    return reply.Get
}

func (t *Table) Put(key string, value interface{}) {
    // check if key is in local partition & proceed normally
    localStore, inLocal := t.getLocal(key)
    if inLocal {
        localStore[key] = value
    }
    // otherwise need to do a remote put
    t.sendRemoteTableOp(PUT, key, value)
}

func (t *Table) Update(key string, value interface{}) {
    // check if key is in local partition & proceed normally
    localStore, inLocal := t.getLocal(key)
    if inLocal {
        originalValue, exists := localStore[key]
        if exists {
            localStore[key] = t.accumulator.accumulate(originalValue, value)
        } else {
            localStore[key] = t.accumulator.init(value)
        }
        return
    }
    // otherwise, buffer updates
    currentVal, inBuffer := t.updateBuffer[key]
    if inBuffer {
        t.updateBuffer[key] = t.accumulator.accumulate(currentVal, value)
    } else {
        t.updateBuffer[key] = t.accumulator.init(value)
    }
}

// flush updates on a single key to remote store
func (t *Table) Flush(key string) {
    // check if key is in buffer, if so send remote update
    val, inBuffer := t.updateBuffer[key]
    if inBuffer {
        t.sendRemoteTableOp(UPDATE, key, val)
    }
}

// returns partition of the table's store that is the partition # of kernelFunction
func (t *Table) GetPartition(partition int) map[string]interface{} {
    return t.Store[partition]
}

func (t *Table) getLocal(key string) (map[string]interface{}, bool) {
    partition := t.partitioner.which(key)
    if t.PartitionMap[partition] == myWorker().me {
        localStore := t.Store[partition]
        return localStore, true
    }
    return make(map[string]interface{}), false
}

func (t *Table) sendRemoteTableOp(op Op, key string, value interface{}) TableOpReply {
    remoteWorker := t.PartitionMap[t.partitioner.which(key)]
    return myWorker().sendRemoteTableOp(remoteWorker, t.Name, op, key, value)
}
