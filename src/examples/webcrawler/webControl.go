package webcrawler

import (
  "fife"
  "encoding/gob"
  "sort"
  "os"
)

// Control function
func webControl(f *fife.Fife, numPartitions int, start_url string) {
    numKernels := numPartitions
    tables := initTables(numPartitions, nil)

    gob.Register(TOFETCH)  //Need to register all interface{} types that we'll send over rpc

    tables[URL_TABLE].Put(start_url, TOFETCH) //seed

    f.Setup(tables)

    f.Run(KERN, numKernels, []interface{}{},
        fife.LocalityConstriant{fife.LOCALITY_REQ, URL_TABLE})

    f.Barrier()

    //Output data. Modeled off word count.

    data := f.CollectData(URL_TABLE)

    // get all the keys
    keys := make([]string, 0)
    for k := range data {
        keys = append(keys, k)
    }

    // sort the keys
    sort.Strings(keys)

    // output data to a file in sorted-key order
    file, err := os.Create("results/web.txt")
    if err != nil { panic(err) }

    for _, k := range keys {
        _, err = file.WriteString(k + "\n")
        if err != nil { panic(err) }
    }

    err = file.Close()
    if err != nil { panic(err) }
}
