package pagerank

import "fife"

func pageRank(f *fife.Fife, numWorkers int) {
    propagationFactor := 0.85

    numPartitions := numWorkers*3

    tables := initTables(numPartitions, nil)

    // TODO: add initial graph table data

    f.Setup(tables)

    for i := 0; i < 50; i++ {
        f.Run("pgKernel", numPartitions, []interface{}{propagationFactor}, 
            fife.LocalityConstriant{fife.LOCALITY_REQ,"graph"})

        f.Barrier()

        nextData := f.CollectData("next")

        tables = initTables(numPartitions, nil)

        tables["curr"].AddData(nextData)

        f.Setup(tables)
    }

}