package parse_test

import (
	"os"
	"testing"

	"github.com/Germanicus1/crizzcrozz/internal/models"
	"github.com/Germanicus1/crizzcrozz/internal/parse"
)

// Utility function to create a mock CSV file for testing.
func createMockCSVFile(content string) (string, func(), error) {
	file, err := os.CreateTemp("", "mock.csv")
	if err != nil {
		return "", nil, err
	}
	_, err = file.WriteString(content)
	if err != nil {
		return "", nil, err
	}
	return file.Name(), func() { os.Remove(file.Name()) }, nil
}

func testReadWordsFromFile(t *testing.T) {
	// Prepare a mock CSV file.
	content := "word,hint\napple,Fruit\nsky,Blue"
	fileName, cleanup, err := createMockCSVFile(content)
	if err != nil {
		t.Fatalf("Failed to create mock CSV file: %s", err)
	}
	defer cleanup()

	t.Run("read a csv file from disc", func(t *testing.T) {
		// Expected result.
		expected := []*models.WordsAndHints{
			{Word: "apple", Hint: "Fruit"},
			{Word: "sky", Hint: "Blue"},
		}
		result, err := parse.ReadWordsFromFile(fileName)
		if result == nil || err != nil {
			t.Fatalf("Unexpected error: %S", err)
		}

		if len(result) != len(expected) {
			t.Fatalf("Incorrect number of results: got %d, want %d", len(result), len(expected))
		}
	})
}
