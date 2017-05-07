package pagerank

import "fife"

var initialData string
var numIterations int

func SetupPageRank(input string, iterations int) {
    initialData = input
    numIterations = iterations
}

func StartWorker(w *fife.Worker, numWorkers int) {
    kernelFunctions := map[string]fife.KernelFunction{"pgKernel":pgKernel}

    numPartitions := numWorkers
    tables := initTables(numPartitions, w)

    w.Setup(kernelFunctions, tables)
}

func StartFife(f *fife.Fife, numWorkers int) {
    numPartitions := numWorkers
    pageRank(f, numPartitions, initialData, numIterations)
}