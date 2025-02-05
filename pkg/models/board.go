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
	// Add logic to determine if a word can be placed here based on
	// current board state.
	return true // Placeholder
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
