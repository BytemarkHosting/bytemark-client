package util

import (
	"bytemark.co.uk/client/lib"
	"fmt"
	"strings"
)

// DiscSpecError represents an error during parse.
type DiscSpecError struct {
	Position  int
	Character rune
}

func (e *DiscSpecError) Error() string {
	return fmt.Sprintf("Disc spec error: Unexpected %c at character %d.", e.Character, e.Position)
}

func ParseDiscSpec(spec string) (*lib.Disc, error) {
	bits := strings.Split(spec, ":")
	size, err := ParseSize(bits[len(bits)-1])
	if err != nil {
		return nil, err
	}
	switch {
	case len(bits) >= 4:
		return nil, &DiscSpecError{}
	case len(bits) == 3:
		return &lib.Disc{Label: bits[0], StorageGrade: bits[1], Size: size}, nil
	case len(bits) == 2:
		return &lib.Disc{StorageGrade: bits[0], Size: size}, nil
	case len(bits) == 1:
		return &lib.Disc{Size: size}, nil
	case len(bits) == 0:
		return nil, &DiscSpecError{}
	}
	return nil, nil
}