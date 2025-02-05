package generator

import (
	"errors"

	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

// Generator defines the interface for generating crossword puzzles.
type Generator interface {
	Generate() error // Generate the crossword puzzle, returning an error if unsuccessful.
}

// BaseGenerator provides a basic structure and common functionality for
// crossword generators.
type BaseGenerator struct {
	Board *models.Board // A reference to the board where the crossword will be generated.
}

// NewBaseGenerator creates a new instance of BaseGenerator with
// specified board boundaries.
func NewBaseGenerator(b *models.Board) *BaseGenerator {
	return &BaseGenerator{
		Board: b,
	}
}

// TODO: Implement a generator

// Generate is a placeholder to satisfy the Generator interface. Specific
// generator implementations should override this method with actual
// logic.
func (bg *BaseGenerator) Generate() error {
	return errors.New("Generate method not implemented")
}

// ValidatePlacement checks if a given word can be placed at a starting
// location in a specified direction.
func (bg *BaseGenerator) ValidatePlacement(word string, start models.Location, direction models.Direction) bool {
	// Placeholder for validation logic
	return true
}
