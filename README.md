# go-sudoku
Sudoku solvers written in go as I have been learning it.

--- example.go

Basic example of how you use the sudoku package, via sudoku.Grid, filling out a puzzle, 
requesting a solve, and how to get the best out of any errors generated by the package.

usage on the command line:
go install src/example.go
bin/example

--- sudoku.go

The package itself with go-routine based solver. Can just import as shown in example.go.

usage in your code:
see example.go

--- sudoku_test.go

Unit tests for the sudoku package.

usage on the command line:
go test github.com/DaveHawes667/go-sudoku/sudoku

--- NOTES
		
Check out the other solvers written in Python & Scala, Rust coming soon. 
