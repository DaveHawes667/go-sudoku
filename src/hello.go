package main

import (
	"fmt"
	"strconv"
	"github.com/go-errors/errors"
)

//General solver error
type SolveError struct{
	m_info string
}

func (e SolveError) Error() string {
	return e.String()
}

func (e SolveError) String() string {
	return e.m_info
}

//General solver interface
type Solver interface {
	Solved() (bool,error)
	reducePossible() (bool,error)
	validate() bool
}

func validate(known []int) bool{
	
	validNums := map[int]bool{1:true,2:true,3:true,4:true,5:true,6:true,7:true,8:true,9:true}

	for _,v := range known{
		if validNums[v]{
			validNums[v] = false
		}else{
			fmt.Println("same number is known in two cells, that is invalid!")
			return false 
		}
	}
	
	return true
}

//Cells which store individual numbers in the grid
type cell struct {
	m_possible map[int]bool 
}

func (c *cell) init(){
	c.m_possible = make(map[int]bool,9)
	for i:=1; i<=9; i++{
		c.m_possible[i] = true
	}
}

func (c *cell) validate() bool{
	if len(c.m_possible) < 1{
		fmt.Println("Less than one possible in a cell, that is invalid")
		return false
	}
	if len(c.m_possible) > 9{
		fmt.Println("More than nine possible in a cell, that is invalid")
		return false
	}
	return true
}

func (c *cell) SetKnownTo(value int){
	for k,_ := range c.m_possible{
		if k != value {
			delete(c.m_possible,k)
		}
	} 
}

func (c *cell) TakeKnownFromPossible(known []int) (bool,error){

	var possibles = len(c.m_possible)
	
	if !c.IsKnown(){		
		for _,v := range known{
			delete(c.m_possible,v)
		}
		
		if !c.validate(){
			errorStr := fmt.Sprintf("TakeKnownFromPossible() made cell invalid: Known:%T:%q", known,known)
			return false,errors.Wrap(SolveError{errorStr},1)
		}
	}
	
	return possibles != len(c.m_possible),nil //did we take any?
}

func (c* cell) IsKnown() bool{
	return len(c.m_possible) == 1
}

func (c *cell) Known() (int, error){
	
	// by convention we delete from the map possibles that are no longer possible
	// so we just need to check map length to see if the cell is solved
	if len(c.m_possible)!=1{
		return 0,errors.Wrap(SolveError{"Value not yet known for this cell"},1)
	}
	
	//Only one key is now considered "possible", it's value should be true, and it should
	//be the only one in the list, return it if that is the case 
	for k,v := range c.m_possible {
		if v{
			return k,nil
		}
	}
	
	return 0,errors.Wrap(SolveError{"Error in cell storage of known values"},1)
}

func (c *cell) Possibles() ([]int,error){
	possibles := make([]int,0,len(c.m_possible))
	
	for k,v := range c.m_possible{
		if v{
			possibles = append(possibles,k)
		}
	}
	
	return possibles,nil
}

func (c cell) String() string {
	val,err := c.Known()
	if(err != nil){
		return "x"
	}
	
	return strconv.Itoa(val)
}

type cellPtrSlice []*cell 

func (cells cellPtrSlice) Solved() (bool, error){
	for _,c := range cells{
		if !c.IsKnown(){
			return false,nil
		}
	}
	
	return true,nil
}

func (cells cellPtrSlice) Known() ([]int, error){
	known := make([]int,0, len(cells))
	for _,c := range cells{
		if c.IsKnown(){
			val,err := c.Known()
			if err != nil{
				return known,err
			}
			known = append(known,val)
		}
	}
	
	return known,nil
}

func (cells cellPtrSlice) TakeKnownFromPossible(known []int) (bool,error){
	
	changed := false
	for _,c := range cells{
		taken, err := c.TakeKnownFromPossible(known)
		if err != nil{
			return false,err
		}
		changed = changed || taken
	}
	
	return changed,nil
}

//Squares which represent one of each of the 9 squares in a grid, each of which 
//references a 3x3 collection of cells.

type square struct {
	m_cells [][]*cell
}

func (s square) Solved() (bool,error) {
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
func (s *square) init() {
	s.m_cells = make([][]*cell,SQUARE_SIZE)
	for i,_ := range s.m_cells{
		s.m_cells[i] = make(cellPtrSlice, SQUARE_SIZE)
	}
}

func (s* square) KnownInSquare() ([]int,error){
	known := make([]int,0,SQUARE_SIZE*SQUARE_SIZE)
	for x,_ := range s.m_cells{
		for y,_ := range s.m_cells[x]{
			c := s.m_cells[x][y]
			if c.IsKnown(){
				val,err := c.Known()
				if err != nil{
					return known,err
				}
				known = append(known,val)	
			}	
		}
	}
	
	return known,nil
}

func (s* square) reducePossible() (bool,error) {
	known,err := s.KnownInSquare()
	reduced := false
	if err != nil {
		return false,err
	}
	
	for x,_ := range s.m_cells{
		cells := s.m_cells[x]
		changed, err := cellPtrSlice(cells).TakeKnownFromPossible(known)
		if err != nil{
			return false,err
		}
		reduced = reduced || changed
	}
	return reduced,nil
}

func (s* square) validate() bool{
	known,err := s.KnownInSquare()
	
	if err != nil {
		fmt.Println("Error from KnownInSquare, that is invalid")
		return false
	}
	
	return validate(known)
}

//A horizontal or vertical line of 9 cells through the entire grid.
type line struct {
	m_cells cellPtrSlice
}

func (l line) Solved() (bool,error) {
	return l.m_cells.Solved()
}

func (l* line) reducePossible() (bool,error) {
	known,err := l.m_cells.Known()
	
	if err != nil {
		return false,err
	}
	
	reduced, err := l.m_cells.TakeKnownFromPossible(known)
	if err != nil{
		return false,err
	}
		
	return reduced,nil
}

func (l* line) validate() bool{
	known,err := l.m_cells.Known()
	
	if err != nil {
		fmt.Println("Error in line cells Known(), that is invalid")
		return false
	}
	return validate(known)
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

func New(puzzle [COL_LENGTH][ROW_LENGTH]int) (*grid, error){
	var g grid
	g.Init();
	g.Fill(puzzle)
	return &g,nil
} 

func (g *grid) validate() bool{
	for x,_ := range g.m_cells{
		for y,_:= range g.m_cells[x]{
			if !g.m_cells[x][y].validate(){
				fmt.Println("cell failed to validate: x:"+strconv.Itoa(x)+" y:"+strconv.Itoa(y))
				fmt.Println(g.m_cells[x][y].m_possible)
				
				return false
			}	
		}
	}
	
	for i,_ := range g.m_sets{
		if !g.m_sets[i].validate(){
			return false
		}
	}
	
	return true
}

func (g *grid) Init() {
	//Init the raw cells themselves that actually store the grid data
	g.m_cells = make([][]cell,COL_LENGTH)
	for i :=0; i< len(g.m_cells); i++{
		g.m_cells[i] = make([]cell, ROW_LENGTH)
		
		for j :=0; j< len(g.m_cells[i]); j++{
			c := &g.m_cells[i][j]
			c.init()
		}
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
				//is this correct?
				gridX := SQUARE_SIZE * (squareIdx % SQUARE_SIZE) + x 
				gridY := SQUARE_SIZE * (squareIdx / SQUARE_SIZE) + y
				
				cellPtr := &g.m_cells[gridX][gridY]
				g.m_squares[squareIdx].m_cells[x][y] = cellPtr
			}
		}
		
	}
	
	g.m_rows = make([]line, ROW_LENGTH)
	g.m_cols = make([]line,COL_LENGTH)
	
	//Make m_sets just a big long list of all the cell grouping structures
	//handy for doing iterations over all different ways of looking at the cells
	g.m_sets = make([]Solver,len(g.m_squares) + len(g.m_rows) + len(g.m_cols))
	
	
	var idx int
	for i := 0; i<len(g.m_squares); i++{
		s:= &g.m_squares[i]
		g.m_sets[idx] = s
		idx++
	}
	
	
	for i:=0; i<len(g.m_rows); i++{
		r:=&g.m_rows[i]
		g.m_sets[idx] = r
		idx++ 
	}
	
	
	for i:= 0; i<len(g.m_cols); i++{
		c := &g.m_cols[i]
		g.m_sets[idx] = c
		idx++ 
	}
	
}

func (g *grid) Fill(puzzle [COL_LENGTH][ROW_LENGTH]int){
	g.Init()
	
	
	for x:=0; x<COL_LENGTH; x++{
		for y:=0; y<ROW_LENGTH; y++{
			var puzzVal = puzzle[y][x]
			
			if puzzVal >=1 && puzzVal<=9{
				g.m_cells[x][y].SetKnownTo(puzzVal)
			}
		}
	}
	
}

func (g grid) Solved() (bool,error) {
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

func (g *grid) squareExclusionReduce() (bool,error){

	return false,nil	
}

func(g *grid) reducePossiblePass() (bool, error){
	changed := false
	
	for pass:=0;pass<2;pass++{
		for i,_ := range g.m_sets{
			reduced,err := g.m_sets[i].reducePossible()
			if err != nil{
				return false,err
			}
			changed = changed || reduced
		}
		sqReduce,err := g.squareExclusionReduce()
		if err != nil{
				return false,err
		}
		changed = changed || sqReduce
	}
	
	return changed,nil
}

func (g *grid) Puzzle() [COL_LENGTH][ROW_LENGTH]int{
	var puzzle [COL_LENGTH][ROW_LENGTH]int
	for x,_ := range puzzle{
		for y,_ := range puzzle[x]{
			if g.m_cells[x][y].IsKnown(){
				var err error
				puzzle[y][x],err = g.m_cells[x][y].Known()
				if err != nil{
					return puzzle
				}
			}
		}
	}
	
	return puzzle
}

func (g *grid) setKnown( x,y, known int) error{
	//should probably check if grid is initialised and return error if it isn't
	g.m_cells[x][y].SetKnownTo(known)
	
	return nil
}

func (g *grid) DuplicateGrid() (*grid,error){
	return New(g.Puzzle())
}

func (g* grid) TotalPossible() (int, error){
	totalPoss := 0
	for x,_ := range g.m_cells{
		for y,_:= range g.m_cells[x]{
			if !g.m_cells[x][y].IsKnown(){
				val,err := g.m_cells[x][y].Possibles()
				if err != nil{
					return 0,err
				}
				numPoss := len(val)
				totalPoss += numPoss
			}
		}
		
	}	
	return totalPoss,nil
}

func (g *grid) GenerateGuessGrids() ([]*grid, error){
	
	totalPoss,err := g.TotalPossible()
	guesses := make([]*grid,0,totalPoss)
	if err != nil{
		return guesses,err
	}
	
	for x,_ := range g.m_cells{
		for y,_:= range g.m_cells[x]{
			if !g.m_cells[x][y].IsKnown(){
				
				possibles,err := g.m_cells[x][y].Possibles()
				if err!=nil{
					return guesses,err
				}
				
				for _,v := range possibles{
					guess,err := g.DuplicateGrid()
					if err!=nil{
						return guesses,err
					}
					err = guess.setKnown(x,y,v)
					if err!=nil{
						return guesses,err
					}
					guesses = append(guesses,guess)
				}
			}
		}
	}
	
	
	return guesses,nil
	
}

type SolveResult struct{
	m_grid 		*grid
	m_solved 	bool
}

func startSolveRoutine(ch chan SolveResult, g *grid) {
	
	defer close(ch)
	//fmt.Print("goroutine solver started")
	res, err := g.Solve()
	fmt.Println("return from solve")
	if err != nil{
		//this error might be expected, we might have sent in an invalid puzzle
		//only care about this response to print or pass on in the root call to solve.
		return
	}
	 
	ch<-*res
}

func (g *grid) Solve() (*SolveResult,error){
	
	var err error
	//solvePasses:=0
	for changed:=true; changed;{
		changed, err = g.reducePossiblePass()
		//fmt.Println("AFTER REDUCE PASS " + strconv.Itoa(solvePasses))
		//solvePasses++
		//fmt.Println(g)
		if err != nil{
			return &SolveResult{nil,false},err
		}
		if !g.validate(){
			fmt.Println("Invalid puzz!")
			return &SolveResult{nil,false},nil //this was probably an invalid guess, just want to stop trying to process this
		}
	}
	
	solved,err := g.Solved()
	if err != nil{
			return &SolveResult{nil,false},err
	}
	
	if solved{
		return &SolveResult{g,true},nil
	}
	
	//fmt.Println("About to generate guess grids")
	guesses,err := g.GenerateGuessGrids()
	//fmt.Println("Generated "+strconv.Itoa(len(guesses))+" potential grids")
	
	if err!=nil{
		return &SolveResult{g,false},err
	}
	
	resChans := make([]chan SolveResult,0,len(guesses))
	
	for _,guess := range guesses{
		ch := make(chan SolveResult)
		go startSolveRoutine(ch,guess)
		resChans = append(resChans,ch)
	}
	
	for i,_ := range resChans{
		for res:= range resChans[i] {
			if res.m_solved{
				return &res,nil
			}else{
				fmt.Println("Failed result returned on channel")
			}
		}
	}
	
	return &SolveResult{nil,false},errors.Wrap(SolveError{"Unable to solve puzzle"},1)
}

func (g grid) String() string {
	var str string
	
	if len(g.m_cells) < COL_LENGTH{
		return "Grid probably not initialised"
	}
	
	for x:=0; x < ROW_LENGTH; x++{
		for y:=0; y < COL_LENGTH; y++{
			
			if len(g.m_cells[y]) < ROW_LENGTH{
				return "Error in grid not all rows correctly initialised"
			}
			
			cell := g.m_cells[y][x]
			str += cell.String()
		}
		str+="\n"
	}
	return str
}

// Main code

func main() {
	
	puzzle := [9][9]int{	{3,0,0,9,6,0,0,0,0},
							{1,4,0,0,0,5,0,9,0},
							{0,0,5,0,0,0,0,0,8},
							{0,0,0,0,5,0,0,2,0},
							{0,0,3,8,0,0,0,1,9},
							{0,0,0,6,4,0,0,3,0},
							{0,0,0,0,0,0,0,0,1},
							{8,0,0,0,2,0,0,0,0},
							{0,0,1,0,0,3,0,0,4},
	}
	
	var g grid
	
	g.Fill(puzzle)
	fmt.Println("Puzzle to solve")
	fmt.Println(g)
	
	res,err := g.Solve()
	if err != nil{
		fmt.Println("Error solving puzzle: " + err.Error())
		fmt.Println("Stacktrace")
		fmt.Println(err.(*errors.Error).ErrorStack())
	}else{
		fmt.Println("")
		
		if res.m_solved && res.m_grid != nil{
			fmt.Println("Solution Found")
			fmt.Println(*res.m_grid)	
		}else{
			fmt.Println("unable to solve puzzle")
		}
		
	}
	
}

