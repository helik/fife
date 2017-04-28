package fife

type Table struct {
    // Store maps partition # to data store
    //     data store maps key (string) to data (interface{})
    Store           map[int]map[string]interface{}
    // PartitionMap maps partition # to worker machine # that stores that partition
    PartitionMap    map[int]int
    Accumulator     Accumulator
    Partitioner     Partitioner

    //private state
    myWorker        Worker  //TODO might be better to put this in common.go?
    //buffer local updates to remotely stored keys. TODO not sure of the format we will want for this
    buffer          map[int]map[string]interface{}
}

//Tells us how to treat updates for a table item 
type Accumulator struct {
    Init        func(value interface{}) interface{}
    Accumulate  func(originalValue interface{}, newValue interface{}) interface{}
}

//Function mapping key to that key's data partition
type Partitioner struct {
    Which func(key string) int
}

func (t *Table) Contains(key string) bool {
    // check if key is in local partition & proceed normally
    // otherwise need to do a remote contains
    panic("unimplemented table.Contains")
    return false
}

func (t *Table) Get(key string) interface{} {
    // check if key is in local partition & proceed normally
    // otherwise need to do a remote get
    panic("unimplemented table.Get")
    return false //junk
}

func (t *Table) Put(key string, value interface{}) {
    // check if key is in local partition & proceed normally
    // otherwise need to do a remote put
}

func (t *Table) Update(key string, value interface{}) {
    // check if key is in local partition & proceed normally
    // otherwise, buffer updates
}

// flush updates on a single key to remote store
func (t *Table) Flush(key string) {
    // check if key is in local partition, do nothing
    // otherwise need to do a remote update
}

// returns partition of the table's store that is the partition # of kernelFunction
func (t *Table) GetPartition(partition int) map[string]interface{} {
    return t.Store[partition]
}

func (t *Table) isLocal(key string) bool {
    _, inLocalStore := t.Store[t.Partitioner.Which(key)]
    return inLocalStore
}
