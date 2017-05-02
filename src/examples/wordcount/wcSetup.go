package wordcount

import (
    "fife"
    "io/ioutil"
)

func initTables(f *fife.Fife, numPartitions int, isMaster bool) map[string]*fife.Table {
    partitioner := createHashedStringPartitioner(numPartitions)

    documents := f.CreateTable(numPartitions, fife.Accumulator{}, 
        partitioner, "documents", nil, isMaster)

    words := f.CreateTable(numPartitions,
        fife.Accumulator{
            Init: func(value interface{}) interface{} {return value},
            Accumulate: func(original interface{}, newVal interface{}) interface{} {
                return makeIntValue(getIntValue(original) + getIntValue(newVal))
                },
            },
        partitioner, "words", nil, isMaster)

    tables := make(map[string]*fife.Table)
    tables[documents.Name] = documents
    tables[words.Name] = words

    return tables
}

func StartWorker(w *fife.Worker, f *fife.Fife, numPartitions int) {
    kernelFunctions := map[string]fife.KernelFunction{"countWords":countWords}

    tables := initTables(f, numPartitions, false)

    w.Setup(kernelFunctions, tables)
}

func StartFife(f *fife.Fife, numPartitions int) {
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

    wordCount(f, fileContentsMap, numPartitions)
}