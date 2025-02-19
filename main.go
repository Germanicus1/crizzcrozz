package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Germanicus1/crizzcrozz/pkg/generators"
	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

func main() {
	// Initialize the board with specific dimensions
	// TODO: make sure width and height are uneven so (0,0) can be in the centre
	width, height := 23, 23

	// Create the board bouindaries for the crosword puzzle
	bounds, _ := models.NewBoundsRectangle(width, height)

	// Create a word pool with some sample words
	// words := []string{"examples", "mamma", "eat", "unmoor", "house", "stomp", "mustard", "school", "shoe"}

	words := []string{"bar", "beispiel", "bezahlen", "cent", "zusammen", "stimmt", "eingeladen", "essen", "euro", "gast", "kellner", "kellnerin", "rechnung", "sagen", "trinkgeld", "kosten", "viel", "zahlen", "karte", "getrennt", "zusammen"}

	for i, word := range words {
		words[i] = strings.ToUpper(word)
	}

	sort.Slice(words, func(j, i int) bool {
		return len(words[i]) < len(words[j])
	})

	board := models.NewBoard(bounds, len(words))

	// initialize a new pool of words.
	newPool := models.NewPool()
	newPool.LoadWords(words)

	generator := generators.NewAsymmetricalGenerator(board, newPool)

	// Generate the crossword
	err := generator.Generate()
	fmt.Println("wordcount:", board.WordCount)
	if err != nil {
		fmt.Println("Failed to generate the crossword:", err)
		printBoard(board)
		return
	}

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
