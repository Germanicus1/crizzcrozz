package generators

import (
	"errors"

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
	// Basic implementation outline:
	// 1. Select a starting point and direction randomly or based on some heuristic.
	// 2. Choose a word from the pool that fits the selected location and direction.
	// 3. Place the word on the board if it fits without conflicts.
	// 4. Repeat until no more words can be placed or the puzzle meets the design criteria.

	// Example placeholder logic (needs actual implementation details):
	for _, word := range ag.WordPool.Words {
		// Attempt to place each word in the pool.
		startLocation := models.Location{X: 0, Y: 0} // FIXME: This should be determined by some logic.
		direction := models.Across                   // FIXME: This should also be decided dynamically.

		if ag.Board.CanPlaceWordAt(startLocation, word, direction) {
			ag.Board.PlaceWordAt(startLocation, word, direction)
		} else {
			// If a word cannot be placed, you might skip or try a different location/direction.
			continue
		}
	}

	// Check if the board is satisfactorily filled or if more words need to be placed.
	if ag.Board.IsComplete() {
		return nil
	}

	return errors.New("failed to generate a complete puzzle")
}

// Additional methods specific to asymmetrical puzzle logic can be added here.
