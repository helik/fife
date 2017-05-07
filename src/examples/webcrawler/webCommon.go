package webcrawler

import (
  "fife"
  "hash/fnv"
  "math"
)

//Same partitioner as word count
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

//The accumulator for url_table takes the max of WebState
func createMaxAccumulator() fife.Accumulator {
  return fife.Accumulator{
    Init: func(value interface{}) interface{} {return value},
    Accumulate: func (initialValue interface{}, newValue interface{}) interface{} {
      if initialValue.(WebState) < newValue.(WebState) {
        return initialValue
      }
      return newValue
    },
  }
}


//possible statuses for web
type WebState int

const (
	DONE WebState = iota
	BLACKLISTED
	FETCHING
  TOFETCH
)

func initTables(numPartitions int, w *fife.Worker) map[string]*fife.Table {
    partitioner := createHashedStringPartitioner(numPartitions)
    accumulator := createMaxAccumulator()

    url_table := fife.MakeTable("url_table", accumulator, partitioner,
        numPartitions, w)

    politeness := fife.MakeTable("politeness", fife.Accumulator{
        }, partitioner, numPartitions, w)



    tables := make(map[string]*fife.Table)
    tables[url_table.Name] = url_table
    tables[politeness.Name] = politeness

    return tables
}
