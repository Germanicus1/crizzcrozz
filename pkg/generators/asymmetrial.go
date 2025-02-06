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

// NewAsymmetricalGenerator creates a new generator for asymmetrical puzzles.
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

	// Iterate through the rest of the words.
	for _, word := range ag.WordPool.Words[1:] {
		placed := false
		for _, location := range ag.FindPlacementLocations(word) {
			// FIXME: remove debug info
			// fmt.Printf("Location %+v for %s\n", location, word)

			if ag.Board.CanPlaceWordAt(location.Start, word, location.Direction) {
				err := ag.Board.PlaceWordAt(location.Start, word, location.Direction)
				fmt.Println("Error:", err)
				placed = true
				break
			}
		}
		if !placed {
			return errors.New("failed to place a word: " + word)
		}
	}

	if !ag.Board.IsComplete() {
		return errors.New("failed to generate a complete puzzle")
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
	}

	// Iterate over each cell in the board
	for y := 0; y < len(ag.Board.Cells); y++ {
		for x := 0; x < len(ag.Board.Cells[y]); x++ {
			tryPlaceWord(x, y, models.Across) // Try horizontal placement
			tryPlaceWord(x, y, models.Down)   // Try vertical placement
		}
	}

	// FIXME: remove debug info
	fmt.Printf("Placements for %s: %+v\n", word, placements)

	// Logic to find potential placements based on existing words on the board.
	return placements
}
