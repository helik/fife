package webcrawler

import (
  "fife"
)

/*

Run go get golang.org/x/net/html to get

Used https://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html
and https://godoc.org/golang.org/x/net/html#pkg-subdirectories
as reference for parsing html for addresses,
and http://www.robotstxt.org/robotstxt.html
as a reference for handling robots.txt files
*/

const TESTING bool = true

//How many worker threads each kernel will run
const NUM_THREADS int = 2

//Hold global state for each kernel
//Shared by fetcher threads, so
type Kernel struct {
  fetch   Fetcher
  pool    chan string
  tables  map[string]*fife.Table
}

//spawned by kernel to get instances
func fetcherThread(k *Kernel){
  for { //loop forever
    url := <- k.pool
    new, err := k.fetch.Fetch(url)
    if err != nil {
      continue //go to next url in pool
    }
    for _, link := range(new){
      k.tables[URL_TABLE].Update(link, TOFETCH)
    }
    k.tables[URL_TABLE].Update(url, DONE)
    k.tables[URL_TABLE].Flush() //might have found work for other workers

  }
}

//One kernel per worker machine.
func fetcherKernel(kernelInstance int, args []interface{}, tables map[string]*fife.Table){
  //set up our kernel
  k := &Kernel{}
  k.fetch = fetcher //fake fetcher
  k.pool = make(chan string)
  k.tables = tables

  url_table := tables[URL_TABLE]

  //start worker threads
  for i := 0; i < NUM_THREADS; i ++ {
    go fetcherThread(k)
  }

  n := 0

  for { //we loop infinitely!
    //for testing, we limit to 100 runs
    n ++
    if TESTING && n > 10 {
      break
    }
    for url, _ := range url_table.GetPartition(kernelInstance) {
        if url_table.Get(url) != TOFETCH {
          continue
        }
        //TODO check domain in robots table
        //TODO check domain in politeness table
        url_table.Update(url, FETCHING)
        k.pool <- url
    }
  }

}
