package models

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

// TODO-mIHdG3: Refactor CanPlaceWordAt() REFACTOR

func (b *Board) CanPlaceWordAt(start Location, word string, direction Direction) bool {
	deltaX, deltaY := getDirectionDeltas(direction) // 1,0 = Across 0,1 = Down
	intersected := false                            // To check if at least one letter overlaps with existing words
	intersectionCount := 0

	// Calculate the position of the last letter in a word depending on
	// the start and the direction.
	x := start.X + (len(word)-1)*deltaX
	y := start.Y + (len(word)-1)*deltaY

	// Check if the word fits on the board
	if isOutOfBound(x, y, b) {
		return false
	}

	for i := 0; i < len(word); i++ {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY
		isIntersection := false

		// Check if the intersection has different letters
		if isCellConflict(x, y, b, string(word[i])) {
			return false
		}

		// Ensure the new word intersects once but not more (avoid
		// writing over an existing word)
		if b.Cells[y][x].Filled && b.Cells[y][x].Character == string(word[i]) {
			intersected = true
			isIntersection = intersected
			intersectionCount++
		}
		if intersectionCount > 1 {
			return false
		}

		// Check parallel positions to prevent adjacent parallel words.
		// Only on letters which are not valid intersections
		if !isIntersection {
			if isParallelPlacement(x, y, direction, b) {
				return false
			}
		}

	}

	// Are cells before or after occupied? Calculate the position of the
	// last letter
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

// Does the word fit on the bnopard?
func (b *Board) DoesWordFitOnBoard(start Location, word string, direction Direction) bool {
	deltaX, deltaY := getDirectionDeltas(direction) // 1,0 = Across 0,1 = Down

	// Calculate the position of the last letter in a word depending on
	// the start and the direction.
	x := start.X + (len(word)-1)*deltaX
	y := start.Y + (len(word)-1)*deltaY

	// Check if the last character fits on the board
	return !isOutOfBound(x, y, b)

	// Ensure the word intersects with existing words on the board
	return true
}

// if isOutOfBound(x, y-1, b) is true and direction when Across, we are
// trying to place a letter at the top of the board.

func isParallelPlacement(x, y int, direction Direction, b *Board) bool {
	// Check directly adjacent cells depending on the word's orientation
	if direction == Across {
		// Check cell above and below each letter, but only within
		// bounds
		aboveIsFilled := false
		belowIsFilled := false

		if !isOutOfBound(x, y-1, b) { // Check above only if it's not out of bounds
			aboveIsFilled = isCellFilled(x, y-1, b)
		}
		if !isOutOfBound(x, y+1, b) { // Check below only if it's not out of bounds
			belowIsFilled = isCellFilled(x, y+1, b)
		}

		return aboveIsFilled || belowIsFilled
	} else if direction == Down {
		// Check cell left and right of each letter, but only within
		// bounds
		leftIsFilled := false
		rightIsFilled := false

		if !isOutOfBound(x-1, y, b) { // Check left only if it's not out of bounds
			leftIsFilled = isCellFilled(x-1, y, b)
		}
		if !isOutOfBound(x+1, y, b) { // Check right only if it's not out of bounds
			rightIsFilled = isCellFilled(x+1, y, b)
		}

		return leftIsFilled || rightIsFilled
	}
	return true
}

func isCellFilled(x, y int, b *Board) bool {
	return b.Cells[y][x].Filled
}

// isOutOfBound checks if a cell fits on the board
func isOutOfBound(x, y int, b *Board) bool {
	// return x >= len(b.Cells[0]) || x < 0 || y >= len(b.Cells) || y <
	// 0
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
func isCellConflict(x, y int, b *Board, char string) bool {
	cell := b.Cells[y][x]
	return cell.Filled && cell.Character != char
}

func (b *Board) PlaceWordAt(start Location, word string, direction Direction) error {
	deltaX, deltaY := getDirectionDeltas(direction)

	// Convert the word to a slice of runes to handle multi-byte characters properly
	runes := []rune(word)

	// Place the word
	for i, r := range runes {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY
		b.Cells[y][x].Character = string(r)
		b.Cells[y][x].Filled = true
	}

	// Record the placed word
	b.PlacedWords = append(b.PlacedWords, PlacedWord{Start: start, Direction: direction, Word: word})
	b.WordCount = len(b.PlacedWords)

	// REM: fmt.Printf("Word '%s' placed at (%d, %d) going %s\n", word,
	// start.X, start.Y, directionString(direction))

	return nil
}

// func (b *Board) isValidWord(word string) bool {
// 	_, exists := b.WordList[word]
// 	return exists
// }

// func (b *Board) checkPerpendicularIntersection(newStart Location, word string, newDirection Direction) bool {
// 	// If there are no placed words, always allow the placement (needed
// 	// for the first word)
// 	if len(b.PlacedWords) == 0 {
// 		return true
// 	}
// 	deltaX, deltaY := getDirectionDeltas(newDirection)
// 	newEndX := newStart.X + deltaX*(len(word)-1)
// 	newEndY := newStart.Y + deltaY*(len(word)-1)

// 	for _, placed := range b.PlacedWords {
// 		placedDeltaX, placedDeltaY := getDirectionDeltas(placed.Direction)
// 		placedEndX := placed.Start.X + placedDeltaX*(len(placed.Word)-1)
// 		placedEndY := placed.Start.Y + placedDeltaY*(len(placed.Word)-1)

// 		// Check if directions are perpendicular
// 		if (deltaX == 0 && placedDeltaX != 0) || (deltaX != 0 && placedDeltaX == 0) {
// 			// Determine if they intersect by checking coordinate ranges
// 			if newStart.X <= placedEndX && newEndX >= placed.Start.X &&
// 				newStart.Y <= placedEndY && newEndY >= placed.Start.Y {
// 				// They intersect and are perpendicular
// 				return true
// 			}
// 		}
// 	}

// 	// No valid perpendicular intersections found
// 	return false
// }

// directionString returns a string representation of a Direction.
// func directionString(direction Direction) string {
// 	switch direction {
// 	case Across:
// 		return "Across"
// 	case Down:
// 		return "Down"
// 	default:
// 		return "Unknown Direction"
// 	}
// }

// Consolidated method to check word validity from a location
// considering all potential word formations
// func (b *Board) isPartOfValidWord(x, y, deltaX, deltaY int) bool {
// 	// Check if a word formed starting at this cell in both directions
// 	// is valid
// 	if checkWordFormed(x, y, deltaX, deltaY, b) || checkWordFormed(x, y, -deltaX, -deltaY, b) {
// 		return true
// 	}
// 	return false
// }

// Helper function to generate a word from a start point in a given
// direction and check its validity
// func checkWordFormed(x, y, deltaX, deltaY int, b *Board) bool {
// 	var word []rune
// 	// Start at the given point and move in the specified direction
// 	for !isOutOfBound(x, y, b) && b.Cells[y][x].Filled {
// 		word = append(word, b.Cells[y][x].Character)
// 		x += deltaX
// 		y += deltaY
// 	}
// 	// Check if the formed word is in the list of valid words
// 	return b.isValidWord(string(word))
// }

// func (b *Board) checkWordValidity(x, y, deltaX, deltaY int) bool {
// 	word := b.extractWord(x, y, deltaX, deltaY)
// 	_, exists := b.WordList[word] // Assuming b.ValidWords is a map containing valid words
// 	return exists
// }

// func (b *Board) extractWord(x, y, deltaX, deltaY int) string {
// 	var word []rune

// 	// Extend backward to the start of the word
// 	for x >= 0 && x < len(b.Cells[0]) && y >= 0 && y < len(b.Cells) && b.Cells[y][x].Filled {
// 		x -= deltaX
// 		y -= deltaY
// 	}
// 	// Move forward one step to the actual start of the word
// 	x += deltaX
// 	y += deltaY

// 	// Now extend forward to extract the whole word
// 	for x >= 0 && x < len(b.Cells[0]) && y >= 0 && y < len(b.Cells) && b.Cells[y][x].Filled {
// 		word = append(word, b.Cells[y][x].Character)
// 		x += deltaX
// 		y += deltaY
// 	}

// 	return string(word)
// }

// func (b *Board) isContinuingWord(x, y, deltaX, deltaY int) bool {
// 	// Check in the placement direction from the current cell
// 	nextX := x + deltaX
// 	nextY := y + deltaY
// 	prevX := x - deltaX
// 	prevY := y - deltaY

// 	return (!isOutOfBound(nextX, nextY, b) && b.Cells[nextY][nextX].Filled) &&
// 		(!isOutOfBound(prevX, prevY, b) && b.Cells[prevY][prevX].Filled)
// }

// func (b *Board) isValidIntersection(x, y, deltaX, deltaY int) bool {
// 	// Assumes that a valid intersection must not extend the same word
// 	// in both perpendicular directions
// 	nextX := x + deltaX
// 	nextY := y + deltaY
// 	prevX := x - deltaX
// 	prevY := y - deltaY

// 	// Check if the cell is a continuation of a word in the placement
// 	// direction or standalone
// 	if isOutOfBound(nextX, nextY, b) || isOutOfBound(prevX, prevY, b) || (!b.Cells[nextY][nextX].Filled && !b.Cells[prevY][prevX].Filled) {
// 		return true // Valid if it's not extending in the same direction
// 	}
// 	return false
// }

// generateWordFromLocation generates a word starting from a given
// location in a specified direction
// func (b *Board) generateWordFromLocation(start Location, deltaX, deltaY int) string {
// 	var word []rune

// 	// Move backwards to the start of the word
// 	x, y := start.X, start.Y
// 	for b.isValidLocation(x-deltaX, y-deltaY) && b.Cells[y-deltaY][x-deltaX].Filled {
// 		x -= deltaX
// 		y -= deltaY
// 	}

// 	// Generate the word forward
// 	for b.isValidLocation(x, y) && b.Cells[y][x].Filled {
// 		x += deltaX
// 		y += deltaY
// 		word = append(word, b.Cells[y][x].Character)
// 	}

// 	return string(word)
// }

// isValidLocation checks if the given coordinates are within the bounds
// of the board
// func (b *Board) isValidLocation(x, y int) bool {
// 	return x >= 0 && y >= 0 && x < len(b.Cells[0]) && y < len(b.Cells)
// }

// IsComplete checks if the board is fully set up with all words placed.
func (b *Board) IsComplete() bool {
	return b.WordCount >= b.TotalWords
}
