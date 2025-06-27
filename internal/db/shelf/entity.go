package shelf

import (
	"errors"
	"fmt"
	"slices"

	"github.com/dkaman/recordbaux/internal/db"
	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/record"
	"github.com/google/uuid"
)

var (
	NameIsEmptyErr = errors.New("shelf name cannot be empty")
	ShelfIsFullErr = errors.New("cannot insert more records, shelf is full")
)

type Entity struct {
	ID      db.ID
	Name    string
	BinSize int
	Bins    [][]*bin.Entity
}

func New(name string, binSize int, opts ...option) (*Entity, error) {
	if name == "" {
		return nil, NameIsEmptyErr
	}

	if binSize < 0 {
		return nil, fmt.Errorf("bin size '%d' is not allowed", binSize)
	}

	u := uuid.New()

	e := &Entity{
		ID:      db.ID(u),
		Name:    name,
		BinSize: binSize,
		Bins:    [][]*bin.Entity{},
	}

	for _, o := range opts {
		err := o(e)
		if err != nil {
			return nil, fmt.Errorf("error in shelf constructor option: %w", err)
		}
	}

	return e, nil
}

func (e *Entity) flattenBins() []*bin.Entity {
	return slices.Concat(e.Bins...)
}

func (e *Entity) nextIndex() int {
	return len(e.flattenBins())
}

func (e *Entity) AddRow() *Entity {
	e.Bins = append(e.Bins, []*bin.Entity{})
	return e
}

func (e *Entity) AddBin(size int, sort string) *Entity {
	row := len(e.Bins) - 1
	next := e.nextIndex()
	b, _ := bin.New(indexToLabel(next), size, sort)
	e.Bins[row] = append(e.Bins[row], b)
	return e
}

func (e *Entity) AddBins(n, size int, sort string) *Entity {
	for range n {
		e.AddBin(size, sort)
	}
	return e
}

func (e *Entity) Insert(r *record.Entity) (*record.Entity, error) {
	next := r

	for _, bin := range e.flattenBins() {
		bumped := bin.Insert(next)
		if bumped == nil {
			next = nil
			break
		}
		next = bumped
	}

	if next != nil {
		return next, ShelfIsFullErr
	}

	return nil, nil
}

func (e *Entity) DimX() int {
	l := 0
	for _, row := range e.Bins {
		if len(row) > l {
			l = len(row)
		}
	}

	return l
}

func (e *Entity) DimY() int {
	return len(e.Bins)
}

func (e *Entity) AllBins() []*bin.Entity {
	return e.flattenBins()
}

// indexToLabel converts 0 -> "A", 25 -> "Z", 26 -> "AA", etc.
func indexToLabel(i int) string {
	label := ""
	for i >= 0 {
		rem := i % 26
		label = string('A'+rem) + label
		i = i/26 - 1
	}
	return label
}
