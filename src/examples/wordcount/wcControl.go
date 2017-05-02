package wordcount

import (
    "fife"
)

// Control function
func wordCount(f *fife.Fife, files map[string]string, numPartitions int) {
    tables := initTables(f, numPartitions, true) // true for isMaster

    for k,v := range files {
        tables["documents"].Put(k, makeDocValue(v))
    }

    // f.Run(countWords, []interface{}{documents, words})

    f.Barrier()
}