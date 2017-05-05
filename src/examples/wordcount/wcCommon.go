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

func initTables(numPartitions int, w *fife.Worker) map[string]*fife.Table {
    partitioner := createHashedStringPartitioner(numPartitions)

    documents := fife.MakeTable("documents", fife.Accumulator{}, partitioner, 
        numPartitions, w)

    words := fife.MakeTable("words", fife.Accumulator{
        Init: func(value interface{}) interface{} {return value},
        Accumulate: func(original interface{}, newVal interface{}) interface{} {
            return original.(int) + newVal.(int)
            },
        }, partitioner, numPartitions, w)

    tables := make(map[string]*fife.Table)
    tables[documents.Name] = documents
    tables[words.Name] = words

    return tables
}