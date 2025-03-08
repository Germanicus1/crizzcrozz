package models

// wordsAndHints defines a struct for mapping words and their hints from
// a CSV file.
type WordsAndHints struct {
	Word string `csv:"word"`
	Hint string `csv:"hint"`
}
