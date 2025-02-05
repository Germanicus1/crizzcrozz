package models

import (
	"errors"
	"fmt"
)

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
	intersected := false // To check if at least one letter overlaps with existing words

	for i := 0; i < len(word); i++ {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY

		if isOutOfBound(x, y, b) {
			return false
		}
		if isCellConflict(x, y, b, rune(word[i])) {
			return false
		}
		if b.Cells[y][x].Filled && b.Cells[y][x].Character == rune(word[i]) {
			intersected = true // Ensure the new word intersects at least once
		}
	}

	// Ensure the word intersects with existing words on the board
	return intersected
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

// PlaceWordAt attempts to place a word on the board at the specified location and direction.
// It returns an error if the placement is not possible.
func (b *Board) PlaceWordAt(start Location, word string, direction Direction) error {
	deltaX, deltaY := getDirectionDeltas(direction)

	for i := 0; i < len(word); i++ {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY

		// Check if the current position is out of the board's bounds.
		if y >= len(b.Cells) || x >= len(b.Cells[y]) {
			return errors.New("placement is out of the board's bounds")
		}

		// Check if the cell is already filled with a different character.
		if b.Cells[y][x].Filled && b.Cells[y][x].Character != rune(word[i]) {
			return fmt.Errorf("conflict at position (%d, %d), the cell is already occupied by a different character", x, y)
		}
	}

	// If all checks are passed, place the word on the board.
	for i := 0; i < len(word); i++ {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY
		b.Cells[y][x].Character = rune(word[i])
		b.Cells[y][x].Filled = true
	}
	b.WordCount++
	fmt.Println("Words placed:", b.WordCount)

	return nil
}

// IsComplete checks if the board is fully set up with all words placed.
func (b *Board) IsComplete() bool {
	return b.WordCount >= b.TotalWords
}
