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

// wordsAndHints defines a struct for mapping words and their hints from
// a CSV file.
type wordsAndHints struct {
	Word string `csv:"word"`
	Hint string `csv:"hint"`
}

// Define custom error for invalid dimensions.
var ErrInvalidDimensions = errors.New("invalid board dimensions")

func main() {
	// width, height, maxRetries, fileName := parseFlags()
	fileName := parseFlags()

	wordsAndHints, err := readWordsFromFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("File does not exist: %v", err)
		} else {
			log.Fatalf("Unknown error occurred: %v", err)
		}
	}

	// fmt.Printf("Read %d words and hints from the file.\n", len(wordsAndHints))

	var words []string
	for _, v := range wordsAndHints {
		cleanWord := strings.TrimSpace(v.Word)
		words = append(words, cleanWord)
	}

	words = sortWordsByLength(words)

	// Find the best board using binary search
	board, bestSize, err := findOptimalBoardSize(words)
	if err != nil {
		log.Fatalf("Failed to generate crossword: %v", err)
	}

	fmt.Printf("Optimal board size found: %dx%d\n", bestSize, bestSize)
	board.PrintBestSolution()
}

// parseFlags returns the filename of the csv-file to parse.
func parseFlags() string {
	var fileName string
	flag.StringVar(&fileName, "f", "vocabulary.csv", "Specify the file with the words and hints. Defaults to vocabulary.csv.")
	flag.Parse()

	// return width, height, r, fileName
	return fileName
}

// setUpBoard initializes a crossword board with given dimensions and a
// list of words. It returns a pointer to the created board or an error
// if the board cannot be created.
func setUpBoard(width, height int, wordCount int) (*models.Board, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("invalid board dimensions (width: %d, height: %d)", width, height)
	}

	bounds, err := models.NewBoundsRectangle(width, height)
	if err != nil {
		return nil, fmt.Errorf("failed to create board boundaries: %w", err)
	}

	fileWriter := &models.OSFileWriter{} // Creating an instance of FileWriter
	board := models.NewBoard(bounds, wordCount, fileWriter)
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

// readWordsFromFile reads words and their hints from a specified CSV
// file. It returns a slice of wordsAndHints structs or an error if the
// file cannot be read.
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
	fmt.Printf("Successfully unmarshalled %d words and hints.\n", len(wordsAndHints))

	return wordsAndHints, nil
}

func sortWordsByLength(words []string) []string {
	sort.Slice(words, func(j, i int) bool {
		return len(words[i]) < len(words[j]) // Sorts words by length.
	})
	return words
}

func findOptimalBoardSize(words []string) (*models.Board, int, error) {
	low := estimateInitialBoardSize(words) / 2 // Start smaller
	high := low * 3                            // Start with a reasonable max size

	var bestBoard *models.Board
	var bestSize int

	for low <= high {
		mid := (low + high) / 2
		fmt.Printf("Trying board size: %dx%d\n", mid, mid)

		board, err := setUpBoard(mid, mid, len(words))
		if err != nil {
			return nil, 0, err
		}

		err = generateCrossword(board, words)
		if err == nil { // Success: all words fit
			bestBoard = board
			bestSize = mid
			high = mid - 1 // Try a smaller size
		} else {
			low = mid + 1 // Increase the size
		}
	}

	if bestBoard == nil {
		return nil, 0, fmt.Errorf("could not fit words on any reasonable board size")
	}

	return bestBoard, bestSize, nil
}

func estimateInitialBoardSize(words []string) int {
	// totalLength := 0
	// longestWord := 0

	// for _, word := range words {
	// 	wordLen := len([]rune(word))
	// 	totalLength += wordLen
	// 	if wordLen > longestWord {
	// 		longestWord = wordLen
	// 	}
	// }

	// // Start with something proportional to total word length
	// estimatedSize := int(float64(totalLength) * 0.7) // 70% of total length
	// if estimatedSize < longestWord {                 // Ensure at least the longest word fits
	// 	estimatedSize = longestWord + 2
	// }

	return int(float64(len([]rune(words[0]))) * 1.5)
}
