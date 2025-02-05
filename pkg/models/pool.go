package models

type Pool struct {
	Words    []string
	ByLength map[int][]string
}

func NewPool() *Pool {
	return &Pool{
		ByLength: make(map[int][]string),
	}
}

// LoadWords loads words into the pool from a given slice of words.

func (p *Pool) LoadWords(words []string) {
	for _, word := range words {
		length := len(word)
		p.Words = append(p.Words, word)
		p.ByLength[length] = append(p.ByLength[length], word)
	}
}
