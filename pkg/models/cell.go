package models

// Cell represents a single square on the crossword puzzle board.
type Cell struct {
	Character  string // The character contained in the cell, if any.
	Filled     bool   // Indicates whether the cell is filled with a character.
	Hint       string // An optional hint associated with this cell.
	Locked     bool   // Indicates whether the cell is locked from being changed.
	UsageCount int    // >1 means it is a intersection. Used for recursion (removal)
}

// NewCell creates a new cell. If a character is provided, the cell is
// marked as filled.
func NewCell(char string, hint string, locked bool) *Cell {
	filled := char != "" // Assuming rune(0) means no character.
	return &Cell{
		Character: char,
		Filled:    filled,
		Hint:      hint,
		Locked:    locked,
	}
}

// NewEmptyCell creates an empty cell with no initial character, no
// hint, and unlocked.
func NewEmptyCell() *Cell {
	return NewCell("", "", false)
}

// SetCharacter sets a character to the cell and marks it as filled.
func (c *Cell) SetCharacter(char string) {
	if !c.Locked {
		c.Character = char
		c.Filled = true
	}
}

// ClearCharacter clears the character from the cell if it's not locked.
func (c *Cell) ClearCharacter() {
	if !c.Locked {
		c.Character = ""
		c.Filled = false
	}
}

// Lock prevents further modifications to the cell.
func (c *Cell) Lock() {
	c.Locked = true
}

// Unlock allows modifications to the cell.
func (c *Cell) Unlock() {
	c.Locked = false
}

// IsEmpty checks if the cell is empty (i.e., does not contain a
// character).
func (c *Cell) IsEmpty() bool {
	return !c.Filled
}

// String returns a string representation of the cell for debugging or
// display purposes.
func (c *Cell) String() string {
	if !c.Filled {
		return "."
	}
	return string(c.Character)
}
