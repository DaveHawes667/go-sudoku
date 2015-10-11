package main

import (
	"fmt"
)

//General solver error
type SolveError struct{
	m_info string
}

func (e SolveError) Error() string {
	return e.m_info
}

//General solver interface
type Solver interface {
	Solved() (bool,*SolveError)
}

//Cells which store individual numbers in the grid
type cell struct {
	m_possible map[int]bool 
}

func (c cell) SetKnownTo(value int){
	for k,_ := range c.m_possible{
		if k != value {
			delete(c.m_possible,k)
		}
	} 
}

func (c cell) TakeKnownFromPossible(known []int){
	for _,v := range known{
		delete(c.m_possible,v)
	}
}

func (c cell) Known() (int, *SolveError){
	
	// by convention we delete from the map possibles that are no longer possible
	// so we just need to check map length to see if the cell is solved
	if len(c.m_possible)!=1{
		return 0,&SolveError{"Value not yet known for this cell"}
	}
	
	//Only one key is now considered "possible", it's value should be true, and it should
	//be the only one in the list, return it if that is the case 
	for k,v := range c.m_possible {
		if v{
			return k,nil
		}
	}
	
	return 0,&SolveError{"Error in cell storage of known values"}
}

type cellPtrSlice []*cell 

func (cells cellPtrSlice) Solved() (bool, *SolveError){
	for _,c := range cells{
		_,err := c.Known();
		if err != nil{
			return false,err
		}
	}
	
	return true,nil
}

//Squares which represent one of each of the 9 squares in a grid, each of which 
//references a 3x3 collection of cells.

type square struct {
	m_cells [][]*cell
}

func (s square) Solved() (bool,*SolveError) {
	for _,r := range s.m_cells{
		solved,err := cellPtrSlice(r).Solved()
		if(err != nil){
			return false,err
		}
		if !solved {
			return false,nil
		}
		
	}
	return true,nil
}
func (s square) init() {
	s.m_cells = make([][]*cell,SQUARE_SIZE)
	for i,_ := range s.m_cells{
		s.m_cells[i] = make(cellPtrSlice, SQUARE_SIZE)
	}
}

//A horizontal or vertical line of 9 cells through the entire grid.
type line struct {
	m_cells cellPtrSlice
}

func (l line) Solved() (bool,*SolveError) {
	return l.m_cells.Solved()
}

//Grid which represents the 3x3 collection of squares which represent the entire puzzle
const ROW_LENGTH = 9
const COL_LENGTH = 9
const NUM_SQUARES = COL_LENGTH
const SQUARE_SIZE = 3

type grid struct {
	m_squares 	[]square
	m_rows		[]line
	m_cols		[]line
	
	m_sets		[]Solver
	m_cells		[][]cell
}

func New(puzzle [COL_LENGTH][ROW_LENGTH]int) (grid, *SolveError){
	var g grid
	g.Init();
	g.Fill(puzzle)
	return g,nil
} 

func (g grid) Init() {
	//Init the raw cells themselves that actually store the grid data
	g.m_cells = make([][]cell,COL_LENGTH)
	for i,_ := range g.m_cells{
		g.m_cells[i] = make([]cell, ROW_LENGTH)
	} 
	
	//Init each of the grouping structures that view portions of the grid
	
	/*
	
	Squares are indexed into the grid as folows
	
	S0 S1 S2
	S3 S4 S5
	S6 S7 S8
	
	*/
	g.m_squares = make([]square,NUM_SQUARES)
	
	for squareIdx :=0; squareIdx<NUM_SQUARES; squareIdx++{
		
		g.m_squares[squareIdx].init()
		for x :=0; x<SQUARE_SIZE; x++{
			for y:= 0; y<SQUARE_SIZE; y++{
				gridX := SQUARE_SIZE * (squareIdx % SQUARE_SIZE) + (x%SQUARE_SIZE)
				gridY := SQUARE_SIZE * (squareIdx % SQUARE_SIZE) + (y%SQUARE_SIZE)
				g.m_squares[squareIdx].m_cells[x][y] = &g.m_cells[gridX][gridY]
			}
		}
		
	}
	
	
	g.m_rows = make([]line, ROW_LENGTH)
	g.m_cols = make([]line,COL_LENGTH)
	
	//Make m_sets just a big long list of all the cell grouping structures
	//handy for doing iterations over all different ways of looking at the cells
	g.m_sets = make([]Solver,len(g.m_squares) + len(g.m_rows) + len(g.m_cols))
	
	var idx int
	for _,s := range g.m_squares{
		g.m_sets[idx] = &s
		idx++
	}
	
	for _,r := range g.m_rows{
		g.m_sets[idx] = &r
		idx++ 
	}
	
	for _,c := range g.m_cols{
		g.m_sets[idx] = &c
		idx++ 
	}
}

func (g grid) Fill(puzzle [COL_LENGTH][ROW_LENGTH]int){
	
}

func (g grid) Solved() (bool,*SolveError) {
	for _,s := range g.m_sets{
		solved,err := s.Solved()
		if err != nil{
			fmt.Println("Error during Solved() check on grid: " + err.Error())
			return false,err
		}
		
		if !solved{
			return false,nil
		}
	}
	
	return true,nil
}

type result struct {
	url string
	body string
}

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

