package generators

import (
	"errors"
	"fmt"
	"time"

	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

// AsymmetricalGenerator generates crossword puzzles without any
// symmetry considerations.
type AsymmetricalGenerator struct {
	*BaseGenerator // to reuse common fields and methods.
	WordPool       *models.Pool
}

func NewAsymmetricalGenerator(board *models.Board, pool *models.Pool) *AsymmetricalGenerator {
	return &AsymmetricalGenerator{
		BaseGenerator: NewBaseGenerator(board),
		WordPool:      pool,
	}
}

func (ag *AsymmetricalGenerator) placeFirstWord() error {
	if ag.Board == nil || len(ag.WordPool.Words) == 0 {
		return fmt.Errorf("uninitialized board or pool, or empty words list")
	}

	// Place the first word at the center horizontally.
	firstWord := ag.WordPool.Words[0]
	midRow := ag.Board.Bounds.Height() / 2
	startCol := (ag.Board.Bounds.Width() - len([]rune(firstWord))) / 2

	err := ag.Board.PlaceWordAt(models.Location{X: startCol, Y: midRow}, firstWord, models.Across)
	if err != nil {
		return errors.New("failed to place the first word")
	}
	return nil
}

// TODO-HneObo:
// Back-tracking. Keeop track of the best possible solution. Right now
// it's all or nothing. The best possible solution is the highest number of
// placed words with the given constraints.

// func (ag *AsymmetricalGenerator) Generate() error {
// 	err := ag.placeFirstWord()
// 	if err != nil {
// 		return err
// 	}

// 	return ag.placeWordsRecursive(1) // Start with the first word
// }

func (ag *AsymmetricalGenerator) Generate() error {
	// REM fmt.Println("Starting crossword generation...") // Debug

	err := ag.placeFirstWord()
	if err != nil {
		fmt.Println("Error placing first word:", err) // Debug
		return err
	}

	// REM fmt.Println("First word placed successfully.") // Debug

	err = ag.placeWordsRecursive(1) // Start from second word
	if err != nil {
		// REM fmt.Println("Word placement failed. Showing best attempt.")
		ag.Board.PrintBestSolution()
		return fmt.Errorf("crossword generation failed: %s", err)
	}

	fmt.Println("Crossword generation completed successfully!") // Debug
	return nil
}

// Recursive function to place words
var backtrackCount int = 0 // Global counter for backtracking

func (ag *AsymmetricalGenerator) placeWordsRecursive(index int) error {
	if index >= len(ag.WordPool.Words) {
		fmt.Println("\nAll words placed successfully!") // Ensure newline when done
		return nil
	}

	word := ag.WordPool.Words[index]
	placements := ag.FindPlacementLocations(word)

	if len(placements) == 0 {
		fmt.Println("\nNo placements found for:", word) // Ensure newline
		return fmt.Errorf("no placements available for word: %s", word)
	}

	failureCount := 0
	maxFailures := len(placements) / 2

	for _, location := range placements {
		if err := ag.Board.PlaceWordAt(location.Start, word, location.Direction); err == nil {
			err := ag.placeWordsRecursive(index + 1)
			if err == nil {
				return nil
			}

			// Backtrack: remove the word and try the next placement
			ag.Board.RemoveWord(location.Start, word, location.Direction)
			failureCount++
			backtrackCount++ // Increment global counter

			// Print backtracking count in place (overwrite previous line)
			fmt.Printf("\rBacktracking: %d", backtrackCount)
			time.Sleep(10 * time.Millisecond) // Small delay for visibility

			if failureCount >= maxFailures {
				// fmt.Println("\nToo many failed placements for word:", word) // Ensure newline
				return fmt.Errorf("too many failed placements for word: %s", word)
			}
		}
	}

	fmt.Println("\nFailed to place:", word) // Ensure newline
	return fmt.Errorf("failed to place word: %s", word)
}

// func (ag *AsymmetricalGenerator) placeWordsRecursive(index int) error {
// 	if index == len(ag.WordPool.Words) { // All words placed successfully
// 		return nil
// 	}

// 	word := ag.WordPool.Words[index]
// 	placements := ag.FindPlacementLocations(word)

// 	for _, location := range placements {
// 		if err := ag.Board.PlaceWordAt(location.Start, word, location.Direction); err == nil {
// 			if err := ag.placeWordsRecursive(index + 1); err == nil {
// 				return nil // Word placed successfully, recursion successful
// 			}
// 			// Backtrack: remove the word and try the next placement
// 			ag.Board.RemoveWord(location.Start, word, location.Direction)
// 		}
// 	}

// 	// Even if no word is placed at this step, save if it's the best state so far
// 	if ag.Board.WordCount > ag.Board.BestWordCount {
// 		ag.Board.SaveBestSolution()
// 	}

// 	return fmt.Errorf("failed to place word: %s", word)
// }

type Placement struct {
	Start     models.Location
	Direction models.Direction
}

// FindPlacementLocations generates a list of possible placement locations for a word.
func (ag *AsymmetricalGenerator) FindPlacementLocations(word string) []Placement {
	var placements []Placement

	// Helper function to try placing a word in one direction
	tryPlaceWord := func(x, y int, dir models.Direction) {
		if ag.Board.CanPlaceWordAt(models.Location{X: x, Y: y}, word, dir) {
			placements = append(placements, Placement{
				Start:     models.Location{X: x, Y: y},
				Direction: dir,
			})
		}
		// fmt.Printf("%v at x: %v, y: %v going %v\n", word, x, y, directionString(dir))
	}

	// Iterate over each cell in the board
	for y := 0; y < len(ag.Board.Cells); y++ {
		for x := 0; x < len(ag.Board.Cells[y]); x++ {
			tryPlaceWord(x, y, models.Across) // Try horizontal placement
			tryPlaceWord(x, y, models.Down)   // Try vertical placement
		}
	}
	return placements
}
