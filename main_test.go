package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/Germanicus1/crizzcrozz/internal/parser"
)

func TestSortWordsByLength(t *testing.T) {
	words := []string{"zoo", "österreich", "ärger", "überraschung"}
	result := sortWordsByLength(words)
	want := []string{"überraschung", "österreich", "ärger", "zoo"}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Incorrect result, got: %s, want: %s", result, want)
	}
}

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

func TestReadWordsFromFile_Success(t *testing.T) {
	// Prepare a mock CSV file.
	content := "word,hint\napple,Fruit\nsky,Blue"
	fileName, cleanup, err := createMockCSVFile(content)
	if err != nil {
		t.Fatalf("Failed to create mock CSV file: %s", err)
	}
	defer cleanup()

	// Test the function.
	result, err := parser.ReadWordsFromFile(fileName)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	// Debug output
	for _, entry := range result {
		t.Logf("Read entry: %+v", entry)
	}

	// Expected result.
	expected := []*wordsAndHints{
		{Word: "apple", Hint: "Fruit"},
		{Word: "sky", Hint: "Blue"},
	}

	// Compare the results.
	if len(result) != len(expected) {
		t.Fatalf("Incorrect number of results: got %d, want %d", len(result), len(expected))
	}
	for i, r := range result {
		if r.Word != expected[i].Word || r.Hint != expected[i].Hint {
			t.Errorf("Incorrect result at index %d, got: %+v, want: %+v", i, r, expected[i])
		}
	}
}

func TestReadWordsFromFile_Error(t *testing.T) {
	// Use a non-existing file name.
	fileName := "nonexistent.csv"
	_, err := readWordsFromFile(fileName)
	if err == nil {
		t.Fatalf("Expected an error for non-existing file, but got none")
	}
}

// Test readWordsFromFile with a malformed CSV
func TestReadWordsFromFile_FailedUnmarshal(t *testing.T) {
	content := "word,hint\napple,Fruit\nsky" // Malformed line missing one column
	fileName, cleanup, err := createMockCSVFile(content)
	if err != nil {
		t.Fatalf("Failed to create mock CSV file: %s", err)
	}
	defer cleanup()

	_, err = readWordsFromFile(fileName)
	if err == nil {
		t.Fatal("Expected an error due to malformed CSV, but got none")
	}
}
