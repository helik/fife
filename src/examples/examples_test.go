package examples

import (
    "testing"
    "fmt"

    "fife"
    "examples/wordcount"
)

func TestWordCount(t *testing.T) {
    fmt.Println("TestWordCount")

    numPartitions := 3

    cfg := fife.Make_config(t, numPartitions)

    // start workers
    for _, w := range(cfg.Workers) {
        wordcount.StartWorker(w, numPartitions)
    }

    // start fife on master
    wordcount.StartFife(cfg.Fife, numPartitions)

    // check output
}