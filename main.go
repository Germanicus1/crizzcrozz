package main

import (
	"fmt"

	generator "github.com/Germanicus1/crizzcrozz/pkg/generators"
	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

func main() {
	// Initialize the board with specific dimensions
	// TODO: make sure width and height are uneven so (0,0) can be in the centre
	width, height := 11, 13
	// Create a word pool with some sample words
	words := []string{"example", "eat"}
	// words := []string{"example", "eagel", "eat", "long", "oxymoron", "car", "house"}
	// totalWords := 10
	// center := models.Location{X: width / 2, Y: height / 2} // Center

	// Create the board bouindaries for the crosword puzzle
	bounds, _ := models.NewBoundsRectangle(width, height)
	board := models.NewBoard(bounds, len(words))

	// initialize a new pool of words.
	newPool := models.NewPool()
	newPool.LoadWords(words)
	board.Pool = newPool // making sure we can access the pool of words from board struct

	generator := generator.NewAsymmetricalGenerator(board, newPool)

	// Generate the crossword
	err := generator.Generate()
	if err != nil {
		// FIXME: remove debug logging
		// printBoard(board)
		fmt.Println("Failed to generate the crossword:", err)
		return
	}

	// FIXME: remove debug logging
	printBoard(board)
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

// TODO: Generate a crossword
