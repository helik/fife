package wordcount

import "fife"

func initTables(numPartitions int, w *fife.Worker) map[string]*fife.Table {
    partitioner := fife.CreateHashedStringPartitioner(numPartitions)

    documents := fife.MakeTable("documents", fife.Accumulator{}, partitioner, 
        numPartitions, w)

    sum := func(v1 interface{}, v2 interface{}) interface{} {
        return v1.(int) + v2.(int)
    }

    iSumAccumulator := fife.CreateSumAccumulator(sum)

    words := fife.MakeTable("words", iSumAccumulator, partitioner, numPartitions, w)

    tables := make(map[string]*fife.Table)
    tables[documents.Name] = documents
    tables[words.Name] = words

    return tables
}