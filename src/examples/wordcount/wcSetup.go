package wordcount

import (
    "fife"
)

var initialData map[string]string

func SetupWordCount(input map[string]string) {
    initialData = input
}

func StartWorker(w *fife.Worker, numWorkers int) {
    numPartitions := numWorkers*2
    kernelFunctions := map[string]fife.KernelFunction{"countWords":countWords}

    tables := initTables(numPartitions, w)

    w.Setup(kernelFunctions, tables)
}

func StartFife(f *fife.Fife, numWorkers int) {
    numPartitions := numWorkers*2

    wordCount(f, numPartitions, initialData)
}