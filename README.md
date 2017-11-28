# Sudoku
A simple sudoku solver.

This is my first Go project where i try to play with the features.

Creating a Sudoku grid:
```
	grid := [][] int {
		...
	}
	g := sudoku.NewGrid(grid,9,3)
```

Solving the grid:
```
	// returns false, if the input grid is not solveable
	// true, otherwise
	solved := g.Solve()
```

Obtain the result:
```
	// returns the solved grid
	grid := g.Grid()
```

Supported algorithm:
- brute force: simply try the "next" free grid position with all possible values
- best value: evaluate the "best" row, column or subgrid and try it with all possible values
- both algorithms work recursively
