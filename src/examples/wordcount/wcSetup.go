package wordcount

import (
    "fife"
    "io/ioutil"
)

func partition_simple(key string) int{
  switch key[0]{
  case 'a':
    return 0
  case 'b':
    return 1
  default:
    return 2
  }
}

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

    simple := fife.MakeTable(fife.Accumulator{}, fife.Partitioner{partition_simple}, numPartitions, "simple", isMaster)
    tables[simple.Name] = simple

    return tables
}

func StartWorker(w *fife.Worker, numPartitions int) {
    kernelFunctions := map[string]fife.KernelFunction{"countWords":CountWords}

    tables := initTables(numPartitions, false)

    w.Setup(kernelFunctions, tables)
}

func StartFife(f *fife.Fife, numPartitions int) {
    // create test input
    fileContentsMap := make(map[string]string)
    // get which files to read
    files, err := ioutil.ReadDir("smalldata")
    if err != nil { panic(err) }
    // read in input files
    for _, file := range files {
        fileContents, err := ioutil.ReadFile("smalldata/"+file.Name())
        if err != nil { panic(err) }
        fileContentsMap[file.Name()] = string(fileContents)
    }

    wordCount(f, fileContentsMap, numPartitions)
}