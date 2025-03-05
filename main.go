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
	findOptimum, width, height, maxRetries, fileName := parseFlags()
	// fileName := parseFlags()

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

	if !findOptimum {
		if width == 1 {
			fmt.Println("Specify a width for the board")
			return
		}

		height = width // Always a square board

		// Track the best attempt
		var bestBoard *models.Board
		maxWordsPlaced := 0

		for attempt := 0; attempt < maxRetries; attempt++ { // Limit attempts to prevent infinite loops
			tempBoard, _ := setUpBoard(width, height, len(words)) // Create a fresh board
			err := generateCrossword(tempBoard, words, maxRetries)

			if err == nil { // Success, all words fit
				bestBoard = tempBoard
				break
			}

			// Update bestBoard if this attempt placed more words
			if tempBoard.WordCount > maxWordsPlaced {
				maxWordsPlaced = tempBoard.WordCount
				bestBoard = tempBoard
			}
		}

		if bestBoard == nil {
			fmt.Println("No words could be placed. Try a larger board.")
			return
		}

		fmt.Printf("Board size: %dx%d | Words placed: %d/%d\n", width, height, bestBoard.WordCount, len(words))
		bestBoard.PrintBestSolution()
		return
	}

	// Find the best board using binary search
	board, bestSize, err := findOptimalBoardSize(words, maxRetries)
	if err != nil {
		log.Fatalf("Failed to generate crossword: %v", err)
	}

	fmt.Printf("Optimal board size found: %dx%d\n", bestSize, bestSize)
	board.PrintBestSolution()
}

// parseFlags returns the filename of the csv-file to parse.
func parseFlags() (bool, int, int, int, string) {
	var fileName string
	var width, height, maxRetries int
	var findOptimum bool
	flag.StringVar(&fileName, "f", "vocabulary.csv", "Specify the file with the words and hints. Defaults to vocabulary.csv.")
	flag.IntVar(&width, "w", 1, "Specify the width of the board. Defaults to 1.")
	flag.IntVar(&height, "h", 1, "Specify the width of the board. Defaults to 1")
	flag.IntVar(&maxRetries, "r", 0, "Specify the max number of retries to build the crossword. Defaults to 1.")
	flag.BoolVar(&findOptimum, "o", false, "Decide if the genereateor has to find th optimum board size.")
	flag.Parse()

	// return width, height, r, fileName
	return findOptimum, width, height, maxRetries, fileName
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
// func generateCrossword(b *models.Board, words []string) error {
// 	newPool := models.NewPool() // Creates a new pool to hold words.
// 	newPool.LoadWords(words)    // Loads words into the pool.

// 	generator := generators.NewAsymmetricalGenerator(b, newPool) // Initializes a new crossword generator.

//		err := generator.Generate() // Attempts to generate the crossword.
//		if err != nil {
//			return err // Returns an error if generation is unsuccessful.
//		}
//		return nil
//	}
func generateCrossword(b *models.Board, words []string, maxRetries int) error {
	newPool := models.NewPool()
	newPool.LoadWords(words)

	generator := generators.NewAsymmetricalGenerator(b, newPool)

	var bestBoard *models.Board
	maxWordsPlaced := 0

	for attempt := 0; attempt < maxRetries; attempt++ {
		fmt.Printf("Attempt %d/%d to generate crossword...\n", attempt+1, maxRetries)

		err := generator.Generate()
		if err == nil { // Success: all words fit
			fmt.Println("Successfully generated crossword.")
			return nil
		}

		// Track the best attempt
		if b.WordCount > maxWordsPlaced {
			maxWordsPlaced = b.WordCount
			bestBoard = b
		}

		fmt.Printf("Retry %d/%d: Words placed: %d/%d\n", attempt+1, maxRetries, b.WordCount, len(words))
	}

	// Ensure something is printed even if all retries fail
	if bestBoard != nil {
		fmt.Println("Could not fit all words. Showing best attempt:")
		bestBoard.PrintBestSolution()
		return fmt.Errorf("crossword generation failed after %d retries", maxRetries)
	}

	fmt.Println("No words could be placed. Try increasing the board size.")
	return fmt.Errorf("crossword generation failed completely")
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

func findOptimalBoardSize(words []string, maxRetries int) (*models.Board, int, error) {
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

		err = generateCrossword(board, words, maxRetries)
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

// printBoard outputs the current state of the crossword board to the
// console. It marks filled cells with their respective characters and
// empty cells with a dot.
func printBoard(b *models.Board) {
	for _, row := range b.Cells {
		for _, cell := range row {
			if cell.Filled {
				fmt.Printf("%v ", string(cell.Character))
			} else {
				fmt.Print(". ") // Prints a dot for unfilled cells.
			}
		}
		fmt.Println() // Ensures each row of the board is printed on a new line.
	}
}
