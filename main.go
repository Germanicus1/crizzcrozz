package main

import (
	"errors"
	"fmt"
	"log"
	"math"
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
	estimate, findOptimalSize, width, _, maxRetries, fileName := parseFlags()
	// fileName := parseFlags()
	if !findOptimalSize && !estimate && width == 1 {
		// FIXME: Improve error handling
		fmt.Println("You need to specify a reasonable width for the board. Use the -f=<size> or -e=TRUE for estimating a size.")
		return
	}

	fmt.Println("estimate:", estimate)
	fmt.Println("findOptimalSize:", findOptimalSize)
	// findOptimalSize = false

	wordsAndHints, err := readWordsFromFile(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Fatalf("File does not exist: %v", err)
		} else {
			log.Fatalf("Unknown error occurred: %v", err)
		}
	}

	cleanedWords := cleanWords(wordsAndHints)
	sortedWords := sortWordsByLength(cleanedWords)

	var bestBoard *models.Board

	if estimate {
		width = estimateInitialBoardSize(sortedWords)
	}
	bestBoard = createBoard(sortedWords, maxRetries, width)
	if bestBoard == nil {
		fmt.Println("No words could be placed. Try a larger board.")
		return
	}

	fmt.Printf("Board size: %dx%d | Words placed: %d/%d\n", bestBoard.Bounds.Width(), bestBoard.Bounds.Height(), bestBoard.BestWordCount, bestBoard.TotalWords)
	bestBoard.PrintBestSolution()

	// TODO: Implement discovery of smallest board size possible with all the
	// words placed.

	// board, bestSize, err := findOptimalBoardSize(words, maxRetries)
	// if err != nil {
	// 	log.Fatalf("Failed to generate crossword: %v", err)
	// }

	// fmt.Printf("Optimal board size found: %dx%d\n", bestSize, bestSize)
	// board.PrintBestSolution()
}

// parseFlags returns the filename of the csv-file to parse.
func parseFlags() (bool, bool, int, int, int, string) {
	var fileName string
	var width, height, maxRetries int
	var findOptimalSize, estimate bool
	flag.StringVar(&fileName, "f", "vocabulary_eng.csv", "Specify the file with the words and hints. Defaults to vocabulary.csv.")
	flag.IntVar(&width, "w", 1, "Specify the width of the board. Defaults to 1.")
	flag.IntVar(&height, "h", 1, "Specify the width of the board. Defaults to 1")
	flag.IntVar(&maxRetries, "r", 1, "Specify the max number of retries to build the crossword. Defaults to 1.")
	flag.BoolVar(&findOptimalSize, "o", false, "Decide if the generator has to find th optimum board size. Default FALSE.")
	flag.BoolVar(&estimate, "e", true, "Decide if the program can calculate an estimated board size. Default FALSE.")
	flag.Parse()

	return estimate, findOptimalSize, width, height, maxRetries, fileName
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
		// REM fmt.Printf("Attempt %d/%d to generate crossword...\n", attempt+1, maxRetries)

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

		// REM fmt.Printf("Retry %d/%d: Words placed: %d/%d\n", attempt+1, maxRetries, b.WordCount, len(words))
	}

	// Ensure something is printed even if all retries fail
	if bestBoard != nil {
		fmt.Println("Could not fit all words. Showing best attempt:")
		// bestBoard.PrintBestSolution()
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

// func estimateInitialBoardSize(words []string) int {
// 	// totalLength := 0
// 	// longestWord := 0

// 	// for _, word := range words {
// 	// 	wordLen := len([]rune(word))
// 	// 	totalLength += wordLen
// 	// 	if wordLen > longestWord {
// 	// 		longestWord = wordLen
// 	// 	}
// 	// }

// 	// // Start with something proportional to total word length
// 	// estimatedSize := int(float64(totalLength) * 0.7) // 70% of total length
// 	// if estimatedSize < longestWord {                 // Ensure at least the longest word fits
// 	// 	estimatedSize = longestWord + 2
// 	// }

// 	return int(float64(len([]rune(words[0]))) * 1.5)
// }

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

func estimateInitialBoardSize(words []string) int {
	wordCount := len(words)
	if wordCount == 0 {
		return 10 // Default minimum size
		// FIXME: This should probably retrun an error.
	}

	longestWord := 0
	totalLength := 0

	for _, word := range words {
		wordLen := len([]rune(word))
		totalLength += wordLen
		if wordLen > longestWord {
			longestWord = wordLen
		}
	}

	averageWordLength := totalLength / wordCount

	// Adjust density factor based on expected intersection
	densityFactor := 1.2 // Increase density factor for better spacing

	// Adjust padding dynamically based on longest word and total words
	padding := int(math.Max(float64(longestWord)*0.3, float64(wordCount)*0.3)) // More words → more padding

	// Estimate board size using an improved formula
	estimatedSize := int(math.Sqrt(float64(wordCount) * float64(averageWordLength) * densityFactor))

	// Ensure it's at least large enough for the longest word + padding
	if estimatedSize < longestWord+padding {
		estimatedSize = longestWord + padding
	}

	// Add extra space for flexibility
	// estimatedSize += 3

	return estimatedSize
}

func cleanWords(wh []*wordsAndHints) []string {
	var words []string
	for _, v := range wh {
		cleanWord := strings.TrimSpace(v.Word)
		words = append(words, cleanWord)
	}
	return words
}

func createBoard(sortedWords []string, maxRetries, width int) *models.Board {
	height := width // Always a square board

	// Track the best attempt
	var bestBoard *models.Board
	maxWordsPlaced := 0

	for attempt := 0; attempt < maxRetries; attempt++ { // Limit attempts to prevent infinite loops
		tempBoard, _ := setUpBoard(width, height, len(sortedWords)) // Create a fresh board
		err := generateCrossword(tempBoard, sortedWords, maxRetries)

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

	return bestBoard
}
