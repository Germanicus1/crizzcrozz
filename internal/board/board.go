package board

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
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
	TotalWords  int
	Pool        *Pool
	FileWriter  FileWriter `json:"-"` // Exclude from JSON serialization. Dependency injection for testing file I/O

	// Track the best solution found
	BestBoard       [][]*Cell
	BestWordCount   int
	BestPlacedWords []PlacedWord
}

type FileWriter interface {
	WriteFile(name string, data []byte, perm os.FileMode) error
}

type OSFileWriter struct{}

func (osw *OSFileWriter) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// NewBoard creates a new board with specified bounds and total words.
func NewBoard(bounds *Bounds, totalWords int, filewriter FileWriter) *Board {
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
		FileWriter: filewriter, // Ensuring FileWriter is always initialize
	}
}

// Save converts the Board struct to JSON and writes it to the file structure.
//
// Returns: Error if marshalling or file writing did not work; nil otherwise.
func (b *Board) Save() error {
	fmt.Println("Debug: Saving board data...") // REM: Debug output
	data, err := json.Marshal(b)
	if err != nil {
		panic(err) //TODO: Handle more gracefully
	}

	err = b.FileWriter.WriteFile("board.json", data, 0644)
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
//	direction - The direction to place the word (e.g., Across (1,0) or Down (0,1) ).
//
// Reports wether the word can be placed according to the rules of the game
// func (b *Board) CanPlaceWordAt(start Location, word string, direction Direction) bool {
// 	// fmt.Printf("Checking placement of '%s' at %d,%d direction: %v\n", word, start.X, start.Y, direction)
// 	deltaX, deltaY := getDirectionDeltas(direction)
// 	intersected := false

// 	// word is a string which is a []byte. A unicode character can occupy 1-4
// 	// bytes. We need to make sure that each letter occupies only 1 cell.
// 	lettersCountInWord := len([]rune(word))

// 	// Step 1: Check if word fits on the board atr a certain position.
// 	if !b.isPlacementWithinBounds(start, lettersCountInWord, deltaX, deltaY) {
// 		return false
// 	}

// 	// Step 2: Check if letters can be placed. A word needs to have at leas one
// 	// intersection, but can have more if they are NOT consecutive.

// 	intersected = b.canPlaceLetters(start, word, deltaX, deltaY)

// 	// Step 3: Check if cells immediately before and after the word are empty to
// 	// prevent contiguous word formation.
// 	// placed at the boarder (0,0)
// 	if !b.isPlacementIsolated(start, lettersCountInWord, deltaX, deltaY) {
// 		return false
// 	}

//		return intersected
//	}
//
// REM debugging
func (b *Board) CanPlaceWordAt(start Location, word string, direction Direction) bool {
	deltaX, deltaY := getDirectionDeltas(direction)
	intersected := false

	runes := []rune(word)

	// Step 1: Ensure word fits within board bounds
	if !b.isPlacementWithinBounds(start, len(runes), deltaX, deltaY) {
		return false
	}

	// Step 2: Check each letter's placement validity
	// for i, char := range runes {
	// 	x := start.X + i*deltaX
	// 	y := start.Y + i*deltaY

	// 	// 🚨 Prevent placement if the cell is locked
	// 	if b.Cells[y][x].LockCount > 0 {
	// 		fmt.Printf("❌ ERROR: Cannot place word '%s' at (%d, %d), cell is locked!\n", word, x, y)
	// 		return false
	// 	}

	// 	// If the cell is already occupied, it must match the letter
	// 	if b.Cells[y][x].Filled {
	// 		if b.Cells[y][x].Character != string(char) {
	// 			fmt.Printf("❌ ERROR: Overwriting existing letter at (%d, %d): '%s' with '%s'\n",
	// 				x, y, b.Cells[y][x].Character, string(char))
	// 			return false
	// 		}
	// 		intersected = true
	// 	}
	// }
	intersected = b.canPlaceLetters(start, word, deltaX, deltaY)

	// 🚨 NEW: Ensure the cells before and after are not locked
	if !b.isPlacementIsolated(start, len(runes), deltaX, deltaY) {
		fmt.Printf("❌ ERROR: Word '%s' at (%d, %d) is not isolated\n", word, start.X, start.Y)
		return false
	}

	// Step 3: Ensure word actually intersects at least once
	return intersected
}

//	isPlacementIsolated checks weather the cell before or after the word is
//	empyt. This only if thew cell is not out of bound.
//
// start - a Location structx with X, Y integers
//
// lettersCountInWord - The number of letters (not bytes) in the word
//
// deltaX/deltaY - Returned  by getDirection(direction). (0,1) for down and
// (1,0) for across.
//
// Returns: weather it's an isolated placement

// func (b *Board) isPlacementIsolated(start Location, lettersInWord, deltaX, deltaY int) bool {
// 	// Check the cell before the word
// 	xBefore, yBefore := start.X-deltaX, start.Y-deltaY
// 	if !isOutOfBound(xBefore, yBefore, b) && isCellFilled(xBefore, yBefore, b) {
// 		return false
// 	}

// 	// Check the cell after the word
// 	xAfter, yAfter := start.X+lettersInWord*deltaX, start.Y+lettersInWord*deltaY
// 	if !isOutOfBound(xAfter, yAfter, b) && isCellFilled(xAfter, yAfter, b) {
// 		return false
// 	}

//		return true
//	}
//
// REM debugging
func (b *Board) isPlacementIsolated(start Location, lettersInWord, deltaX, deltaY int) bool {
	isIsolated := true

	// Check the cell before the word
	xBefore, yBefore := start.X-deltaX, start.Y-deltaY
	if !isOutOfBound(xBefore, yBefore, b) && isCellFilled(xBefore, yBefore, b) {
		fmt.Printf("❌ ERROR: Cell BEFORE word at (%d, %d) is occupied!\n", xBefore, yBefore)
		isIsolated = false
	}

	// Check the cell after the word
	xAfter, yAfter := start.X+lettersInWord*deltaX, start.Y+lettersInWord*deltaY
	if !isOutOfBound(xAfter, yAfter, b) && isCellFilled(xAfter, yAfter, b) {
		fmt.Printf("❌ ERROR: Cell AFTER word at (%d, %d) is occupied!\n", xAfter, yAfter)
		isIsolated = false
	}

	if isIsolated {
		fmt.Printf("✅ Word at (%d, %d) is properly isolated.\n", start.X, start.Y)
	} else {
		fmt.Printf("❌ Word at (%d, %d) is NOT isolated! Placement should be blocked.\n", start.X, start.Y)
	}

	return isIsolated
}

// New function to handle character placement checks
func (b *Board) canPlaceLetters(start Location, word string, deltaX, deltaY int) bool {
	intersectedWord := false
	intersectionCount := 0
	runes := []rune(word)
	consecutiveIntersections := 0 // Track consecutive intersections

	// if word == "wieder" {
	// 	fmt.Println(word)
	// }
	// if word == "kaputt" {
	// 	fmt.Println(word)
	// }

	for i := 0; i < len(runes); i++ {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY
		cellIsIntersection := false

		// Check if it intersects correctly with an existing letter
		if b.isValidIntersection(x, y, string(runes[i])) {
			intersectedWord = true
			cellIsIntersection = true
			intersectionCount++
			consecutiveIntersections++
		} else {
			consecutiveIntersections = 0 // Reset count if not consecutive
		}

		// If there are two or more consecutive intersections, placement is invalid
		if consecutiveIntersections >= 2 {
			return false
		}

		// Check for invalid adjacent placements (parallel words)
		if !cellIsIntersection && isParallelPlacement(x, y, getDirectionFromDeltas(deltaX, deltaY), b) {
			return false
		}
	}

	return intersectedWord
}

// Check weather a cell is filled and if the character is the same as the
// letter we are trying to fill it with.
func (b *Board) isValidIntersection(x, y int, char string) bool {
	// if b.Cells[y][x].Filled {
	// 	return b.Cells[y][x].Character == char
	// }
	// return true
	return b.Cells[y][x].Filled && b.Cells[y][x].Character == char
}

// Helper function to determine the direction from deltas
func getDirectionFromDeltas(deltaX, deltaY int) Direction {
	if deltaX == 1 && deltaY == 0 {
		return Across
	}
	return Down
}

// isPlacementWithinBounds checks if the placement of the last character of the
// word at the calculated position is within the boundaries of the board. It
// calculates the position based on the starting location, the length of the
// word, and the direction in which the word is placed. Parameters:
//
//	start- The starting location (x, y) for placing the word.
//	word - The word to be placed.
//	deltaX  - The horizontal direction increment (1 for rightward, -1 for leftward, 0 for none.)
//	deltaY  - The vertical direction increment (1 for downward, -1 for upward, 0 for none).
//
// Reports wether the last character of the word fits within the board boundaries.
func (b *Board) isPlacementWithinBounds(start Location, charsCount int, deltaX, deltaY int) bool {
	// Calculate the coordinates of the last letter in the word based on the
	// initial position, the length of the word, and the direction of placement
	// (deltaX, deltaY). The '-1' in the calculation accounts for the zero-based
	// index of the first character at the start position.
	x := start.X + (charsCount-1)*deltaX
	y := start.Y + (charsCount-1)*deltaY

	// Return the negation of isOutOfBound to check if the calculated position
	// of the last letter is within the board. isOutOfBound typically returns
	// true if the coordinates are outside the board's limits.
	return !isOutOfBound(x, y, b)
}

// isParallelPlacement checks if there are already filled cells directly
// adjacent to a given cell in the board depending on the orientation of the
// word being placed. This function is used to prevent adjacent parallel words
// from touching each other. Parameters:
//
// x, y - The coordinates of the cell to check.
//
// direction - The orientation of the word being placed (Across or Down).
//
// b - A pointer to the Board structure on which the word is being placed.
//
// Reports wether there are filled cells adjacent to the specified cell in the
// specified direction.
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
	}
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

// PlaceWordAt places a word on the board at a specified location and in a given
// direction. It updates the board's cells to include the new word and records
// this action in the board's history of placed words. Parameters:
//
// start - The starting location (x, y) where the first character of the word
// will be placed.
//
// word - The word to be placed on the board.
//
// direction - The direction in which the word will be placed (Across or Down).
//
// Returns: An error if the placement is invalid (not implemented here, returns
// nil by default).
// func (b *Board) PlaceWordAt(start Location, word string, direction Direction) error {
// 	// Obtain the deltas for the direction to determine how to increment the
// 	// position for each character.
// 	// REM debugging
// 	fmt.Printf("📌 Attempting to place: %s at (%d, %d) %v\n", word, start.X, start.Y, direction)

// 	deltaX, deltaY := getDirectionDeltas(direction)

// 	// Convert the word into a slice of runes to properly handle multi-byte
// 	// characters, which are common in languages that use characters beyond the
// 	// standard ASCII set.
// 	runes := []rune(word)

// 	for i, r := range runes {
// 		x := start.X + i*deltaX
// 		y := start.Y + i*deltaY
// 		cell := b.Cells[y][x]
// 		cell.Character = string(r)
// 		cell.Filled = true
// 		cell.UsageCount++

// 		// REM debugging
// 		fmt.Printf("  ✅ Placed '%s' at (%d, %d), UsageCount: %d\n", string(r), x, y, cell.UsageCount)

// 	}

// 	b.PlacedWords = append(b.PlacedWords, PlacedWord{Start: start, Direction: direction, Word: word})
// 	b.WordCount++

// 	//TODO: error handling
// 	// REM debugging
// 	fmt.Printf("✅ Successfully placed: %s at (%d, %d) %v\n", word, start.X, start.Y, direction)

//		return nil
//	}
//
// REM debugging
func (b *Board) PlaceWordAt(start Location, word string, direction Direction) error {
	fmt.Printf("📌 Placing word: %s at (%d, %d) %v\n", word, start.X, start.Y, direction)

	deltaX, deltaY := getDirectionDeltas(direction)
	runes := []rune(word)

	for i, r := range runes {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY

		cell := b.Cells[y][x]
		cell.Character = string(r)
		cell.Filled = true
		cell.UsageCount++

		fmt.Printf("  ✅ Placed '%s' at (%d, %d), UsageCount: %d\n", string(r), x, y, cell.UsageCount)
	}

	// 🚨 NEW: Lock cells before and after the word
	b.lockAdjacentCells(start, len(runes), deltaX, deltaY)

	b.PlacedWords = append(b.PlacedWords, PlacedWord{Start: start, Direction: direction, Word: word})
	b.WordCount++
	return nil
}

// DEPRECIATED: IsComplete checks if the board is fully set up with all words
// placed.
func (b *Board) IsComplete() bool {
	return b.WordCount >= b.TotalWords
}

// Assuming each cell knows which word it belongs to (you might need to adjust your data structures)
// func (b *Board) RemoveWord(start Location, word string, direction Direction) {
// 	deltaX, deltaY := getDirectionDeltas(direction)
// 	runes := []rune(word)
// 	for i := range runes {
// 		x := start.X + i*deltaX
// 		y := start.Y + i*deltaY
// 		cell := b.Cells[y][x]
// 		// Only decrement usage if this word was actually using the cell
// 		if cell.Filled && cell.Character == string(runes[i]) {
// 			cell.UsageCount--
// 		}

// 		if cell.UsageCount == 0 {
// 			cell.Character = "" // Clear the character only if no other word is using this cell
// 			cell.Filled = false
// 		}
// 		// REM debug info
// 		if x == 14 && y == 13 && (word == "wieder" || b.Cells[14][13].Filled == true) {
// 			fmt.Printf("CELL (14,13): %v, %v\n", b.Cells[14][13].Filled, b.Cells[14][13].Character)
// 		}

// 	}
// 	// Remove the word from PlacedWords and update WordCount
// 	for index, placed := range b.PlacedWords {
// 		if placed.Start == start && placed.Word == word && placed.Direction == direction {
// 			b.PlacedWords = append(b.PlacedWords[:index], b.PlacedWords[index+1:]...)
// 			break
// 		}
// 	}
// 	b.WordCount--
// }

// REM debugging
func (b *Board) RemoveWord(start Location, word string, direction Direction) {
	fmt.Printf("🔄 Backtracking: Removing word %s from (%d, %d) %v\n", word, start.X, start.Y, direction)

	deltaX, deltaY := getDirectionDeltas(direction)
	runes := []rune(word)

	for i := range runes {
		x := start.X + i*deltaX
		y := start.Y + i*deltaY
		cell := b.Cells[y][x]

		cell.UsageCount--
		if cell.UsageCount == 0 {
			cell.Character = ""
			cell.Filled = false
		}
	}

	// 🚨 NEW: Unlock cells before and after the word
	b.unlockAdjacentCells(start, len(runes), deltaX, deltaY)

	// Remove word from placed words list
	for index, placed := range b.PlacedWords {
		if placed.Start == start && placed.Word == word && placed.Direction == direction {
			b.PlacedWords = append(b.PlacedWords[:index], b.PlacedWords[index+1:]...)
			break
		}
	}

	b.WordCount--
}

func (b *Board) SaveBestSolution() {
	b.BestWordCount = b.WordCount

	// Deep copy the board cells
	b.BestBoard = make([][]*Cell, len(b.Cells))
	for i := range b.Cells {
		b.BestBoard[i] = make([]*Cell, len(b.Cells[i]))
		for j := range b.Cells[i] {
			cellCopy := *b.Cells[i][j] // Copy struct
			b.BestBoard[i][j] = &cellCopy
		}
	}

	// Copy the placed words
	b.BestPlacedWords = make([]PlacedWord, len(b.PlacedWords))
	copy(b.BestPlacedWords, b.PlacedWords) //Directly assiging b.Cells, later changes will affect the stored best board.
}

// PrintBestSolution outputs the crossword board to the console. It marks filled
// cells with their respective characters and empty cells with a dot.
// func (b *Board) PrintBestSolution() {
// 	if b.BestBoard == nil {
// 		fmt.Println("No valid crossword solution found.")
// 		return
// 	}

//		// fmt.Println("\nBest Solution Found:")
//		for y := 0; y < len(b.BestBoard); y++ {
//			for x := 0; x < len(b.BestBoard[0]); x++ {
//				cell := b.BestBoard[y][x]
//				if cell != nil && cell.Filled {
//					fmt.Print(strings.ToUpper(cell.Character), " ")
//				} else {
//					fmt.Print(". ")
//				}
//			}
//			fmt.Println()
//		}
//	}
//
// REM debug version
func (b *Board) PrintBestSolution() {
	if b.BestBoard == nil {
		fmt.Println("❌ No valid crossword solution found.")
		return
	}

	fmt.Println("\n🔎 Final Board Rendering:")
	for y := 0; y < len(b.BestBoard); y++ {
		for x := 0; x < len(b.BestBoard[0]); x++ {
			cell := b.BestBoard[y][x]
			if cell != nil && cell.Filled {
				fmt.Print(strings.ToUpper(cell.Character), " ")
			} else {
				if cell.LockCount > 0 {
					fmt.Print("☐ ")
				} else {
					fmt.Print("・")
				}
			}
		}
		fmt.Println()
	}

	fmt.Println("\n📌 Words officially placed (from BestPlacedWords):")
	for _, placed := range b.BestPlacedWords {
		fmt.Printf("  → %s at (%d, %d) %v\n", placed.Word, placed.Start.X, placed.Start.Y, placed.Direction)
	}

	fmt.Printf("\n✅ Words placed: %d / %d\n", len(b.BestPlacedWords), b.TotalWords)
	if len(b.BestPlacedWords) < b.TotalWords {
		fmt.Println("⚠️ WARNING: Not all words were placed correctly!")
	}
	fmt.Println("=======================================")
}

func (b *Board) lockAdjacentCells(start Location, lettersInWord, deltaX, deltaY int) {
	// Lock the cell before the word
	xBefore, yBefore := start.X-deltaX, start.Y-deltaY
	if !isOutOfBound(xBefore, yBefore, b) {
		b.Cells[yBefore][xBefore].LockCount++
		fmt.Printf("🔒 Locking cell before word at (%d, %d), LockCount: %d\n", xBefore, yBefore, b.Cells[yBefore][xBefore].LockCount)
	}

	// Lock the cell after the word
	xAfter, yAfter := start.X+lettersInWord*deltaX, start.Y+lettersInWord*deltaY
	if !isOutOfBound(xAfter, yAfter, b) {
		b.Cells[yAfter][xAfter].LockCount++
		fmt.Printf("🔒 Locking cell after word at (%d, %d), LockCount: %d\n", xAfter, yAfter, b.Cells[yAfter][xAfter].LockCount)
	}
}

func (b *Board) unlockAdjacentCells(start Location, lettersInWord, deltaX, deltaY int) {
	// Unlock the cell before the word
	xBefore, yBefore := start.X-deltaX, start.Y-deltaY
	if !isOutOfBound(xBefore, yBefore, b) {
		b.Cells[yBefore][xBefore].LockCount--
		fmt.Printf("🔓 Unlocking cell before word at (%d, %d), LockCount: %d\n", xBefore, yBefore, b.Cells[yBefore][xBefore].LockCount)
	}

	// Unlock the cell after the word
	xAfter, yAfter := start.X+lettersInWord*deltaX, start.Y+lettersInWord*deltaY
	if !isOutOfBound(xAfter, yAfter, b) {
		b.Cells[yAfter][xAfter].LockCount--
		fmt.Printf("🔓 Unlocking cell after word at (%d, %d), LockCount: %d\n", xAfter, yAfter, b.Cells[yAfter][xAfter].LockCount)
	}
}
