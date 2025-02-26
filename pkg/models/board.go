package models

import (
	"encoding/json"
	"os"
)

type PlacedWord struct {
	Start     Location
	Direction Direction
	Word      string
}

// Board represents the entire state of the crossword puzzle.
type Board struct {
	Bounds      *Bounds
	PlacedWords []PlacedWord
	Cells       [][]*Cell
	WordList    map[string]bool
	WordCount   int
	TotalWords  int // Total number of words that need to be placed on the board for completion.
	Pool        *Pool
}

// NewBoard creates a new board with specified bounds and total words.
func NewBoard(bounds *Bounds, totalWords int) *Board {
	width := bounds.Width()
	height := bounds.Height()
	cells := make([][]*Cell, height)
	for i := range cells {
		cells[i] = make([]*Cell, width)
		for j := range cells[i] {
			cells[i][j] = NewEmptyCell()
		}
	}
	return &Board{
		Bounds:     bounds,
		Cells:      cells,
		TotalWords: totalWords,
		WordCount:  0,
	}
}

// Save converts the Board struct to JSON and writes it to the file structure.
//
// Returns: Error if marshalling or file writing did not work; nil otherwise.
func (b *Board) Save() error {
	data, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("board.json", data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// CanPlaceWordAt determines if a word can be legally placed on the board at a
// specified location and direction. It checks that the word does not overflow
// the board, does not conflict with existing letters, intersects correctly with
// at least one existing letter, and does not create invalid adjacent parallel
// words. Parameters:
//
//	start - The starting location (x, y) for placing the word.
//	word - The word to be placed.
//	direction - The direction to place the word (e.g., horizontal or vertical).
//
// Returns:
//
//	true if the word can be placed according to the rules of the game; false otherwise.
func (b *Board) CanPlaceWordAt(start Location, word string, direction Direction) bool {
	deltaX, deltaY := getDirectionDeltas(direction) // Get the direction deltas to determine how to increment the position.
	intersected := false                            // Flag to track if the word intersects at least once with existing words.
	intersectionCount := 0                          // Counter for the number of intersections with existing words.

	// Check if the placement of the entire word would be within the board's
	// bounds.
	if !b.isPlacementWithinBounds(start, word, deltaX, deltaY) {
		return false
	}

	// Loop through each character in the word to check placement rules.
	for i := 0; i < len(word); i++ {
		x := start.X + i*deltaX // Calculate x position of the current character.
		y := start.Y + i*deltaY // Calculate y position of the current character.
		isIntersection := false // Local flag to check if the current character intersects with a filled cell.

		// Check if placing the current character causes a conflict with
		// different letters already placed.
		if isCellConflict(x, y, b, string(word[i])) {
			return false
		}

		// Check if the current placement intersects correctly without
		// overlapping incorrectly.
		if b.Cells[y][x].Filled && b.Cells[y][x].Character == string(word[i]) {
			intersected = true
			isIntersection = true
			intersectionCount++
			if intersectionCount > 1 { // Ensure only one intersection to avoid multiple overlaps.
				return false
			}
		}

		// Check for invalid adjacent parallel placements.
		if !isIntersection {
			if isParallelPlacement(x, y, direction, b) {
				return false
			}
		}
	}

	// Check if cells immediately before and after the word are unoccupied to
	// prevent contiguous word formation.
	xBefore, yBefore := start.X-deltaX, start.Y-deltaY
	xAfter, yAfter := start.X+len(word)*deltaX, start.Y+len(word)*deltaY

	if isOutOfBound(xBefore, yBefore, b) || isCellFilled(xBefore, yBefore, b) || isOutOfBound(xAfter, yAfter, b) || isCellFilled(xAfter, yAfter, b) {
		return false
	}

	// Ensure the word intersects at least once with existing words on the
	// board.
	return intersected
}

// isPlacementWithinBounds checks if the placement of the last character of the
// word at the calculated position is within the boundaries of the board. It
// calculates the position based on the starting location, the length of the
// word, and the direction in which the word is placed. Parameters:
//
// start - The starting location (x, y) for placing the word.
//
// word - The word to be placed.
//
// deltaX - The horizontal direction increment (1 for rightward, -1 for
// leftward, 0 for none.)
//
// deltaY - The vertical direction increment (1 for downward, -1 for upward, 0
// for noåne).
//
// Returns:
//
//	true if the last character of the word fits within the board boundaries; false otherwise.
func (b *Board) isPlacementWithinBounds(start Location, word string, deltaX, deltaY int) bool {
	// Calculate the coordinates of the last letter in the word based on the
	// initial position, the length of the word, and the direction of placement
	// (deltaX, deltaY). The '-1' in the calculation accounts for the zero-based
	// index of the first character at the start position.
	x := start.X + (len(word)-1)*deltaX
	y := start.Y + (len(word)-1)*deltaY

	// Return the negation of isOutOfBound to check if the calculated position
	// of the last letter is within the board. isOutOfBound typically returns
	// true if the coordinates are outside the board's limits.
	return !isOutOfBound(x, y, b)
}

// isParallelPlacement checks if there are already filled cells directly
// adjacent to a given cell in the board depending on the orientation of the
// word being placed. This function is used to prevent adjacent parallel words
// from touching each other.
//
// Parameters: x, y - The coordinates of the cell to check.
//
// direction - The orientation of the word being placed (Across or Down).
//
// b - A pointer to the Board structure on which the word is being placed.
//
// Returns: true if there are filled cells adjacent to the specified cell in the
// specified direction, false otherwise.
func isParallelPlacement(x, y int, direction Direction, b *Board) bool {
	// Check directly adjacent cells depending on the word's orientation
	if direction == Across {
		// Check cell above and below each letter, but only within bounds. This
		// block handles the horizontal placement of words.
		aboveIsFilled := false
		belowIsFilled := false

		// Check above the current cell if it's not out of bounds.
		if !isOutOfBound(x, y-1, b) {
			aboveIsFilled = isCellFilled(x, y-1, b)
		}
		// Check below the current cell if it's not out of bounds.
		if !isOutOfBound(x, y+1, b) {
			belowIsFilled = isCellFilled(x, y+1, b)
		}

		// Return true if either the cell directly above or below is filled.
		return aboveIsFilled || belowIsFilled
	} else if direction == Down {
		// Check cell left and right of each letter, but only within bounds.
		// This block handles the vertical placement of words.
		leftIsFilled := false
		rightIsFilled := false

		// Check to the left of the current cell if it's not out of bounds.
		if !isOutOfBound(x-1, y, b) {
			leftIsFilled = isCellFilled(x-1, y, b)
		}
		// Check to the right of the current cell if it's not out of bounds.
		if !isOutOfBound(x+1, y, b) {
			rightIsFilled = isCellFilled(x+1, y, b)
		}

		// Return true if either the cell directly left or right is filled.
		return leftIsFilled || rightIsFilled
	}
	// If no specific direction is applicable, default to true as a safe
	// fallback.
	return true
}

func isCellFilled(x, y int, b *Board) bool {
	return b.Cells[y][x].Filled
}

// isOutOfBound checks if a cell fits on the board
func isOutOfBound(x, y int, b *Board) bool {
	return x < 0 || y < 0 || x >= len(b.Cells[0]) || y >= len(b.Cells)
}

// getDirectionDeltas determines the increments (deltaX, deltaY) for x and y
// based on the direction.
func getDirectionDeltas(direction Direction) (int, int) {
	if direction == Across {
		return 1, 0 // horizontal
	}
	return 0, 1 // vertical
}

// isCellConflict checks if there are any conflicting characters. The conflict
// tested for is: overlapping character matches the current charater or not.
//
// Returns:
//
// true if conflict
func isCellConflict(x, y int, b *Board, char string) bool {
	cell := b.Cells[y][x]
	return cell.Filled && cell.Character != char
}

// PlaceWordAt places a word on the board at a specified location and in a given
// direction. It updates the board's cells to include the new word and records
// this action in the board's history of placed words.
//
// Parameters:
//
// start - The starting location (x, y) where the first character of the word
// will be placed.
//
// word - The word to be placed on the board.
//
// direction - The direction in which the word will be placed (Across or Down).
//
// Returns:
//
// An error if the placement is invalid (not implemented here, returns nil by
// default).
func (b *Board) PlaceWordAt(start Location, word string, direction Direction) error {
	// Obtain the deltas for the direction to determine how to increment the
	// position for each character.
	deltaX, deltaY := getDirectionDeltas(direction)

	// Convert the word into a slice of runes to properly handle multi-byte
	// characters, which are common in languages that use characters beyond the
	// standard ASCII set.
	runes := []rune(word)

	// Iterate over each rune in the slice, placing each character on the board.
	for i, r := range runes {
		x := start.X + i*deltaX // Calculate the x-coordinate for the character based on its position in the word.
		y := start.Y + i*deltaY // Calculate the y-coordinate for the character based on its position in the word.

		// Place the character in the specified cell and mark the cell as
		// filled.
		b.Cells[y][x].Character = string(r)
		b.Cells[y][x].Filled = true
	}

	// Record the action of placing the word by appending a new PlacedWord
	// struct to the board's PlacedWords slice. This includes details about the
	// word's starting position, direction, and the word itself.
	b.PlacedWords = append(b.PlacedWords, PlacedWord{Start: start, Direction: direction, Word: word})

	// Update the total count of words placed on the board.
	b.WordCount = len(b.PlacedWords)

	// Return nil to indicate successful placement. In a real application, error
	// handling could be added to deal with situations where the word cannot be
	// placed (e.g., out of bounds or overlapping conflicts).
	return nil
}

// IsComplete checks if the board is fully set up with all words placed.
func (b *Board) IsComplete() bool {
	return b.WordCount >= b.TotalWords
}
