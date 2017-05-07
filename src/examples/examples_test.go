package examples

import (
    "testing"
    "fmt"

    "fife"
    "examples/wordcount"
    "examples/webcrawler"
    "examples/pagerank"

    "io/ioutil"
    "os"
)

func TestWordCount(t *testing.T) {
    fmt.Println("TestWordCount")

    numWorkers := 3
    cfg := fife.Make_config(t, numWorkers)

    // start workers
    for _, w := range(cfg.Workers) {
        wordcount.StartWorker(w, numWorkers)
    }

    // create test input
    fileContentsMap := make(map[string]string)
    // get which files to read
    files, err := ioutil.ReadDir("data")
    if err != nil { panic(err) }
    // read in input files
    for _, file := range files {
        fileContents, err := ioutil.ReadFile("data/"+file.Name())
        if err != nil { panic(err) }
        fileContentsMap[file.Name()] = string(fileContents)
    }

    wordcount.SetupWordCount(fileContentsMap)

    // start fife on master
    wordcount.StartFife(cfg.Fife, numWorkers)

    // check and make sure that the results were correct
    ref, err := ioutil.ReadFile("results/wc-ref.txt")
    if err != nil { panic(err) }

    actual, err := ioutil.ReadFile("results/wc.txt")
    if err != nil { panic(err) }

    if string(ref) != string(actual) {
        t.Fatalf("incorrect wc results")
    }
    fmt.Println("...passed")
}

func TestWebCrawler(t *testing.T) {
    f :=  webcrawler.RealFetcher{}
    f.Fetch("https://godoc.org/golang.org/x/net/html")
}

func TestPageRankSimple(t *testing.T) {
    fmt.Println("TestPageRankSimple")

    numWorkers := 3
    cfg := fife.Make_config(t, numWorkers)

    for _, w := range(cfg.Workers) {
        pagerank.StartWorker(w, numWorkers)
    }

    os.Remove("results/pg.txt")

    for i := range make([]int, 6) {

        pagerank.SetupPageRank("A:C,D\nB:C,D\nC:A,B,D\nD:A", i)

        pagerank.StartFife(cfg.Fife, numWorkers)
    }

    // check and make sure that the results were correct
    ref, err := ioutil.ReadFile("results/pg-ref.txt")
    if err != nil { panic(err) }

    actual, err := ioutil.ReadFile("results/pg.txt")
    if err != nil { panic(err) }

    if string(ref) != string(actual) {
        t.Fatalf("incorrect pg results")
    }
    fmt.Println("...passed")
}
