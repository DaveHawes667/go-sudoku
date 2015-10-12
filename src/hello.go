package main

import (
	"fmt"
	"strconv"
)

//General solver error
type SolveError struct{
	m_info string
}

func (e SolveError) Error() string {
	return e.m_info
}

func (e SolveError) String() string {
	return e.Error()
}

//General solver interface
type Solver interface {
	Solved() (bool,*SolveError)
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

func (c *cell) SetKnownTo(value int){
	for k,_ := range c.m_possible{
		if k != value {
			delete(c.m_possible,k)
		}
	} 
}

func (c *cell) TakeKnownFromPossible(known []int) bool{
	
	var possibles = len(c.m_possible)
	
	for _,v := range known{
		delete(c.m_possible,v)
	}
	
	return possibles != len(c.m_possible) //did we take any?
}

func (c* cell) IsKnown() bool{
	return len(c.m_possible) == 1
}

func (c *cell) Known() (int, *SolveError){
	
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

func (c cell) String() string {
	val,err := c.Known()
	if(err != nil){
		return "x"
	}
	
	return strconv.Itoa(val)
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
func (s *square) init() {
	s.m_cells = make([][]*cell,SQUARE_SIZE)
	for i,_ := range s.m_cells{
		s.m_cells[i] = make(cellPtrSlice, SQUARE_SIZE)
	}
}

func (s* square) KnownInSquare() ([]int,*SolveError){
	known := make([]int,0,SQUARE_SIZE*SQUARE_SIZE)
	for x,_ := range s.m_cells{
		for y,_ := range s.m_cells[x]{
			c := s.m_cells[x][y]
			val,err := c.Known()
			if err != nil{
				return known,err
			}
			known = append(known,val)
		}
	}
	
	return known,nil
}

func (s* square) ReducePossible() (bool,*SolveError) {
	known,err := s.KnownInSquare()
	reduced := false
	if err != nil {
		return false,err
	}
	
	for x,_ := range s.m_cells{
		for y,_ := range s.m_cells[x]{
			c := s.m_cells[x][y]
			if !c.IsKnown(){
				if c.TakeKnownFromPossible(known){
					reduced = true
				}
			}
			
		}
	}
	return reduced,nil
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

func New(puzzle [COL_LENGTH][ROW_LENGTH]int) (*grid, *SolveError){
	var g grid
	g.Init();
	g.Fill(puzzle)
	return &g,nil
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
	fmt.Println(g)
}

