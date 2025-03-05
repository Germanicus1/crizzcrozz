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

func (ag *AsymmetricalGenerator) placeWordsRecursive(index int) error {
	if index >= len(ag.WordPool.Words) {
		fmt.Println("\nAll words placed successfully!")
		ag.Board.SaveBestSolution() // Ensure final board is saved
		return nil
	}

	// Stop if we've backtracked too much
	if backtrackCount >= maxBacktracks {
		fmt.Println("\nMax backtracking limit reached. Stopping word placement.")
		ag.Board.SaveBestSolution() // Save before stopping
		return fmt.Errorf("stopping due to excessive backtracking")
	}

	word := ag.WordPool.Words[index]
	// REM fmt.Printf("Placing word #%d: %s\n", index, word)

	placements := ag.FindPlacementLocations(word)

	if len(placements) == 0 {
		// REM fmt.Println("No placements found for:", word)
		return fmt.Errorf("no placements available for word: %s", word)
	}

	failureCount := 0
	maxFailsPerWord := len(placements) / 2

	for _, location := range placements {
		if err := ag.Board.PlaceWordAt(location.Start, word, location.Direction); err == nil {
			err := ag.placeWordsRecursive(index + 1)
			if err == nil {
				return nil
			}

			// Backtrack: remove the word and try the next placement
			ag.Board.RemoveWord(location.Start, word, location.Direction)
			backtrackCount++

			// Save the best attempt before stopping
			if ag.Board.WordCount > ag.Board.BestWordCount {
				ag.Board.SaveBestSolution()
			}

			failureCount++
			if failureCount >= maxFailsPerWord {
				// REM fmt.Println("\nToo many failed placements for word:", word)
				return fmt.Errorf("too many failed placements for word: %s", word)
			}
		}
	}

	// REM fmt.Printf("\nFailed to place '%s'. Stopping recursion.\n", word)
	return fmt.Errorf("failed to place word: %s", word)
}

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
