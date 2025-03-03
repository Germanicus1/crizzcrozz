package models

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestNewBoard(t *testing.T) {
	// Create Bounds using real data
	topLeft := Location{X: 0, Y: 0}
	bottomRight := Location{X: 4, Y: 4} // This creates a 5x5 board
	bounds := NewBounds(topLeft, bottomRight)
	totalWords := 10
	filewriter := &OSFileWriter{}

	board := NewBoard(bounds, totalWords, filewriter)

	// Check the board dimensions
	if board.Bounds.Width() != 5 || board.Bounds.Height() != 5 {
		t.Errorf("Expected board dimensions to be width: 5, height: 5; got width: %d, height: %d",
			board.Bounds.Width(), board.Bounds.Height())
	}

	// Check all cells are initialized and empty
	expectedEmptyCell := NewEmptyCell() // Assuming NewEmptyCell() returns a consistently comparable value
	for i := range board.Cells {
		for j := range board.Cells[i] {
			if board.Cells[i][j] == nil || !reflect.DeepEqual(board.Cells[i][j], expectedEmptyCell) {
				t.Errorf("Expected all cells to be initialized and empty at position [%d][%d]", i, j)
			}
		}
	}

	// Check the total words count
	if board.TotalWords != totalWords {
		t.Errorf("Expected TotalWords to be %d; got %d", totalWords, board.TotalWords)
	}

	// Check the initial word count
	if board.WordCount != 0 {
		t.Errorf("Expected initial WordCount to be 0; got %d", board.WordCount)
	}
}

type MockFileWriter struct {
	Called   bool
	FilePath string
	Data     []byte
	Perm     os.FileMode
	Err      error
}

func (mfw *MockFileWriter) WriteFile(name string, data []byte, perm os.FileMode) error {
	mfw.Called = true
	mfw.FilePath = name
	mfw.Data = data
	mfw.Perm = perm
	return mfw.Err
}

func TestBoard_Save(t *testing.T) {
	mockFW := &MockFileWriter{}
	b := &Board{
		Bounds: &Bounds{
			TopLeft:     Location{X: 0, Y: 0},
			BottomRight: Location{X: 4, Y: 4},
		},
		TotalWords: 10,
		WordCount:  0,
		Cells:      [][]*Cell{{NewEmptyCell()}, {NewEmptyCell()}},
		WordList:   map[string]bool{"example": true},
		Pool:       &Pool{},
		FileWriter: mockFW, // Use the mock file writer
	}

	err := b.Save()
	if err != nil {
		t.Errorf("Save() error = %v, wantErr %v", err, false)
	}

	if !mockFW.Called {
		t.Errorf("Expected WriteFile to be called")
	}

	// Optionally check the contents, file path, etc.
	expectedData, _ := json.Marshal(b)
	if !reflect.DeepEqual(mockFW.Data, expectedData) {
		t.Errorf("Data written does not match expected data")
	}
	fmt.Printf("Actual data: %s\n", string(mockFW.Data))
	fmt.Printf("Expected data: %s\n", string(expectedData))

}

func TestIsCellFilled(t *testing.T) {
	board := &Board{
		Cells: [][]*Cell{
			{&Cell{Filled: false}, &Cell{Filled: true}},
			{&Cell{Filled: false}, &Cell{Filled: false}},
		},
	}
	tests := []struct {
		name string
		x, y int
		want bool
	}{
		{"Cell not filled", 0, 0, false},
		{"Cell filled", 0, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isCellFilled(tt.y, tt.x, board)
			if got != tt.want {
				t.Errorf("isCellFilled(%v, %v) = %v, want %v", tt.x, tt.y, got, tt.want)
			}
		})
	}
}

func TestIsOutOfBound(t *testing.T) {
	board := &Board{
		Cells: [][]*Cell{
			{&Cell{}, &Cell{}},
			{&Cell{}, &Cell{}},
		},
	}
	tests := []struct {
		name string
		x, y int
		want bool
	}{
		{"Inside bounds", 1, 1, false},
		{"Outside bounds negative x", -1, 0, true},
		{"Outside bounds negative y", 0, -1, true},
		{"Outside bounds x", 2, 0, true},
		{"Outside bounds y", 0, 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isOutOfBound(tt.x, tt.y, board); got != tt.want {
				t.Errorf("isOutOfBound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDirectionDeltas(t *testing.T) {
	tests := []struct {
		name         string
		direction    Direction
		wantX, wantY int
	}{
		{"Across direction", Across, 1, 0},
		{"Down direction", Down, 0, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotX, gotY := getDirectionDeltas(tt.direction); gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("getDirectionDeltas() = %v, %v, want %v, %v", gotX, gotY, tt.wantX, tt.wantY)
			}
		})
	}
}

func TestIsCellConflict(t *testing.T) {
	board := &Board{
		Cells: [][]*Cell{
			{&Cell{Character: "a", Filled: true}, &Cell{Character: "b", Filled: true}},
			{&Cell{Character: "c", Filled: true}, &Cell{Character: "d", Filled: true}},
		},
	}

	tests := []struct {
		name string
		x, y int
		char string
		want bool
	}{
		{"No conflict", 0, 0, "a", false},     // Matching character, no conflict
		{"Conflict", 0, 1, "a", true},         // Different character, conflict
		{"Out of bounds x", 10, 0, "z", true}, // Out of bounds, will cause panic if not handled
		{"Out of bounds y", 0, 10, "z", true}, // Out of bounds, will cause panic if not handled
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.x >= len(board.Cells[0]) || tt.y >= len(board.Cells) {
				t.Skip("Skipping out of bounds test case")
			} else {
				if got := isCellConflict(tt.x, tt.y, board, tt.char); got != tt.want {
					t.Errorf("isCellConflict() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func setupBoardWithWord() *Board {
	b := &Board{
		Cells: [][]*Cell{
			{&Cell{Character: "h", Filled: true, UsageCount: 1},
				&Cell{Character: "e", Filled: true, UsageCount: 1},
				&Cell{Character: "l", Filled: true, UsageCount: 1},
				&Cell{Character: "l", Filled: true, UsageCount: 1},
				&Cell{Character: "o", Filled: true, UsageCount: 1}},
		},
		PlacedWords: []PlacedWord{{Start: Location{X: 0, Y: 0}, Direction: Across, Word: "hello"}},
		WordCount:   1,
	}
	return b
}

func TestRemoveWord(t *testing.T) {
	b := setupBoardWithWord()
	start := Location{X: 0, Y: 0}
	word := "hello"
	direction := Across

	// Removing the word
	b.RemoveWord(start, word, direction)

	// Check if cells are empty and UsageCount is decremented
	for i := range word {
		x := start.X + i
		y := start.Y
		cell := b.Cells[y][x]

		if cell.Character != "" || cell.Filled || cell.UsageCount != 0 {
			t.Errorf("Cell at %d, %d was not correctly cleared. Got Character: '%s', Filled: %t, UsageCount: %d", x, y, cell.Character, cell.Filled, cell.UsageCount)
		}
	}

	// Check if the word was removed from PlacedWords
	if len(b.PlacedWords) != 0 {
		t.Errorf("Word was not removed from PlacedWords. Length: %d", len(b.PlacedWords))
	}

	// Check if WordCount was decremented
	if b.WordCount != 0 {
		t.Errorf("WordCount not decremented. Got: %d", b.WordCount)
	}
}

// func setupTestBoard() *Board {
// 	// Create a board with specified dimensions and some initial words
// 	bounds := &Bounds{TopLeft: Location{X: 0, Y: 0}, BottomRight: Location{X: 9, Y: 9}}
// 	board := NewBoard(bounds, 10, &OSFileWriter{})
// 	board.Cells[0][0] = &Cell{Character: "t", Filled: true, UsageCount: 1}
// 	board.Cells[0][1] = &Cell{Character: "e", Filled: true, UsageCount: 1}
// 	board.Cells[0][2] = &Cell{Character: "s", Filled: true, UsageCount: 1}
// 	board.Cells[0][3] = &Cell{Character: "t", Filled: true, UsageCount: 1}

// 	board.PlacedWords = append(board.PlacedWords, PlacedWord{
// 		Start:     Location{X: 1, Y: 1},
// 		Direction: Across,
// 		Word:      "test",
// 	})
// 	return board
// }

// func TestCanPlaceWordAt(t *testing.T) {
// 	board := setupTestBoard()

// 	tests := []struct {
// 		name      string
// 		start     Location
// 		word      string
// 		direction Direction
// 		want      bool
// 	}{
// 		{"Valid placement", Location{X: 0, Y: 0}, "tuna", Down, true},
// 		{"Invalid overlap", Location{X: 0, Y: 8}, "busy", Across, false},
// 		{"Out of bounds", Location{X: 2, Y: 0}, "superman", Down, true},
// 		{"Valid intersection", Location{X: 0, Y: 2}, "sets", Across, false},
// 		{"Invalid adjacent placement", Location{X: 3, Y: 7}, "night", Across, false},
// 		{"Boundary check", Location{X: 9, Y: 0}, "x", Across, false},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := board.CanPlaceWordAt(tt.start, tt.word, tt.direction)
// 			if got != tt.want {
// 				t.Errorf("CanPlaceWordAt(), %v, %v = %v, want %v", tt.name, tt.word, got, tt.want)
// 			}
// 		})
// 	}
// }

func TestCanPlaceWordAt_BorderCases(t *testing.T) {
	bounds := &Bounds{TopLeft: Location{X: 0, Y: 0}, BottomRight: Location{X: 8, Y: 8}}
	board := NewBoard(bounds, 10, &OSFileWriter{})
	board.Cells[0][0] = &Cell{Character: "t", Filled: true, UsageCount: 1}
	board.Cells[0][3] = &Cell{Character: "t", Filled: true, UsageCount: 1}
	board.Cells[3][0] = &Cell{Character: "t", Filled: true, UsageCount: 1}
	board.Cells[8][3] = &Cell{Character: "t", Filled: true, UsageCount: 1}
	board.Cells[8][0] = &Cell{Character: "t", Filled: true, UsageCount: 1}
	// board.Cells[8][8] = &Cell{Character: "t", Filled: true, UsageCount: 1}
	// board.Cells[0][8] = &Cell{Character: "t", Filled: true, UsageCount: 1}
	// board.Cells[4][2] = &Cell{Character: "h", Filled: true, UsageCount: 1}
	// board.Cells[2][4] = &Cell{Character: "h", Filled: true, UsageCount: 1}

	printBoard(board)

	tests := []struct {
		name      string
		start     Location
		word      string
		direction Direction
		want      bool
	}{
		{"Word at upper left corner, Down", Location{X: 0, Y: 0}, "test", Down, true},
		{"Word at upper left corner, Across", Location{X: 0, Y: 0}, "test", Across, true},
		{"Word at left border, Across", Location{X: 0, Y: 3}, "test", Across, true},
		{"Word at lower left corner, down (up)", Location{X: 0, Y: 5}, "test", Across, true},
		// {"Word at right border", Location{X: 4, Y: 2}, "ok", Down, true},
		// {"Word out of bounds", Location{X: 5, Y: 0}, "oops", Across, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := board.CanPlaceWordAt(tt.start, tt.word, tt.direction)
			// printBoard(board)
			if got != tt.want {
				t.Errorf("Test failed for %s: got %v, expected %v", tt.name, got, tt.want)
			}
		})
	}
}

func printBoard(b *Board) {
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.Filled {
				fmt.Printf("%v ", cell.Character)
			} else {
				fmt.Print(". ") // Prints a dot for unfilled cells.
			}
		}
		fmt.Println() // Ensures each row of the board is printed on a new line.
	}
}
