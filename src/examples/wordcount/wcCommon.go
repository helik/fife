package wordcount

import "fife"

func initTables(numPartitions int, w *fife.Worker) map[string]*fife.Table {
    partitioner := fife.CreateHashedStringPartitioner(numPartitions)

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