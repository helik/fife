package fife

import (
    "hash/fnv"
    "math"
)

//Function mapping key to that key's data partition
//TODO Currently, we break if partitioner returns something larger than nPartitions - 1
//because partitions 0 through npartitions - 1 are allocated in fife.partitionTables
//Is that ok? Should we do a safety %nPartitions whenever someone calls which()?
type Partitioner struct {
    Which func(key string) int
}

func abs(x uint32) int {
    return int(math.Abs(float64(int(x))))
}

func CreateHashedStringPartitioner(numPartitions int) Partitioner {
    return Partitioner{func(s string) int {
        h := fnv.New32a()
        h.Write([]byte(s))
        return abs(h.Sum32()) % numPartitions
    }}
}