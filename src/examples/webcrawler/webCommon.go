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

    url_table := fife.MakeTable("url_table",
        fife.CreateMaxAccumulator(func(a interface{}, b interface{}) bool{
		      return a.(WebState) < b.(WebState)
	        }),
        partitioner,
        numPartitions, w)

    politeness := fife.MakeTable("politeness",
      fife.CreateMaxAccumulator(func(a interface{}, b interface{}) bool{
        return a.(int) > b.(int)
        }),
      partitioner, numPartitions, w)

    robots := fife.MakeTable("robots", fife.FirstAccumulator, partitioner,
      numPartitions, w)

    tables := make(map[string]*fife.Table)
    tables[url_table.Name] = url_table
    tables[politeness.Name] = politeness
    tables[robots.Name] = robots

    return tables
}
