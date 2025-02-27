package models

type Pool struct {
	Words    []string
	ByLength map[int][]string
	WordSet  map[string]bool
}

func NewPool() *Pool {
	return &Pool{
		ByLength: make(map[int][]string),
		WordSet:  make(map[string]bool),
	}
}

// LoadWords loads words into the pool from a given slice of words.

func (p *Pool) LoadWords(words []string) {
	for _, word := range words {
		runes := []rune(word)
		length := len(runes)
		p.Words = append(p.Words, word)
		p.ByLength[length] = append(p.ByLength[length], word)
		p.WordSet[word] = true // Add the word to the set for quick validation
	}
}

// Exists checks if a word is in the pool.
func (p *Pool) Exists(word string) bool {
	_, exists := p.WordSet[word]
	return exists
}
