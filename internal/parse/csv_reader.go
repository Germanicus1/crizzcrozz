package parse

import (
	"fmt"
	"os"

	"github.com/Germanicus1/crizzcrozz/internal/models"
	"github.com/gocarina/gocsv"
)

// readWordsFromFile reads words and their hints from a specified CSV
// file. It returns a slice of wordsAndHints structs or an error if the
// file cannot be read.
func ReadWordsFromFile(fileName string) ([]*models.WordsAndHints, error) {
	csvFile, err := os.OpenFile(fileName, os.O_RDWR, os.ModePerm) // Opens the CSV file for reading.
	if err != nil {
		return nil, err // Returns an error if the file cannot be opened.
	}
	defer csvFile.Close() // Ensures the file is closed after the operation.

	var wordsAndHints []*models.WordsAndHints
	if err := gocsv.UnmarshalFile(csvFile, &wordsAndHints); err != nil {
		return nil, err // Returns an error if the CSV data cannot be parsed.
	}
	fmt.Printf("Successfully unmarshalled %d words and hints.\n", len(wordsAndHints))

	return wordsAndHints, nil
}
