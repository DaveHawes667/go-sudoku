package main

import (
	"fmt"
	"github.com/DaveHawes667/go-sudoku/sudoku"
)

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
	
	var g sudoku.grid
	
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