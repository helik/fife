package webcrawler

/*

Run go get golang.org/x/net/html to get

Used https://schier.co/blog/2015/04/26/a-simple-web-scraper-in-go.html
and https://godoc.org/golang.org/x/net/html#pkg-subdirectories
as reference for parsing html for addresses,
and http://www.robotstxt.org/robotstxt.html
as a reference for handling robots.txt files 
*/

import (
  "golang.org/x/net/html"
  "net/http"
  "fmt"
)

//given the address of a page, read the page and parse its html for other web pages
func ReadPage(url string) []string {
  resp, _ := http.Get(url)

  tokens := html.NewTokenizer(resp.Body)
  for item := tokens.Next(); item != html.ErrorToken; item = tokens.Next() {
    if item == html.StartTagToken{
      token := tokens.Token()
      if token.Data == "a" {//tag for a url
        fmt.Printf("token %v, link: ",token)
        for _, attribute := range token.Attr { //iterate over attributes till we find url
          if attribute.Key == "href" {
            fmt.Printf("%v\n",attribute.Val )
            //TODO when we test with https://godoc.org/golang.org/x/net/html,
            //we can see that some are not real links; eg /builtin#byte
            break
          }
        }
      }
    }
  }
  //TODO how to deal with err?
  return []string{""}
}

//Return a list of dissalowed addresses in this domain
func checkRobots(address string) []string {
  panic("check robots not implemented")
  return []string{""}
}
