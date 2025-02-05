package main

import (
	"fmt"

	generator "github.com/Germanicus1/crizzcrozz/pkg/generators"
	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

func main() {
	// Initialize the board with specific dimensions
	width, height := 21, 19
	// Create a word pool with some sample words
	words := []string{"example", "crossword", "puzzle", "symmetry", "generation", "horizontal", "vertical", "language", "programming", "golang"}
	// totalWords := 10
	// center := models.Location{X: width / 2, Y: height / 2} // Center

	// Create the board bouindaries for the crosword puzzle
	bounds, _ := models.NewBoundsRectangle(width, height)
	board := models.NewBoard(bounds, len(words))

	// FIXME: remove debug logging
	fmt.Println(board.Cells[0])

	// initialize a new pool of words.
	newPool := models.NewPool()
	newPool.LoadWords(words)

	generator := generator.NewAsymmetricalGenerator(board, newPool)

	// Generate the crossword
	err := generator.Generate()
	if err != nil {
		fmt.Println("Failed to generate the crossword:", err)
		return
	}

	// fmt.Println(b)
}

// printBoard prints the crossword board to the console.
func printBoard(b *models.Board) {
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.Filled {
				fmt.Printf("%c ", cell.Character)
			} else {
				fmt.Print(". ")
			}
		}
		fmt.Println()
	}
}

// DONE: Print empty board
// TODO: Generate a crossword
