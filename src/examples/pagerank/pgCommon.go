package pagerank

import "fife"

func initTables(numPartitions int, w *fife.Worker) map[string]*fife.Table {
    tables := make(map[string]*fife.Table)

    partitioner := fife.CreateHashedStringPartitioner(numPartitions)

    // graph maps a PageId (string) to a list of PageIds ([]string)
    graph := fife.MakeTable("graph", fife.Accumulator{}, partitioner, numPartitions, w)

    // curr maps a PageId (string) to a rank (float64)
    curr := fife.MakeTable("curr", fife.FSumAccumulator, partitioner, numPartitions, w)

    // next maps a PageId (string) to a rank (float64)
    next := fife.MakeTable("next", fife.FSumAccumulator, partitioner, numPartitions, w)

    tables[graph.Name] = graph
    tables[curr.Name] = curr
    tables[next.Name] = next

    return tables
}
