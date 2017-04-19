package fife

type Table struct {
    // Store maps partition # to data store
    //     data store maps key (string) to data (interface{})
    Store           map[int]map[string]interface{}
    // PartitionMap maps partition # to worker machine # that stores that partition
    PartitionMap    map[int]int
    Accumulator     Accumulator
    Partitioner     Partitioner
}

type Accumulator struct {
    Init        func(value interface{}) interface{}
    Accumulate  func(originalValue interface{}, newValue interface{}) interface{}
}

type Partitioner struct {
    Which func(key string) int
}

func (t *Table) Contains(key string) bool {
    // check if key is in local partition & proceed normally
    // otherwise need to do a remote contains
}

func (t *Table) Get(key string) interface{} {
    // check if key is in local partition & proceed normally
    // otherwise need to do a remote get
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

func (t *Table) isLocal(key string) {
    _, inLocalStore := t.Store[t.Partitioner.Which(key)]
    return inLocalStore
}