package shelf

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dkaman/recordbaux/internal/db/bin"
	"github.com/dkaman/recordbaux/internal/db/record"
)

var (
	NameIsEmptyErr  = errors.New("shelf name cannot be empty")
	ShelfIsFullErr  = errors.New("cannot insert more records, shelf is full")
	ShapeIsEmptyErr = errors.New("get shape was called before the shape was initialized, i don't think this can happen")
)

type Shape struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Entity struct {
	Shape []byte `gorm:"type:jsonb"`

	ID      uint `gorm:"primaryKey"`
	Name    string
	BinSize int
	Bins    []*bin.Entity `gorm:"foreignKey:ShelfID;references:ID"`
}

func New(name string, binSize int, opts ...option) (*Entity, error) {
	if name == "" {
		return nil, NameIsEmptyErr
	}

	if binSize < 0 {
		return nil, fmt.Errorf("bin size '%d' is not allowed", binSize)
	}

	e := &Entity{
		Name:    name,
		BinSize: binSize,
		Bins:    []*bin.Entity{},
	}

	for _, o := range opts {
		err := o(e)
		if err != nil {
			return nil, fmt.Errorf("error in shelf constructor option: %w", err)
		}
	}

	return e, nil
}

func (e *Entity) nextIndex() int {
	return len(e.Bins)
}

func (e *Entity) SetShape(s Shape) error {
	d, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("error marshalling shape to json: %w", err)
	}
	e.Shape = d
	return nil
}

func (e *Entity) GetShape() (Shape, error) {
	var s Shape

	if len(e.Shape) == 0 {
		return s, ShapeIsEmptyErr
	}

	err := json.Unmarshal(e.Shape, &s)
	if err != nil {
		return s, fmt.Errorf("error unmarshalling shape data: %w", err)
	}

	return s, nil
}

func (e *Entity) AddBin(size int, sort string) *Entity {
	next := e.nextIndex()
	b, _ := bin.New(indexToLabel(next), size, sort)
	e.Bins = append(e.Bins, b)
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

	for _, bin := range e.Bins {
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

// implementing the tabler interface to change default name so it's not
// entitites
func (e *Entity) TableName() string {
	return "shelves"
}

func (e *Entity) TotalRecords() int {
	total := 0
	for _, b := range e.Bins {
		total += len(b.Records)
	}
	return total
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
