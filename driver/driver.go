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
	g1 := sudoku.NewGrid(grid,9,3)
	g2 := sudoku.NewGrid(grid,9,3)
	fmt.Printf("Original:\n%v\n", g1)

	g1.SolveBruteForce()
	fmt.Printf("Solved via brute force:\n%v\n", g1)

	g2.SolveBestVal()
	fmt.Printf("Solved via best value:\n%v\n", g2)
}
