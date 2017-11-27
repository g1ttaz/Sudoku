package sudoku

import (
	"strconv"
)

// BitSet is a set of integers
type BitSet struct {
	bitset uint64
}

func (b *BitSet) clear() {
	b.bitset = 0
}

func (b *BitSet) isSet(bit int) bool {
	return (b.bitset & (1 << uint(bit))) != 0
}

func (b *BitSet) set(bit int) {
	b.bitset |= (1 << uint(bit))
}

func (b *BitSet) reset(bit int) {
	b.bitset &^= (1 << uint(bit))
}

func (b *BitSet) count(len int) int {
	number := 0
	for bit := 0; bit <= len; bit++ {
		if b.isSet(bit) {
			number++
		}
	}
	return number
}

// Grid describes the sudoku grid
type Grid struct {
	grid            [][]int
	gridSize        int
	subGridSize     int
	valuesInRows    []BitSet
	valuesInCols    []BitSet
	valuesInSubGrid []BitSet
}

func getSubGrid(row, col, subGridSize int) int {
	return subGridSize*(row/subGridSize) + col/subGridSize
}

func getSubGridOffset(subGrid, subGridSize int) (rowOffset, colOffset int) {
	rowOffset, colOffset = subGridSize*(subGrid/subGridSize), subGridSize*(subGrid%subGridSize)
	return
}

func (g *Grid) copy(gnew [][]int) {
	for _, bitset := range g.valuesInRows {
		bitset.clear()
	}
	for _, bitset := range g.valuesInCols {
		bitset.clear()
	}
	for _, bitset := range g.valuesInSubGrid {
		bitset.clear()
	}
	for row := range gnew {
		for col, val := range gnew[row] {
			g.setValue(row,col,val)
		}
	}
}

func (g *Grid) setValue(row, col, val int) {
	g.grid[row][col] = val
	if val != 0 {
		g.valuesInRows[row].set(val)
		g.valuesInCols[col].set(val)
		g.valuesInSubGrid[getSubGrid(row, col, g.subGridSize)].set(val)
	}
}

// MakeGrid creates a sudoku grid
func MakeGrid(grid [][]int, gridSize int, subGridSize int) Grid {
	newGrid := make([][]int, gridSize)
	pixels := make([]int, gridSize*gridSize)
	for row := range newGrid {
		newGrid[row], pixels = pixels[:gridSize], pixels[gridSize:]
	}
	bitsets := make([]BitSet, 3*gridSize)
	valuesInRows := bitsets[:gridSize] 
	valuesInCols := bitsets[gridSize:2*gridSize]
	valuesInSubGrid := bitsets[2*gridSize:3*gridSize]

	for row := range grid {
		for col, val := range grid[row] {
			newGrid[row][col] = val
			if val != 0 {
				valuesInRows[row].set(val)
				valuesInCols[col].set(val)
				valuesInSubGrid[getSubGrid(row, col, subGridSize)].set(val)
			}
		}
	}
	
	return Grid{newGrid, gridSize, subGridSize, valuesInRows, valuesInCols, valuesInSubGrid}
}

func (g Grid) splitLine() string {
	stringified := "+"
	for col := 0; col < g.gridSize; col++ {
		stringified = stringified + "--"
		if col%g.subGridSize == g.subGridSize-1 {
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
			if col%g.subGridSize == g.subGridSize-1 {
				stringified = stringified + " |"
			}
		}
		stringified = stringified + "\n"
		if row%g.subGridSize == g.subGridSize-1 {
			stringified = stringified + g.splitLine()
			stringified = stringified + "\n"
		}
	}
	return stringified
}

func (g *Grid) solved() bool {
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
func (g *Grid) consistentInRows() bool {
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
func (g *Grid) consistentInCols() bool {
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
func (g *Grid) consistentInSubGrids() bool {
	for subGrid := 0; subGrid < g.gridSize; subGrid++ {
		subGridRow, subGridCol := getSubGridOffset(subGrid, g.subGridSize)
		bitset := BitSet{0}
		for row := subGridRow; row < subGridRow+g.subGridSize; row++ {
			for col := subGridCol; col < subGridCol+g.subGridSize; col++ {
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
func (g *Grid) consistent() bool {
	return g.consistentInRows() && g.consistentInCols() && g.consistentInSubGrids()
}

func (g *Grid) bruteForce() bool {
	for row := range g.grid {
		for col, val := range g.grid[row] {
			if val == 0 {
				for i := 1; i <= g.gridSize; i++ {
					gnew := MakeGrid(g.grid, g.gridSize, g.subGridSize)
					gnew.setValue(row,col,i)
					if gnew.consistent() {
						solved := gnew.bruteForce()
						if solved {
							g.copy(gnew.grid)
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

func bestCount(bitSetArray []BitSet, gridSize int) (bestIndex, bestCount int) {
	for idx, bitset := range bitSetArray {
		if no := bitset.count(gridSize); no > bestCount && no < gridSize {
			bestIndex, bestCount = idx, no
		}
	}
	return
}

func (g *Grid) recurse(row, col int) bool {
	for val := 1; val <= g.gridSize; val++ {
		gnew := MakeGrid(g.grid, g.gridSize, g.subGridSize)
		gnew.setValue(row, col, val)
		if gnew.consistent() {
			// fmt.Println("recurse(). row=", row, " col=", col, " val=", val, " grid=", gnew.grid, " valuesInRows=", gnew.valuesInRows)
			solved := gnew.bestVal()
			if solved {
				g.copy(gnew.grid)
				return true
			}
		}
	}
	return false
}

func (g *Grid) bestValRow(row int) bool {
	bestCol, valuesSetCols := -1, 0
	for col, val := range g.grid[row] {
		if val == 0 {
			if no := g.valuesInCols[col].count(g.gridSize); no > valuesSetCols /*&& no < g.gridSize*/ {
				bestCol, valuesSetCols = col, no
			}
		}
	}

	// fmt.Println("bestValRow(). row=", row, " bestCol=", bestCol, " valuesSetCols=", valuesSetCols, "grid=", g.grid)
	if bestCol == -1 {
		return true
	}
	return g.recurse(row, bestCol)
}

func (g *Grid) bestValCol(col int) bool {
	bestRow, valuesSetRows := -1, 0
	for row := 0; row < g.gridSize; row++ {
		val := g.grid[row][col]
		if val == 0 {
			if no := g.valuesInRows[row].count(g.gridSize); no > valuesSetRows /*&& no < g.gridSize*/ {
				bestRow, valuesSetRows = row, no
			}
		}
	}

	// fmt.Println("bestValCol(). col=", col, " bestRow=", bestRow, " valuesSetRows=", valuesSetRows, " grid=", g.grid)
	if bestRow == -1 {
		return true
	}
	return g.recurse(bestRow, col)
}

func (g *Grid) bestValSubGrid(subGrid int) bool {
	bestRow, bestCol, valuesSet := -1, -1, 0
	rowOffset, colOffset := getSubGridOffset(subGrid, g.subGridSize)

	for row := rowOffset; row < rowOffset+g.subGridSize; row++ {
		for col := colOffset; col < colOffset+g.subGridSize; col++ {
			val := g.grid[row][col]
			if val == 0 {
				if no := g.valuesInRows[row].count(g.gridSize); no > valuesSet /*&& no < g.gridSize*/ {
					bestRow, bestCol, valuesSet = row, col, no
				}
				if no := g.valuesInCols[col].count(g.gridSize); no > valuesSet /*&& no < g.gridSize*/ {
					bestRow, bestCol, valuesSet = row, col, no
				}
			}
		}
	}

	// fmt.Println("bestValSubGrid(). subGrid=", subGrid, " bestRow=", bestRow, " bestCol=", bestCol, "grid=", g.grid)
	if bestRow == -1 || bestCol == -1 {
		return true
	}
	return g.recurse(bestRow, bestCol)
}

func (g *Grid) bestVal() bool {
	bestRow, valuesSetRows := bestCount(g.valuesInRows, g.gridSize)
	bestCol, valuesSetCols := bestCount(g.valuesInCols, g.gridSize)
	bestSubGrid, valuesSetSubGrid := bestCount(g.valuesInSubGrid, g.gridSize)

	// fmt.Println("bestVal(): grid=", g.grid, "bestRow=", bestRow, "bestCol=", bestCol, "bestSubGrid=", bestSubGrid,
	//	"valuesSetRows=", valuesSetRows, "valuesSetCols=", valuesSetCols,
	//	"valuesSetSubGrid=", valuesSetSubGrid)

	if valuesSetRows >= valuesSetCols {
		if valuesSetRows >= valuesSetSubGrid {
			return g.bestValRow(bestRow)
		}
	} else {
		if valuesSetCols >= valuesSetSubGrid {
			return g.bestValCol(bestCol)
		}
	}
	return g.bestValSubGrid(bestSubGrid)
}

// Solve solves a sudoku grid
func (g *Grid) Solve() bool {
	if !g.consistent() {
		return false
	}
	return g.bestVal()
}
