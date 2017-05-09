package webcrawler

import (
  "fife"
)

//Kernel function names
const KERN string = "fetcherKernel"

//Table names
const POLITENESS string = "politeness"
const URL_TABLE string = "url_table"
const ROBOTS string = "robots"

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

    url_table := fife.MakeTable(URL_TABLE,
        fife.CreateMaxAccumulator(func(a interface{}, b interface{}) bool{
		      return a.(WebState) < b.(WebState)
	        }),
        partitioner,
        numPartitions, w)

    //politeness holds the time
    politeness := fife.MakeTable(POLITENESS,
      fife.CreateMaxAccumulator(func(a interface{}, b interface{}) bool{
        return a.(int) > b.(int)
        }),
      partitioner, numPartitions, w)

    robots := fife.MakeTable(ROBOTS, fife.FirstAccumulator, partitioner,
      numPartitions, w)

    tables := make(map[string]*fife.Table)
    tables[url_table.Name] = url_table
    tables[politeness.Name] = politeness
    tables[robots.Name] = robots

    return tables
}
