package webcrawler

import (
  "fife"
)

//possible statuses for web
type WebState int

const (
	DONE WebState = iota
	BLACKLISTED
	FETCHING
  TOFETCH
)

func initTables(numPartitions int, w *fife.Worker) map[string]*fife.Table {
    partitioner := fife.CreateHashedStringPartitioner(numPartitions)
    accumulator := fife.CreateMaxAccumulator(func(a interface{}, b interface{}) bool{
		     return a.(WebState) < b.(WebState)
	  })

    url_table := fife.MakeTable("url_table", accumulator, partitioner,
        numPartitions, w)

    politeness := fife.MakeTable("politeness", fife.Accumulator{
        }, partitioner, numPartitions, w)



    tables := make(map[string]*fife.Table)
    tables[url_table.Name] = url_table
    tables[politeness.Name] = politeness

    return tables
}
