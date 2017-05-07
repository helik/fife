package examples

import (
    "testing"
    "fmt"

    "fife"
    "examples/wordcount"
    "examples/webcrawler"
    "examples/pagerank"

    "io/ioutil"
)

func TestWordCount(t *testing.T) {
    fmt.Println("TestWordCount")

    numWorkers := 3
    cfg := fife.Make_config(t, numWorkers)

    // start workers
    for _, w := range(cfg.Workers) {
        wordcount.StartWorker(w, numWorkers)
    }

    // start fife on master
    wordcount.StartFife(cfg.Fife, numWorkers)

    // check and make sure that the results were correct
    ref, err := ioutil.ReadFile("results/ref-wc.txt")
    if err != nil { panic(err) }

    actual, err := ioutil.ReadFile("results/wc.txt")
    if err != nil { panic(err) }

    if string(ref) != string(actual) {
        t.Fatalf("incorrect wc results")
    }
    fmt.Println("...passed")
}

func TestWebCrawler(t *testing.T) {
    webcrawler.ReadPage("https://godoc.org/golang.org/x/net/html")
}

func TestPageRank(t *testing.T) {
    fmt.Println("TestPageRank")

    numWorkers := 3
    cfg := fife.Make_config(t, numWorkers)

    for _, w := range(cfg.Workers) {
        pagerank.StartWorker(w, numWorkers)
    }

    pagerank.StartFife(cfg.Fife, numWorkers)

    // TODO: check output
}
