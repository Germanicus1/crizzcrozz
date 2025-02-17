package main

import (
	"fmt"

	"github.com/Germanicus1/crizzcrozz/pkg/generators"
	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

func main() {
	// Initialize the board with specific dimensions
	// TODO: make sure width and height are uneven so (0,0) can be in the centre
	width, height := 15, 15

	// Create the board bouindaries for the crosword puzzle
	bounds, _ := models.NewBoundsRectangle(width, height)

	// Create a word pool with some sample words
	words := []string{"examples", "eat", "unmoor", "mamma", "house", "stomp", "mustard", "school", "shoe"}

	board := models.NewBoard(bounds, len(words))

	// initialize a new pool of words.
	newPool := models.NewPool()
	newPool.LoadWords(words)

	// REM: debug info
	fmt.Println("newpool.words:", newPool.Words)
	fmt.Println("newpool.wordSet:", newPool.WordSet)
	fmt.Println("newpool.ByLength:", newPool.ByLength)

	generator := generators.NewAsymmetricalGenerator(board, newPool)

	// Generate the crossword
	err := generator.Generate()
	fmt.Println("wordcount:", board.WordCount)
	if err != nil {
		fmt.Println("Failed to generate the crossword:", err)
		printBoard(board)
		return
	}

	// REM: printBoard debug info
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
