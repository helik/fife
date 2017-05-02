package wordcount

import (
    "fife"
    "hash/fnv"
    "math"
)

type IntValue struct {
    value   int
}

func makeIntValue(val int) IntValue {
    return IntValue{val}
}

func getIntValue(valObj interface{}) int {
    return valObj.(IntValue).value
}

type DocValue struct {
    value   string
}

func makeDocValue(val string) DocValue {
    return DocValue{val}
}

func getDocValue(valObj interface{}) string {
    return valObj.(DocValue).value
}

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