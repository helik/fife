package wordcount

import (
    "fife"
    "io/ioutil"
)

func StartWorker(w *fife.Worker, numWorkers int) {
    numPartitions := numWorkers*2
    kernelFunctions := map[string]fife.KernelFunction{"countWords":countWords}

    tables := initTables(numPartitions, w)

    w.Setup(kernelFunctions, tables)
}

func StartFife(f *fife.Fife, numWorkers int) {
    numPartitions := numWorkers*2
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