package slowmgr

import (
	"errors"
	"io"
)

type Parser interface {
	Parse(r io.Reader) (*Slow, error)
}

// DefaultParser reads from a file and detect slow queries
type DefaultParser struct{}

func (p *DefaultParser) Parse(r io.Reader) (*Slow, error) {
	return nil, errors.New("NOT IMPLEMENTED")
}
