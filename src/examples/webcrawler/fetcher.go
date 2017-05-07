package webcrawler

//
// Fetcher
// Used by kernel to get URLs. Has both a fake fetcher and a real.
// Based off in-class
//

import (
  "golang.org/x/net/html"
  "net/http"
  "fmt"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (urls []string, err error)
  Robots(url string) ()
}

//Implements fetcher interface with real data
type RealFetcher struct {

}

//given the address of a page, read the page and parse its html for other web pages
func (f RealFetcher) Fetch(url string) ([]string, error) {

  resp, err := http.Get(url)
  urls := []string{}

  tokens := html.NewTokenizer(resp.Body)
  for item := tokens.Next(); item != html.ErrorToken; item = tokens.Next() {
    if item == html.StartTagToken{
      token := tokens.Token()
      if token.Data == "a" {//tag for a url
        for _, attribute := range token.Attr { //iterate over attributes till we find url
          if attribute.Key == "href" {
            urls = append(urls, attribute.Val)
            //TODO when we test with https://godoc.org/golang.org/x/net/html,
            //we can see that some are not real links; eg /builtin#byte
            break
          }
        }
      }
    }
  }
  //TODO how to deal with err?
  fmt.Println(urls)
  return urls, err
}

//Return a list of dissalowed addresses in this domain
func checkRobots(address string) []string {
  panic("check robots not implemented")
  return []string{""}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) ([]string, error) {
	if res, ok := f[url]; ok {
		return res.urls, nil
	}
	return nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
