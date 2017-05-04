package wordcount

import (
    "fife"
    "io/ioutil"
)

func initTables(numPartitions int, isMaster bool) map[string]*fife.Table {
    partitioner := createHashedStringPartitioner(numPartitions)

    documents := fife.MakeTable(fife.Accumulator{}, partitioner, numPartitions, 
        "documents", isMaster)

    words := fife.MakeTable(fife.Accumulator{
        Init: func(value interface{}) interface{} {return value},
        Accumulate: func(original interface{}, newVal interface{}) interface{} {
            return original.(int) + newVal.(int)
            },
        }, partitioner, numPartitions, "words", isMaster)

    tables := make(map[string]*fife.Table)
    tables[documents.Name] = documents
    tables[words.Name] = words

    return tables
}

func StartWorker(w *fife.Worker, numWorkers int) {
    kernelFunctions := map[string]fife.KernelFunction{"countWords":countWords}

    tables := initTables(numWorkers*2, false)

    w.Setup(kernelFunctions, tables)
}

func StartFife(f *fife.Fife, numWorkers int) {
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

    wordCount(f, fileContentsMap, numWorkers*2)
}