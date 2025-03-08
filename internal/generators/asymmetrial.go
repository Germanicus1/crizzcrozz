package generators

import (
	"errors"
	"fmt"

	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

var backtrackCount int = 0  // Global counter for backtracking
const maxBacktracks = 50000 // Limit for stopping excessive backtracking

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

func (ag *AsymmetricalGenerator) Generate() error {
	backtrackCount = 0 // Reset counter before recursion starts

	// REM fmt.Println("Starting crossword generation...")

	err := ag.placeFirstWord()
	if err != nil {
		fmt.Println("Error placing first word:", err)
		return err
	}

	err = ag.placeWordsRecursive(1) // Start from the second word
	// fmt.Println("\nBacktracking limit reached or crossword generation failed.")
	return err
}

// REM debugging
func (ag *AsymmetricalGenerator) placeWordsRecursive(index int) error {
	if index >= len(ag.WordPool.Words) {
		fmt.Println("\n✅ All words placed successfully! Saving best solution...")
		ag.Board.SaveBestSolution()
		return nil
	}

	word := ag.WordPool.Words[index]
	fmt.Printf("\n🔍 Trying to place word #%d: %s\n", index+1, word)

	placements := ag.FindPlacementLocations(word)
	fmt.Printf("➡️ Available placements for %s: %d\n", word, len(placements))

	if len(placements) == 0 {
		fmt.Printf("❌ No placements found for: %s\n", word)
		return fmt.Errorf("no placements available for word: %s", word)
	}

	for _, location := range placements {
		// 🚨 Check if the placement is valid before proceeding
		if !ag.Board.CanPlaceWordAt(location.Start, word, location.Direction) {
			fmt.Printf("⚠️ Skipping invalid placement: %s at (%d, %d) %v (Would overwrite another word)\n",
				word, location.Start.X, location.Start.Y, location.Direction)
			continue
		}

		if err := ag.Board.PlaceWordAt(location.Start, word, location.Direction); err == nil {
			fmt.Printf("✅ Successfully placed word: %s at (%d, %d) %v\n",
				word, location.Start.X, location.Start.Y, location.Direction)

			err := ag.placeWordsRecursive(index + 1)
			if err == nil {
				return nil
			}

			// Backtrack: remove the word and try the next placement
			ag.Board.RemoveWord(location.Start, word, location.Direction)
			fmt.Printf("🔄 Backtracking: Removed word %s from (%d, %d) %v\n",
				word, location.Start.X, location.Start.Y, location.Direction)
		}
	}

	fmt.Printf("⚠️ Failed to place word: %s\n", word)
	return fmt.Errorf("failed to place word: %s", word)
}

// func (ag *AsymmetricalGenerator) placeWordsRecursive(index int) error {
// 	if index >= len(ag.WordPool.Words) {
// 		fmt.Println("\nAll words placed successfully!")
// 		ag.Board.SaveBestSolution() // Ensure final board is saved
// 		return nil
// 	}

// 	word := ag.WordPool.Words[index]
// 	// REM debug info
// 	// fmt.Printf("Trying to place word: %s\n", word)

// 	placements := ag.FindPlacementLocations(word)
// 	// REM debug info
// 	// fmt.Printf("Available placements for %s: %d\n", word, len(placements))

// 	if len(placements) == 0 {
// 		// REM debug info
// 		// fmt.Printf("No placements found for: %s\n", word)
// 		return fmt.Errorf("no placements available for word: %s", word)
// 	}

// 	for _, location := range placements {
// 		if err := ag.Board.PlaceWordAt(location.Start, word, location.Direction); err == nil {
// 			// REM debug info
// 			// fmt.Printf("Placed word: %s at (%d, %d) %v\n", word, location.Start.X, location.Start.Y, location.Direction)

// 			err := ag.placeWordsRecursive(index + 1)
// 			if err == nil {
// 				return nil
// 			}

// 			// Backtrack: remove the word and try the next placement
// 			ag.Board.RemoveWord(location.Start, word, location.Direction)
// 			// REM debug info
// 			// fmt.Printf("Removed word: %s from (%d, %d) %v (Backtracking)\n", word, location.Start.X, location.Start.Y, location.Direction)
// 		}
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
