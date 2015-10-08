package main

import (
	"fmt"
//	"strings"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func startCrawlSite(url string, depth int, fetcher Fetcher, ch chan result, resMap map[string]string){
	defer close(ch)
	crawlSite(url,depth,fetcher,ch, resMap)
	
}

func crawlSite(url string, depth int, fetcher Fetcher, ch chan result, resMap map[string]string){
	
	if(depth <=0){
		return
	}
	
	if _,ok := resMap[url];ok {
		return //already crawled this url
	}
	
	ch<- result{url,""};
	
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		//fmt.Println(err)
		return
	}
	
	ch<- result{url,body};
	
	for _, u := range urls {
		crawlSite(u, depth-1, fetcher, ch, resMap)
	}

		
}

type result struct {
	url string
	body string
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	
	if depth <= 0 {
		return
	}	
	
	resultsMap := make(map[string]string) 
	
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		//fmt.Println(err)
		return
	}
	
	resultsMap[url] = body
	resultChs := make([]chan result,0,depth)
	
	
	for _, u := range urls {
		siteResCh := make(chan result)
		resultChs = append(resultChs,siteResCh)
		
		go startCrawlSite(u, depth-1, fetcher, siteResCh, resultsMap)
	}
	
	for _,resCh :=  range resultChs {
		
			for v1 := range resCh {
				resultsMap[v1.url] = v1.body
		}
	}
							   
	
	
	for k,v := range resultsMap{
		fmt.Println(k + " " + v)
	}
	return
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
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

