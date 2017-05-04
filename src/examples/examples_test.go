package examples

import (
    "testing"
    "fmt"
    "fife"
    "examples/wordcount"
)

func TestWordCount(t *testing.T) {
    fmt.Println("TestWordCount")

    numWorkers := 3

    cfg := fife.Make_config(t, numWorkers)

    // start workers
    for _, w := range(cfg.Workers) {
        wordcount.StartWorker(w, numWorkers)
    }

    // // start fife on master
    wordcount.StartFife(cfg.Fife, numWorkers)
}
