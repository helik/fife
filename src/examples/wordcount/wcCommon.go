package wordcount

import (
    "fife"
    "hash/fnv"
    "math"
)

func abs(x uint32) int {
    return int(math.Abs(float64(int(x))))
}

func createHashedStringPartitioner(numPartitions int) fife.Partitioner {
    return fife.Partitioner{func(s string) int {
        h := fnv.New32a()
        h.Write([]byte(s))
        return abs(h.Sum32()) % numPartitions
    }}
}