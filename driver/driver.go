package main

import (
	"fmt"
	"github.com/g1ttaz/Sudoku/sudoku"
)

func main() {
	grid := [][] int {
		{0,0,0, 0,0,0, 2,1,8},
		{8,0,0, 0,5,7, 0,0,4},
		{4,3,0, 0,0,1, 0,0,0},

		{2,0,0, 0,0,0, 5,3,7},
		{0,5,0, 1,0,8, 0,4,0},
		{0,4,0, 2,3,0, 0,0,0},

		{0,0,9, 0,0,0, 7,0,0},
		{0,0,2, 0,6,3, 0,0,0},
		{0,0,0, 0,0,0, 8,0,6},		
	}
	g := sudoku.NewGrid(grid,9,3)
	fmt.Printf("Original:\n%v\n", g)

	g.Solve()
	fmt.Printf("Solved:\n%v\n", g)
}
