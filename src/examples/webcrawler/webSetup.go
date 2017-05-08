package webcrawler

import (
  "fife"
)

func StartWorker(w *fife.Worker, numWorkers int) {
    numPartitions := numWorkers //we will have one kernel per worker, and one partition per kernel
    kernelFunctions := map[string]fife.KernelFunction{KERN:fetcherKernel}

    tables := initTables(numPartitions, w)

    w.Setup(kernelFunctions, tables)
}

func StartFife(f *fife.Fife, numWorkers int) {
    numPartitions := numWorkers
    tables := initTables(numPartitions, nil) //TODO these two lines are in the kernel
                                             //function in the wc example. Should we standardize that?
    f.Setup(tables)

    webControl(f, numPartitions)
}
