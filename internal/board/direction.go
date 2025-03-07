package board

// Direction indicates the direction in which a word is placed on the
// crossword board. It is used to specify whether a word runs
// horizontally or vertically.
type Direction int

const (
	// Across (0) indicates a horizontal placement, from left to right.
	Across Direction = iota
	// Down (1) indicates a vertical placement, from top to bottom.
	Down
)
