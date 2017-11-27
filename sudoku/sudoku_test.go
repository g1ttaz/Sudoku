package sudoku

import (
	"testing"
)

func testEqual(a, b [][]int) bool {

	if a == nil && b == nil { 
		return true 
	}

	if a == nil || b == nil { 
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] == nil && b[i] == nil { 
			return true 
		}

		if a[i] == nil || b[i] == nil { 
			return false
		}

		if len(a[i]) != len(b[i]) {
			return false
		}

		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}

	return true
}

func TestSudokuGrid(t *testing.T) {
	for _, c := range []struct {
		in, want [][]int
	}{
		{
			[][] int {
				{0,0,0, 0,0,0, 2,1,8},
				{8,0,0, 0,5,7, 0,0,4},
				{4,3,0, 0,0,1, 0,0,0},

				{2,0,0, 0,0,0, 5,3,7},
				{0,5,0, 1,0,8, 0,4,0},
				{0,4,0, 2,3,0, 0,0,0},

				{0,0,9, 0,0,0, 7,0,0},
				{0,0,2, 0,6,3, 0,0,0},
				{0,0,0, 0,0,0, 8,0,6},
			},
			[][] int {
				{7,9,5, 3,4,6, 2,1,8},
				{8,2,1, 9,5,7, 3,6,4},
				{4,3,6, 8,2,1, 9,7,5},

				{2,1,8, 6,9,4, 5,3,7},
				{9,5,3, 1,7,8, 6,4,2},
				{6,4,7, 2,3,5, 1,8,9},

				{1,6,9, 4,8,2, 7,5,3},
				{5,8,2, 7,6,3, 4,9,1},
				{3,7,4, 5,1,9, 8,2,6},
			},
		},
	} {
		g := MakeGrid(c.in,9,3)
		g.Solve()
		solved := g.grid
		if !testEqual(solved,c.want) {
			t.Errorf("Sudoku(%q) == %q, want %q", c.in, solved, c.want)
		}
	}
}
