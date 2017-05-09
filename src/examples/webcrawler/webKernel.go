package webcrawler

import (
  "fife"
  "log"
)

/*

Run go get golang.org/x/net/html to get

Used https://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html
and https://godoc.org/golang.org/x/net/html#pkg-subdirectories
as reference for parsing html for addresses,
and http://www.robotstxt.org/robotstxt.html
as a reference for handling robots.txt files
*/

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
      log.Printf("err\n")
      continue //go to next url in pool TODO should we not try this one again?
    }
    for _, link := range(new){
      k.tables[URL_TABLE].Update(link, TOFETCH)
    }
    k.tables[URL_TABLE].Update(url, DONE)
    k.tables[URL_TABLE].Flush() //might have found work for other workers

    log.Printf(url)
  }
}

//One kernel per worker machine.
func fetcherKernel(kernelInstance int, args []interface{}, tables map[string]*fife.Table){
  //set up our kernel
  k := &Kernel{}
  k.fetch = fetcher //fake fetcher
  k.pool = make(chan string)
  k.tables = tables

//  politeness := tables[POLITENESS]
  url_table := tables[URL_TABLE]
//  robots := tables[ROBOTS]

  //start worker threads
  for i := 0; i < NUM_THREADS; i ++ {
    go fetcherThread(k)
  }

  for { //we loop infinitely! TODO this will be hella sad irl.
    //should limit for testing purposes
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


  url_table.Flush() //TODO where should we flush? inside for loop?
}
