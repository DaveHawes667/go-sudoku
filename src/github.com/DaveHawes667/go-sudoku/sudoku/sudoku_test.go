package sudoku

import (
		"testing"
		"strconv"
)

func TestGrid(t *testing.T){
	puzzles := [1][9][9]int{	
								{
									{3,0,0,9,6,0,0,0,0},
									{1,4,0,0,0,5,0,9,0},
									{0,0,5,0,0,0,0,0,8},
									{0,0,0,0,5,0,0,2,0},
									{0,0,3,8,0,0,0,1,9},
									{0,0,0,6,4,0,0,3,0},
									{0,0,0,0,0,0,0,0,1},
									{8,0,0,0,2,0,0,0,0},
									{0,0,1,0,0,3,0,0,4},
								},
	}
	
	solutions := [1][9][9]int{	
								{
									{3,7,2,9,6,8,1,4,5},
									{1,4,8,7,3,5,6,9,2},
									{9,6,5,2,1,4,3,7,8},
									{4,1,7,3,5,9,8,2,6},
									{6,5,3,8,7,2,4,1,9},
									{2,8,9,6,4,1,5,3,7},
									{5,3,6,4,9,7,2,8,1},
									{8,9,4,1,2,6,7,5,3},
									{7,2,1,5,8,3,9,6,4},
								},
	}
	
	for i,puzz := range puzzles{
		var g Grid
		g.Fill(puzz)
		
		res,err := g.Solve()
		if err != nil{
			t.Errorf("Test " + strconv.Itoa(i) + ": " + err.Error())
		}else{
			if res.Solved() && res.Grid() != nil{
				if !res.Grid().KnownEquals(solutions[i]){
					t.Errorf("Test " + strconv.Itoa(i) + ": Solution found, but did not match expected solution instead got\n" + res.Grid().String())
				}	
			}else{
				t.Errorf("Test " + strconv.Itoa(i) + ": unable to solve puzzle")
			}
			
		}
	}
}