package models

import (
	"fmt"
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

// CanPlaceWordAt checks if a word can be placed at a specific location
// in a given direction.
func (b *Board) CanPlaceWordAt(start Location, word string, direction Direction) bool {
	deltaX, deltaY := getDirectionDeltas(direction) // 1,0 = Across 0,1 = Down
	intersected := false                            // To check if at least one letter overlaps with existing words
	intersectionCount := 0

	// REM: CanPlaceWordAt, debug info
	// fmt.Printf("CanPlaceWordAt: Testing word %s at start location %+v, in direction %v\n", word, start, direction)

	// Calculate the position of the last letter in a word depending on the start and the direction.
	x := start.X + (len(word)-1)*deltaX
	y := start.Y + (len(word)-1)*deltaY

	// Check if the word fits on the board
	if isOutOfBound(x, y, b) {
		// // REM: isOutOfBound, debug info
		// fmt.Printf("%v with %v char too long starting at %v going %v\n", word, len(word), start, directionString(direction))
		return false
	}

	// REFACTOR: Only check possible placements.
	for i := 0; i < len(word); i++ {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY

		// Check if the intersection has different letters
		if isCellConflict(x, y, b, rune(word[i])) {
			return false
		}

		// Ensure the new word intersects once but not more (avoid
		// writing over an existing word)
		if b.Cells[y][x].Filled && b.Cells[y][x].Character == rune(word[i]) {
			intersected = true
			intersectionCount = +1
		}
		if intersectionCount > 1 {
			return false
		}

		// Check parallel positions to prevent adjacent parallel words
		// if isParallelPlacement(x, y, deltaX, deltaY, b) {
		// 	return false
		// }

	}

	// // REM: Remove debug info
	// // fmt.Printf("CanPlaceWordAt.intersected: %v\n", intersected)

	// TODO: Cells before or after occupied?
	// Calculate the position of the last letter
	xEnd := start.X + (len(word)-1)*deltaX
	yEnd := start.Y + (len(word)-1)*deltaY

	xBefore, yBefore := 0, 0

	if start.X > 0 {
		xBefore = start.X - deltaX
	}

	if start.Y > 0 {
		yBefore = start.Y - deltaY
	}

	// check for bound violation
	xAfter := xEnd + deltaX
	yAfter := yEnd + deltaY

	if isOutOfBound(xBefore, yBefore, b) || isCellFilled(xBefore, yBefore, b) || (isOutOfBound(xAfter, yAfter, b) || isCellFilled(xAfter, yAfter, b)) {
		return false
	}

	// Ensure the word intersects with existing words on the board
	return intersected
}

func isParallelPlacement(x, y int, deltaX, deltaY int, b *Board) bool {
	// Check directly adjacent cells depending on the word's orientation
	if deltaX == 1 && deltaY == 0 { // Horizontal placement
		if isOutOfBound(x, y-1, b) || isOutOfBound(x, y+1, b) {
			return false
		}
		return !(isCellFilled(x, y-1, b) || isCellFilled(x, y+1, b))
	} else if deltaX == 0 && deltaY == 1 { // Vertical placement
		if isOutOfBound(x-1, y, b) || isOutOfBound(x+1, y, b) {
			return false
		}
		return !(isCellFilled(x-1, y, b) || isCellFilled(x+1, y, b))
	}
	return true
}

// // Helper function to check if a cell is out of bounds or filled
// func isCellOutOfBoundsOrFilled(x, y int, b *Board) bool {
// 	return x < 0 || x >= len(b.Cells[0]) || y < 0 || y >= len(b.Cells) || b.Cells[y][x].Filled
// }

func isCellFilled(x, y int, b *Board) bool {
	return b.Cells[y][x].Filled
}

// isOutOfBound checks if a cell fits on the board
func isOutOfBound(x, y int, b *Board) bool {
	// return x >= len(b.Cells[0]) || x < 0 || y >= len(b.Cells) || y < 0
	return x < 0 || y < 0 || x >= len(b.Cells[0]) || y >= len(b.Cells)
}

// getDirectionDeltas determines the increments (deltaX, deltaY) for x
// and y based on the direction.
func getDirectionDeltas(direction Direction) (int, int) {
	if direction == Across {
		return 1, 0 // horizontal
	}
	return 0, 1 // vertical
}

// isCellConflict checks if there are any conflicting characters. The
// conflict tested for is: overlapping character matches the current
// charater or not. Returns true if conflict
func isCellConflict(x, y int, b *Board, char rune) bool {
	cell := b.Cells[y][x]
	return cell.Filled && cell.Character != char
}

func (b *Board) PlaceWordAt(start Location, word string, direction Direction) error {
	deltaX, deltaY := getDirectionDeltas(direction)
	// var completeWord []rune

	// Attempt to place the word
	// for i := 0; i < len(word); i++ {
	// 	x := start.X + i*deltaX
	// 	y := start.Y + i*deltaY

	// 	if isOutOfBound(x, y, b) {
	// 		fmt.Printf("Out of bounds at (%d, %d)\n", x, y)
	// 		return fmt.Errorf("out of bounds at (%d, %d)", x, y)
	// 	}
	// 	if isCellConflict(x, y, b, rune(word[i])) {
	// 		fmt.Printf("Cell conflict at (%d, %d)\n", x, y)
	// 		return fmt.Errorf("cell conflict at (%d, %d)", x, y)
	// 	}
	// 	if !b.checkPerpendicularIntersection(start, word, direction) {
	// 		fmt.Printf("Invalid intersection at (%d, %d)\n", start.X, start.Y)
	// 		return fmt.Errorf("invalid intersection at (%d, %d)", start.X, start.Y)
	// 	}
	// 	// // Check for valid parallel placements
	// 	// if !b.isValidParallelPlacement(x, y, deltaX, deltaY, word) {
	// 	// 	return fmt.Errorf("invalid parallel word formation at position (%d, %d)", x, y)
	// 	// }
	// 	completeWord = append(completeWord, rune(word[i])) // Build the word being placed

	// }

	// Place the word
	for i := 0; i < len(word); i++ {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY
		b.Cells[y][x].Character = rune(word[i])
		b.Cells[y][x].Filled = true
	}

	// Record the placed word
	b.PlacedWords = append(b.PlacedWords, PlacedWord{Start: start, Direction: direction, Word: word})
	b.WordCount = len(b.PlacedWords)
	fmt.Printf("Word '%s' placed at (%d, %d) going %s\n", word, start.X, start.Y, directionString(direction))

	return nil
}

func (b *Board) isValidWord(word string) bool {
	_, exists := b.WordList[word]
	return exists
}

func (b *Board) checkPerpendicularIntersection(newStart Location, word string, newDirection Direction) bool {
	// If there are no placed words, always allow the placement (needed for the first word)
	if len(b.PlacedWords) == 0 {
		return true
	}
	deltaX, deltaY := getDirectionDeltas(newDirection)
	newEndX := newStart.X + deltaX*(len(word)-1)
	newEndY := newStart.Y + deltaY*(len(word)-1)

	for _, placed := range b.PlacedWords {
		placedDeltaX, placedDeltaY := getDirectionDeltas(placed.Direction)
		placedEndX := placed.Start.X + placedDeltaX*(len(placed.Word)-1)
		placedEndY := placed.Start.Y + placedDeltaY*(len(placed.Word)-1)

		// Check if directions are perpendicular
		if (deltaX == 0 && placedDeltaX != 0) || (deltaX != 0 && placedDeltaX == 0) {
			// Determine if they intersect by checking coordinate ranges
			if newStart.X <= placedEndX && newEndX >= placed.Start.X &&
				newStart.Y <= placedEndY && newEndY >= placed.Start.Y {
				// They intersect and are perpendicular
				return true
			}
		}
	}

	// No valid perpendicular intersections found
	return false
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

// Consolidated method to check word validity from a location considering all potential word formations
func (b *Board) isPartOfValidWord(x, y, deltaX, deltaY int) bool {
	// Check if a word formed starting at this cell in both directions is valid
	if checkWordFormed(x, y, deltaX, deltaY, b) || checkWordFormed(x, y, -deltaX, -deltaY, b) {
		return true
	}
	return false
}

// Helper function to generate a word from a start point in a given direction and check its validity
func checkWordFormed(x, y, deltaX, deltaY int, b *Board) bool {
	var word []rune
	// Start at the given point and move in the specified direction
	for !isOutOfBound(x, y, b) && b.Cells[y][x].Filled {
		word = append(word, b.Cells[y][x].Character)
		x += deltaX
		y += deltaY
	}
	// Check if the formed word is in the list of valid words
	return b.isValidWord(string(word))
}

func (b *Board) checkWordValidity(x, y, deltaX, deltaY int) bool {
	word := b.extractWord(x, y, deltaX, deltaY)
	_, exists := b.WordList[word] // Assuming b.ValidWords is a map containing valid words
	return exists
}

func (b *Board) extractWord(x, y, deltaX, deltaY int) string {
	var word []rune

	// Extend backward to the start of the word
	for x >= 0 && x < len(b.Cells[0]) && y >= 0 && y < len(b.Cells) && b.Cells[y][x].Filled {
		x -= deltaX
		y -= deltaY
	}
	// Move forward one step to the actual start of the word
	x += deltaX
	y += deltaY

	// Now extend forward to extract the whole word
	for x >= 0 && x < len(b.Cells[0]) && y >= 0 && y < len(b.Cells) && b.Cells[y][x].Filled {
		word = append(word, b.Cells[y][x].Character)
		x += deltaX
		y += deltaY
	}

	return string(word)
}

func (b *Board) isContinuingWord(x, y, deltaX, deltaY int) bool {
	// Check in the placement direction from the current cell
	nextX := x + deltaX
	nextY := y + deltaY
	prevX := x - deltaX
	prevY := y - deltaY

	return (!isOutOfBound(nextX, nextY, b) && b.Cells[nextY][nextX].Filled) &&
		(!isOutOfBound(prevX, prevY, b) && b.Cells[prevY][prevX].Filled)
}

func (b *Board) isValidIntersection(x, y, deltaX, deltaY int) bool {
	// Assumes that a valid intersection must not extend the same word in both perpendicular directions
	nextX := x + deltaX
	nextY := y + deltaY
	prevX := x - deltaX
	prevY := y - deltaY

	// Check if the cell is a continuation of a word in the placement direction or standalone
	if isOutOfBound(nextX, nextY, b) || isOutOfBound(prevX, prevY, b) || (!b.Cells[nextY][nextX].Filled && !b.Cells[prevY][prevX].Filled) {
		return true // Valid if it's not extending in the same direction
	}
	return false
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

// // isValidWord checks if the word is in the allowed word list
// func (b *Board) isValidWord(word string, pool *Pool) bool {
// 	return pool.Exists(word)
// }

// IsComplete checks if the board is fully set up with all words placed.
func (b *Board) IsComplete() bool {
	return b.WordCount >= b.TotalWords
}
