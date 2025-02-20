package generators

import (
	"errors"
	"fmt"

	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

// AsymmetricalGenerator generates crossword puzzles without any
// symmetry considerations.
type AsymmetricalGenerator struct {
	*BaseGenerator // to reuse common fields and methods.
	WordPool       *models.Pool
}

// TODO-Ak2yMR: Clean up NewAsymmetricalGenerator
// TODO-0MsCHp: Refactor NewAsymmetricalGenerator

func NewAsymmetricalGenerator(board *models.Board, pool *models.Pool) *AsymmetricalGenerator {
	return &AsymmetricalGenerator{
		BaseGenerator: NewBaseGenerator(board),
		WordPool:      pool,
	}
}

// Generate implements the Generate method for generating asymmetrical crossword puzzles.
func (ag *AsymmetricalGenerator) Generate() error {
	// Place the first word at the center horizontally.
	firstWord := ag.WordPool.Words[0]
	midRow := ag.Board.Bounds.Height() / 2
	midCol := (ag.Board.Bounds.Width() - len(firstWord)) / 2

	err := ag.Board.PlaceWordAt(models.Location{X: midCol, Y: midRow}, firstWord, models.Across)
	if err != nil {
		return errors.New("failed to place the first word")
	}

	// // Initialize a queue with all the words
	wordQueue := make([]string, len(ag.WordPool.Words)-1)
	copy(wordQueue, ag.WordPool.Words[1:])
	maxRetries := 5                 // to prevent infinit loops
	retries := make(map[string]int) // to keep track of the number oif retries per string

	// Iterate through the rest of the words.
	for len(wordQueue) > 0 {
		word := wordQueue[0]
		wordQueue = wordQueue[1:] // taking of the first word of the list. Will be added again if placement was unsucessful
		placed := false
		placements := ag.FindPlacementLocations(word)

		for _, location := range placements {
			if ag.Board.CanPlaceWordAt(location.Start, word, location.Direction) {
				// FIXME-JM0DEa: this is checking a second time if the word can be placed. Why? See CanPlaceWordAt()
				err := ag.Board.PlaceWordAt(location.Start, word, location.Direction)
				if err != nil {
					fmt.Println("Error:", err)
					break
				}
				placed = true
				break
			}
		}

		fmt.Println("wordQueue: ", wordQueue)

		// TODO-zKYKMy: Decide what to do with words that cvould not be placed.
		// Backtrace?
		if !placed && retries[word] <= maxRetries {
			wordQueue = append(wordQueue, word) // Re-queue the word at the end
			retries[word]++
		}
	}

	if !ag.Board.IsComplete() {
		return errors.New("Generate: failed to generate a complete puzzle")
	}
	return nil
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
