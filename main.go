package main

import (
	"fmt"
	"sort"
	"strings"

	"flag"

	"github.com/Germanicus1/crizzcrozz/pkg/generators"
	"github.com/Germanicus1/crizzcrozz/pkg/models"
)

//TODO-yUQLxC: Read words from file

func main() {
	width, height, wordList := parseFlags()
	words := processWordList(wordList)

	board, err := setUpBoard(width, height, words)
	if err != nil {
		fmt.Println("Failed to generate the crossword:", err)
		return
	}

	if err := generateCrossword(board, words); err != nil {
		fmt.Println("Failed to generate the crossword:", err)
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

func parseFlags() (int, int, string) {
	var width, height int
	var wordList string
	flag.IntVar(&width, "width", 23, "The width of the board")
	flag.IntVar(&height, "height", 0, "The width of the board")
	flag.StringVar(&wordList, "words", "", "A comma separated list of words")
	flag.Parse()

	if height == 0 {
		height = width
	}

	return width, height, wordList
}

func processWordList(wordList string) []string {
	var words []string
	if wordList != "" {
		words = strings.Split(wordList, ",")
		for i, v := range words {
			words[i] = strings.TrimSpace(v)
		}
	} else {
		words = []string{"bar", "beispiel", "bezahlen", "cent", "zusammen", "stimmt", "eingeladen", "essen", "euro", "gast", "kellner", "kellnerin", "rechnung", "sagen", "trinkgeld", "kosten", "viel", "zahlen", "karte", "getrennt", "zusammen"}
	}

	for i, word := range words {
		words[i] = strings.ToUpper(word)
	}

	sort.Slice(words, func(j, i int) bool {
		return len(words[i]) < len(words[j])
	})
	return words
}

func setUpBoard(width, height int, words []string) (*models.Board, error) {
	// Create the board bouindaries for the crossword puzzle
	bounds, err := models.NewBoundsRectangle(width, height)
	if err != nil {
		return nil, err
	}
	board := models.NewBoard(bounds, len(words))
	return board, nil
}

func generateCrossword(b *models.Board, words []string) error {
	// initialize a new pool of words.
	newPool := models.NewPool()
	newPool.LoadWords(words)

	generator := generators.NewAsymmetricalGenerator(b, newPool)

	// Generate the crossword
	err := generator.Generate()
	// fmt.Println("wordcount:", b.WordCount)
	if err != nil {
		fmt.Println("Failed to generate the crossword:", err)
		// printBoard(board)
		return err
	}
	return nil
}
