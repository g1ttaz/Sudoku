package sudoku

import (
	"strconv"
)

type BitSet struct {
	bitset uint64
}

func (b* BitSet) clear() {
	b.bitset = 0
}

func (b* BitSet) isSet(bit int) bool {
	ok := (b.bitset & (1 << uint(bit))) != 0
	return ok
}

func (b* BitSet) set(bit int) {
	b.bitset |= (1 << uint(bit))
}

func (b* BitSet) reset(bit int) {
	b.bitset &^= (1 << uint(bit))
}


type Grid struct {
	grid      [][]int
	gridSize    int
	subGridSize int
}

func MakeGrid(grid [][]int, gridSize int, subGridSize int) Grid {
	newGrid := make([][]int, gridSize)
	for row := range newGrid {
		newGrid[row] = make([]int, gridSize)
		for col, val := range grid[row] {
			newGrid[row][col] = val
		}
	}
	return Grid{newGrid, gridSize, subGridSize}
}

func (g Grid) splitLine() string {
	stringified := "+"
	for col := 0; col < g.gridSize; col++ {
		stringified = stringified + "--"
		if col % g.subGridSize == g.subGridSize - 1 {
			stringified = stringified + "-+"
		}
	}
	return stringified
}

func (g Grid) String() string {
	stringified := g.splitLine()
	stringified = stringified + "\n"
	for row := range g.grid {
		stringified = stringified + "|"
		for col, val := range g.grid[row] {
			stringified = stringified + " " + strconv.Itoa(val)
			if col % g.subGridSize == g.subGridSize - 1 {
				stringified = stringified + " |"
			}
		}
		stringified = stringified + "\n"
		if row % g.subGridSize == g.subGridSize - 1 {
			stringified = stringified + g.splitLine()
			stringified = stringified + "\n"		
		}
	}
	return stringified
}

func (g* Grid) copyGrid(gnew Grid) {
	for row := range g.grid {
		for col, _ := range g.grid[row] {
			g.grid[row][col] = gnew.grid[row][col]
		}
	}
}

func (g* Grid) solved() bool {
	for row := range g.grid {
		for val := range g.grid[row] {
			if val == 0 {
				return false
			}
		}
	}
	return true
}

// check whether consistent in all rows
func (g* Grid) consistentInRows() bool {
	for row := range g.grid {
		bitset := BitSet{0}
		for _, val := range g.grid[row] {
			if val != 0 {
				if bitset.isSet(val) {
					return false
				}
				bitset.set(val)
			}
		}
	}
	return true
}	

// check whether consistent in all columns
func (g* Grid) consistentInCols() bool {
	for col := 0; col < g.gridSize; col++ {
		bitset := BitSet{0}
		for row := range g.grid {
			val := g.grid[row][col]
			if val != 0 {
				if bitset.isSet(val) {
					return false
				}
				bitset.set(val)
			}
		}
	}
	return true
}	

// check whether consistent in all subgrids
func (g* Grid) consistentInSubGrids() bool {
	for subGrid := 0; subGrid < g.gridSize; subGrid++ {
		subGridRow, subGridCol := g.subGridSize * (subGrid / g.subGridSize), g.subGridSize * (subGrid % g.subGridSize)
		bitset := BitSet{0}
		for row := subGridRow; row < subGridRow + g.subGridSize; row++ {
			for col := subGridCol; col < subGridCol + g.subGridSize; col++ {
				val := g.grid[row][col]
				if val != 0 {
					if bitset.isSet(val) {
						return false
					}
					bitset.set(val)
				}
			}
		}
	}
	return true
}	

// checks whether grid is consistent
func (g* Grid) consistent() bool {
	return g.consistentInRows() && g.consistentInCols() && g.consistentInSubGrids()
}

func (g* Grid) nextStepRow(row int) (number, column, value int) {
	number, column, value = 0, 0, 0
	bitset := BitSet{0}
	for idx, val := range g.grid[row] {
		if val == 0 {
			number += 1
			column = idx
		} else {
			bitset.set(val)
		}
	}
	if number == (g.gridSize - 1) {
		for i := 1; i <= g.gridSize; i++ {
			if !bitset.isSet(i) {
				value = i
			}
		}
	} else {
		column, value = 0, 0
	}
	return
}

func (g* Grid) nextStepCol(col int) (number, row, value int) {
	number, row, value = 0, 0, 0
	bitset := BitSet{0}
	for idx := 0; idx < g.gridSize; idx++ {
		value = g.grid[idx][col]
		if value == 0 {
			number += 1
			row = idx
		} else {
			bitset.set(value)
		}
	}
	if number == (g.gridSize - 1) {
		for i := 1; i <= g.gridSize; i++ {
			if !bitset.isSet(i) {
				value = i
			}
		}
	} else {
		row = 0
	}
	return
}

func (g* Grid) nextStep() bool {
	found := false
	for row := range g.grid {
		numberOfFreePositions, col, value := g.nextStepRow(row)
		if numberOfFreePositions == 1 {
			g.grid[row][col] = value
			found = true
		}
	}
	for col := 0; col < g.gridSize; col++ {
		numberOfFreePositions, row, value := g.nextStepCol(col)
		if numberOfFreePositions == 1 {
			g.grid[row][col] = value
			found = true
		}
	}
	return found
}

func (g* Grid) bruteForce() bool {
	for row := range g.grid {
		for col, val := range g.grid[row] {
			if val == 0 {
				for i := 1; i <= g.gridSize; i++ {
					gnew := MakeGrid(g.grid,g.gridSize,g.subGridSize)
					gnew.grid[row][col] = i
					if gnew.consistent() {
						solved := gnew.bruteForce()
						if solved {
							g.copyGrid(gnew)
							return true
						}
					}
				}
				return false
			}
		}
	}
	return true
}

func (g* Grid) Solve() bool {
	if !g.consistent() {
		return false
	}
	return g.bruteForce()
}
