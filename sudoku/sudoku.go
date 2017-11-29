package sudoku

import (
	"bytes"
	"fmt"
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
	consistent      bool
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
		if g.valuesInRows[row].isSet(val) {
			g.consistent = false;
		}
		g.valuesInRows[row].set(val)

		if g.valuesInCols[col].isSet(val) {
			g.consistent = false;
		}
		g.valuesInCols[col].set(val)

		subgrid := getSubGrid(row, col, g.subGridSize)
		if g.valuesInSubGrid[subgrid].isSet(val) {
			g.consistent = false;
		}
		g.valuesInSubGrid[subgrid].set(val)
	}
}

// NewGrid creates a sudoku grid
func NewGrid(grid [][]int, gridSize int, subGridSize int) *Grid {
	newGrid := make([][]int, gridSize)
	pixels := make([]int, gridSize*gridSize)
	for row := range newGrid {
		newGrid[row], pixels = pixels[:gridSize], pixels[gridSize:]
	}
	bitsets := make([]BitSet, 3*gridSize)
	valuesInRows := bitsets[:gridSize] 
	valuesInCols := bitsets[gridSize:2*gridSize]
	valuesInSubGrid := bitsets[2*gridSize:3*gridSize]

	consistent := true
	for row := range grid {
		for col, val := range grid[row] {
			newGrid[row][col] = val
			if val != 0 {
				if valuesInRows[row].isSet(val) {
					consistent = false;
				}
				valuesInRows[row].set(val)
				
				if valuesInCols[col].isSet(val) {
					consistent = false;
				}
				valuesInCols[col].set(val)
				
				subgrid := getSubGrid(row, col, subGridSize)
				if valuesInSubGrid[subgrid].isSet(val) {
					consistent = false;
				}
				valuesInSubGrid[subgrid].set(val)
			}
		}
	}
	
	return &Grid{newGrid, gridSize, subGridSize, valuesInRows, valuesInCols, valuesInSubGrid, consistent}
}

// Grid returns the sudoku grid
func (g *Grid) Grid() [][]int {
	return g.grid
}

// String returns a stringified representation of the sudoku grid
func (g Grid) String() string {
	// construct split line
	var b bytes.Buffer
	b.WriteString("+")
	for col := 0; col < g.gridSize; col++ {
		b.WriteString("--")
		if col%g.subGridSize == g.subGridSize-1 {
			b.WriteString("-+")
		}
	}
	b.WriteString("\n")
	splitline := b.String()

	// construct grid
	var b2 bytes.Buffer
	b2.WriteString(splitline)
	for row := range g.grid {
		b2.WriteString("|")
		for col, val := range g.grid[row] {
			fmt.Fprintf(&b2, " %d", val)
			if col%g.subGridSize == g.subGridSize-1 {
				b2.WriteString(" |")
			}
		}
		b2.WriteString("\n")
		if row%g.subGridSize == g.subGridSize-1 {
			b2.WriteString(splitline)
		}
	}
	return b2.String()
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

type NextValFunc func (g *Grid) bool

func (g *Grid) recurse(row, col int, nextValFunc NextValFunc) bool {
	for val := 1; val <= g.gridSize; val++ {
		if g.valuesInRows[row].isSet(val) || g.valuesInCols[col].isSet(val) {
			continue;
		}
		if g.valuesInSubGrid[getSubGrid(row, col, g.subGridSize)].isSet(val) {
			continue;
		}
		gnew := NewGrid(g.grid, g.gridSize, g.subGridSize)
		gnew.setValue(row, col, val)
		if gnew.consistent {
			// fmt.Println("recurse(). row=", row, " col=", col, " val=", val, " grid=", gnew.grid)
			solved := nextValFunc(gnew)
			if solved {
				g.copy(gnew.grid)
				return true
			}
		}
	}
	return false
}

func (g *Grid) bruteForcePosition() (rowPos, colPos int) {
	for row := range g.grid {
		for col, val := range g.grid[row] {
			if val == 0 {
				rowPos, colPos = row, col
				return
			}
		}
	}
	return -1, -1
}

func (g *Grid) bruteForceVal() bool {
	row, col := g.bruteForcePosition()
	if row == -1 || col == -1 {
		return true
	}
	return g.recurse(row, col, (*Grid).bruteForceVal)
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
	return g.recurse(row, bestCol, (*Grid).bestVal)
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
	return g.recurse(bestRow, col, (*Grid).bestVal)
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
	return g.recurse(bestRow, bestCol, (*Grid).bestVal)
}

func bestCount(bitSetArray []BitSet, gridSize int) (bestIndex, bestCount int) {
	for idx, bitset := range bitSetArray {
		if no := bitset.count(gridSize); no > bestCount && no < gridSize {
			bestIndex, bestCount = idx, no
		}
	}
	return
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

// SolveBruteForce solves a sudoku grid via brute force algorithm
func (g *Grid) SolveBruteForce() bool {
	if !g.consistent {
		return false
	}
	return g.bruteForceVal()
}

// SolveBestVal solves a sudoku grid via best value algorithm
func (g *Grid) SolveBestVal() bool {
	if !g.consistent {
		return false
	}
	return g.bestVal()
}
