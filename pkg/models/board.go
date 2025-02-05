package models

// Board represents the entire state of the crossword puzzle.
type Board struct {
	Bounds     *Bounds
	Cells      [][]*Cell
	WordCount  int
	TotalWords int // Total number of words that need to be placed on the board for completion.
}

// NewBoard creates a new board with specified bounds and total words.
func NewBoard(bounds *Bounds, totalWords int) *Board {
	width := bounds.Width()
	height := bounds.Height()
	cells := make([][]*Cell, height)
	for i := range cells {
		cells[i] = make([]*Cell, width)
		for j := range cells[i] {
			cells[i][j] = NewEmptyCell() // Assuming NewEmptyCell initializes an empty cell.
		}
	}
	return &Board{
		Bounds:     bounds,
		Cells:      cells,
		TotalWords: totalWords,
	}
}

// CanPlaceWordAt checks if a word can be placed at a specific location
// in a given direction.
func (b *Board) CanPlaceWordAt(start Location, word string, direction Direction) bool {
	deltaX, deltaY := getDirectionDeltas(direction)

	for i := 0; i < len(word); i++ {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY

		if isOutOfBound(x, y, b) {
			return false
		}
		if isCellConflict(x, y, b, rune(word[i])) {
			return false
		}
	}
	return true
}

// getDirectionDeltas determines the increments (deltaX, deltaY) for x
// and y based on the direction.
func getDirectionDeltas(direction Direction) (int, int) {
	if direction == Across {
		return 1, 0
	}
	return 0, 1
}

// isOutOfBound checks if a cell fits on the board
func isOutOfBound(x, y int, b *Board) bool {
	return x >= len(b.Cells[0]) || y >= len(b.Cells)
}

// isCellConflict checks if there are any conflicts character. The
// conflict tested for is: overlapping character matches the current
// charater or not.
func isCellConflict(x, y int, b *Board, char rune) bool {
	cell := b.Cells[y][x]
	return cell.Filled && cell.Character != char
}

// PlaceWordAt places a word on the board at the specified location and
// direction.
func (b *Board) PlaceWordAt(start Location, word string, direction Direction) {
	// Logic to place the word on the board.
	b.WordCount++
}

// IsComplete checks if the board is fully set up with all words placed.
func (b *Board) IsComplete() bool {
	return b.WordCount >= b.TotalWords
}
