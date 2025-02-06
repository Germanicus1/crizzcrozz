package models

import (
	"errors"
	"fmt"
)

// Board represents the entire state of the crossword puzzle.
type Board struct {
	Bounds     *Bounds
	Cells      [][]*Cell
	WordList   map[string]bool
	WordCount  int
	TotalWords int // Total number of words that need to be placed on the board for completion.
	Pool       *Pool
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

	// FIXME: Remove debug info
	// fmt.Printf("CanPlaceWordAt: Testing word %s at start location %+v, in direction %v\n", word, start, direction)

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

	// FIXME: Remove debug info
	// fmt.Printf("CanPlaceWordAt.intersected: %v\n", intersected)

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

		// Check if the cell is out of bounds
		if isOutOfBound(x, y, b) {
			return errors.New("word placement out of bounds")
		}

		// Check for cell conflicts where a cell is already filled with a different character
		if isCellConflict(x, y, b, rune(word[i])) {
			return fmt.Errorf("conflict at position (%d, %d), the cell is occupied by a different character", x, y)
		}

		// Check for valid parallel placement
		if !b.isValidParallelPlacement(x, y, deltaX, deltaY) {
			return fmt.Errorf("invalid parallel word formation at position (%d, %d)", x, y)
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

	//FIXME: remove debug info
	fmt.Printf("Word '%s' placed at (%d, %d) going %s\n", word, start.X, start.Y, directionString(direction))

	return nil
}

// directionString returns a string representation of a Direction.
func directionString(direction Direction) string {
	switch direction {
	case Across:
		return "Across"
	case Down:
		return "Down"
	default:
		return "Unknown Direction"
	}
}

func (b *Board) isValidParallelPlacement(x, y, deltaX, deltaY int) bool {
	// Check perpendicular directions
	perpendicularDeltaX, perpendicularDeltaY := deltaY, deltaX // Swap deltas for perpendicular check
	adjacentCells := []Location{
		{X: x - perpendicularDeltaX, Y: y - perpendicularDeltaY},
		{X: x + perpendicularDeltaX, Y: y + perpendicularDeltaY},
	}

	for _, loc := range adjacentCells {
		if !isOutOfBound(loc.X, loc.Y, b) && b.Cells[loc.Y][loc.X].Filled {
			// If adjacent cell is filled, ensure it does not form an invalid new word
			if !b.isPartOfValidWord(loc, perpendicularDeltaX, perpendicularDeltaY) {
				return false
			}
		}
	}

	return true
}

func (b *Board) isPartOfValidWord(loc Location, deltaX, deltaY int) bool {
	// Generate potential words in both positive and negative directions
	word1 := b.generateWordFromLocation(loc, deltaX, deltaY)
	word2 := b.generateWordFromLocation(loc, -deltaX, -deltaY)

	// Check if either generated word is valid
	return b.isValidWord(word1, b.Pool) || b.isValidWord(word2, b.Pool)
}

// generateWordFromLocation generates a word starting from a given location in a specified direction
func (b *Board) generateWordFromLocation(start Location, deltaX, deltaY int) string {
	var word []rune

	// Move backwards to the start of the word
	x, y := start.X, start.Y
	for b.isValidLocation(x-deltaX, y-deltaY) && b.Cells[y-deltaY][x-deltaX].Filled {
		x -= deltaX
		y -= deltaY
	}

	// Generate the word forward
	for b.isValidLocation(x, y) && b.Cells[y][x].Filled {
		x += deltaX
		y += deltaY
		word = append(word, b.Cells[y][x].Character)
	}

	return string(word)
}

// isValidLocation checks if the given coordinates are within the bounds of the board
func (b *Board) isValidLocation(x, y int) bool {
	return x >= 0 && y >= 0 && x < len(b.Cells[0]) && y < len(b.Cells)
}

// isValidWord checks if the word is in the allowed word list
func (b *Board) isValidWord(word string, pool *Pool) bool {
	return pool.Exists(word)
}

// IsComplete checks if the board is fully set up with all words placed.
func (b *Board) IsComplete() bool {
	return b.WordCount >= b.TotalWords
}
