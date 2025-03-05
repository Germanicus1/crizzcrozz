package models

import "fmt"

// Location defines a point in a 2D grid
type Location struct {
	X, Y int
}

// Bounds defines a rectangular region in the crossword puzzle
type Bounds struct {
	TopLeft     Location
	BottomRight Location
}

// NewBounds creates new Bounds given top-left and bottom-right coordinates
func NewBounds(topLeft, bottomRight Location) *Bounds {
	return &Bounds{
		TopLeft:     topLeft,
		BottomRight: bottomRight,
	}
}

func NewBoundsRectangle(width, height int) (*Bounds, error) {
	if width < 1 || height < 1 {
		return nil, fmt.Errorf("width and height must be positive numbers")
	}

	halfSizeX := width / 2
	halfSizeY := height / 2

	return NewBounds(
		Location{X: -halfSizeX, Y: -halfSizeY},
		Location{X: -halfSizeX + width - 1, Y: -halfSizeY + height - 1},
	), nil
}

// Width returns the width of the bounds
func (b *Bounds) Width() int {
	return b.BottomRight.X - b.TopLeft.X + 1
}

// Height returns the height of the bounds
func (b *Bounds) Height() int {
	return b.BottomRight.Y - b.TopLeft.Y + 1
}

// Contains checks if a location is within the bounds
func (b *Bounds) Contains(loc Location) bool {
	return loc.X >= b.TopLeft.X && loc.X <= b.BottomRight.X &&
		loc.Y >= b.TopLeft.Y && loc.Y <= b.BottomRight.Y
}
