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

func (ag *AsymmetricalGenerator) Generate() error {
	err := ag.placeFirstWord()
	if err != nil {
		return err
	}

	return ag.placeWordsRecursive(0) // Start with the first word
}

// Recursive function to place words
func (ag *AsymmetricalGenerator) placeWordsRecursive(index int) error {
	if index == len(ag.WordPool.Words) { // All words placed successfully
		return nil
	}

	word := ag.WordPool.Words[index]
	placements := ag.FindPlacementLocations(word)

	for _, location := range placements {
		if err := ag.Board.PlaceWordAt(location.Start, word, location.Direction); err == nil {
			if err := ag.placeWordsRecursive(index + 1); err == nil {
				return nil // Word placed successfully, recursion successful
			}
			// Backtrack: remove the word and try the next placement
			ag.Board.RemoveWord(location.Start, word, location.Direction)
		}
	}

	// Even if no word is placed at this step, save if it's the best state so far
	if ag.Board.WordCount > ag.Board.BestWordCount {
		ag.Board.SaveBestSolution()
	}

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
