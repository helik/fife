package wordcount

import (
    "fife"
    "sort"
    "os"
    "strconv"
)

// Control function
func wordCount(f *fife.Fife, files map[string]string, numPartitions int) {
    tables := initTables(numPartitions, true) // true for isMaster

    for k,v := range files {
        tables["documents"].Put(k, v)
    }

    f.Setup(tables)

    var args interface{}
    f.Run("countWords", numPartitions, []interface{}{args}, fife.LocalityConstriant{fife.NONE,""})

    f.Barrier()

    data := f.CollectData("words")

    // get all the keys
    keys := make([]string, 0)
    for k := range data {
        keys = append(keys, k)
    }

    // sort the keys
    sort.Strings(keys)

    // output data to a file in sorted-key order
    file, err := os.Create("results/wc.txt")
    if err != nil { panic(err) }

    for _, k := range keys {
        s := k + ":" + strconv.Itoa(int(data[k].(int))) + "\n"
        _, err = file.WriteString(s)
        if err != nil { panic(err) }
    }

    err = file.Close()
    if err != nil { panic(err) }
}
