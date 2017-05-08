package pagerank

import "fife"

func initTables(numPartitions int, w *fife.Worker) map[string]*fife.Table {
    tables := make(map[string]*fife.Table)

    partitioner := fife.CreateHashedStringPartitioner(numPartitions)

    sum := func(v1 interface{}, v2 interface{}) interface{} {
        return v1.(float64) + v2.(float64)
    }

    fSumAccumulator := fife.CreateSumAccumulator(sum)

    // graph maps a PageId (string) to a list of PageIds ([]string)
    graph := fife.MakeTable("graph", fife.Accumulator{}, partitioner, numPartitions, w)

    // curr maps a PageId (string) to a rank (float64)
    curr := fife.MakeTable("curr", fSumAccumulator, partitioner, numPartitions, w)

    // next maps a PageId (string) to a rank (float64)
    next := fife.MakeTable("next", fSumAccumulator, partitioner, numPartitions, w)

    tables[graph.Name] = graph
    tables[curr.Name] = curr
    tables[next.Name] = next

    return tables
}