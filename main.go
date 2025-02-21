package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"flag"

	"github.com/Germanicus1/crizzcrozz/pkg/generators"
	"github.com/Germanicus1/crizzcrozz/pkg/models"
	"github.com/gocarina/gocsv"
)

// wordsAndHints defines a struct for mapping words and their hints from a CSV file.
type wordsAndHints struct {
	Word string `csv:"word"`
	Hint string `csv:"hint"`
}

// Define custom error for invalid dimensions.
var ErrInvalidDimensions = errors.New("invalid board dimensions")

func main() {
	width, height := parseFlags()
	fileName := "vocabulary.csv" // Fallback file name FIXME-0DTbqN: The filename needs to come from the command line too

	wordsAndHints, err := readWordsFromFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("File does not exist: %v", err)
		} else {
			log.Fatalf("Unknown error occurred: %v", err)
		}
	}

	var words []string // Initializes a slice to store words extracted from wordsAndHints.
	for _, v := range wordsAndHints {
		words = append(words, v.Word) // Populates the words slice with words for the crossword.
	}

	// Sets up the crossword board with the specified dimensions and words.
	board, err := setUpBoard(width, height, words)
	if err != nil {
		// Check if the error is due to invalid dimensions
		if errors.Is(err, ErrInvalidDimensions) {
			log.Fatalf("Cannot create a board with invalid dimensions: %v", err)
		} else {
			log.Printf("Failed to generate the crossword: %v", err)
			if board != nil {
				printBoard(board) // Attempts to print the current state of the board before exiting.
			}
			log.Fatal("Exiting due to unrecoverable error setting up the board.")
		}
	}
	// Attempts to generate the crossword puzzle using the board setup.
	if err := generateCrossword(board, words); err != nil {
		fmt.Println("Failed to generate the crossword:", err)
		printBoard(board) // Prints the board even if the puzzle generation fails.
		return
	}
	printBoard(board) // Prints the successfully generated crossword board.
}

// printBoard outputs the current state of the crossword board to the console.
// It marks filled cells with their respective characters and empty cells with a dot.
func printBoard(b *models.Board) {
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.Filled {
				fmt.Printf("%c ", cell.Character) // Prints the character for filled cells.
			} else {
				fmt.Print(". ") // Prints a dot for unfilled cells.
			}
		}
		fmt.Println() // Ensures each row of the board is printed on a new line.
	}
}

// parseFlags parses the width and height command-line arguments. It
// returns the parsed dimensions, using width for height if height is
// not specified. Defaults to a square of 23
func parseFlags() (int, int) {
	var width, height int
	flag.IntVar(&width, "width", 23, "Specify the width of the board. Default is 23.")
	flag.IntVar(&height, "height", 0, "Specify the height of the board. Defaults to the value of width if not set.")
	flag.Parse()

	if height == 0 {
		height = width
	}

	return width, height
}

// processWordList cleans up and sorts a list of words.
// It trims whitespace, converts words to uppercase, and sorts by length.
func processWordList(wordList string) []string {
	var words []string
	if wordList != "" {
		words = strings.Split(wordList, ",")
		for i, v := range words {
			words[i] = strings.TrimSpace(v) // Trims any surrounding whitespace.
		}
	} else {
		// Fallback word list if none is provided.
		words = []string{"bar", "beispiel", "bezahlen", "cent", "zusammen", "stimmt", "eingeladen", "essen", "euro", "gast", "kellner", "kellnerin", "rechnung", "sagen", "trinkgeld", "kosten", "viel", "zahlen", "karte", "getrennt", "zusammen"}
	}

	for i, word := range words {
		words[i] = strings.ToUpper(word) // Converts all words to uppercase.
	}

	sort.Slice(words, func(j, i int) bool {
		return len(words[i]) < len(words[j]) // Sorts words by length.
	})
	return words
}

// setUpBoard initializes a crossword board with given dimensions and a list of words.
// It returns a pointer to the created board or an error if the board cannot be created.
func setUpBoard(width, height int, words []string) (*models.Board, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid board dimensions (width: %d, height: %d)", width, height)
	}

	bounds, err := models.NewBoundsRectangle(width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to create board boundaries: %w", err)
	}

	board := models.NewBoard(bounds, len(words))
	if board == nil {
		return nil, fmt.Errorf("failed to initialize the crossword board")
	}

	return board, nil
}

// generateCrossword tries to populate the crossword board with words.
// It returns an error if the crossword generation fails.
func generateCrossword(b *models.Board, words []string) error {
	newPool := models.NewPool() // Creates a new pool to hold words.
	newPool.LoadWords(words)    // Loads words into the pool.

	generator := generators.NewAsymmetricalGenerator(b, newPool) // Initializes a new crossword generator.

	err := generator.Generate() // Attempts to generate the crossword.
	if err != nil {
		return err // Returns an error if generation is unsuccessful.
	}
	return nil
}

// readWordsFromFile reads words and their hints from a specified CSV file.
// It returns a slice of wordsAndHints structs or an error if the file cannot be read.
func readWordsFromFile(fileName string) ([]*wordsAndHints, error) {
	csvFile, err := os.OpenFile(fileName, os.O_RDWR, os.ModePerm) // Opens the CSV file for reading.
	if err != nil {
		return nil, err // Returns an error if the file cannot be opened.
	}
	defer csvFile.Close() // Ensures the file is closed after the operation.

	var wordsAndHints []*wordsAndHints
	if err := gocsv.UnmarshalFile(csvFile, &wordsAndHints); err != nil {
		return nil, err // Returns an error if the CSV data cannot be parsed.
	}

	return wordsAndHints, nil
}
