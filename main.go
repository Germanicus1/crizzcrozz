package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/Germanicus1/crizzcrozz/pkg/generators"
	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

// TODO-4po_rD: Implement command line arguments
// TODO-dNww-d: Refactor main.go

func main() {
	var width int
	height := 0
	flag.IntVar(&width, "width", 23, "The width of the board")
	flag.IntVar(&height, "height", height, "The width of the board")
	flag.Parse()

	if height == 0 {
		height = width
	}
	// Create the board bouindaries for the crosword puzzle
	bounds, _ := models.NewBoundsRectangle(width, height)

	// Create a word pool with some sample words
	//TODO-EiNHD4: Add words from the command line.
	//TODO-yUQLxC: Read words from file

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
