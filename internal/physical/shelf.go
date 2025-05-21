package physical

import (
	"errors"
)

var (
	// sentinel errors
	ShelfCapacityExceededErr = errors.New("attempt was made to insert into a full shelf")
)

// option type to keep shelf construction as minimal as possible
type shelfOpt func(*Shelf) error

type Shelf struct {
	// shape interface to calculate capacity, we can use this for casting
	// later to draw with bubble tea
	Shape    Shape

	// function that will be applied during the sorting process to sort
	// records
	sortFunc sortFunc

	// shelf name for distinction
	Name       string

	// representation of the organized collection
	Bins []*Bin
}

// New is the constructor for the shelf struct. it uses the functional option
// pattern to keep its argument list short, just a name is required to get
// started, but users can specify several options using the option functions
// provided in this package.
func New(name string, opts ...shelfOpt) (*Shelf, error) {
	s := &Shelf{
		Name: name,
	}

	for _, o := range opts {
		err := o(s)
		if err != nil {
			return nil, err
		}
	}

	f := s.sortFunc
	if f == nil {
		f = AlphaByArtist
	}

	if s.Shape != nil {
		n := s.Shape.NumBins()
		sz := s.Shape.BinSize()
		ids := generateBinIDs(n)


		s.Bins = make([]*Bin, n)

		for i, id := range ids {
			b, err := newBin(id, sz, WithSortFunc(f))
			if err != nil {
				return nil, err
			}

			s.Bins[i] = b
		}
	} else {
		s.Bins = make([]*Bin, 1)

		// lol big bin, this should be changed somehow
		b, err := newBin("A", 10000, WithSortFunc(f))
		if err != nil {
			return nil, err
		}

		s.Bins[0] = b
	}

	return s, nil
}

// shelf API

// Insert adds a record to the collection. it should do this by starting at the
// first bin, attempting to insert into that bin, and then cascade the
// reamainder down the bins until there is no remainder left. if there is still
// remainder by the time the last bin is reached, a ShelfCapacityExceededErr is
// returned

func (s *Shelf) Insert(r *Record) (*Record, error) {
	next := r

	for _, bin := range s.Bins {
		bumped := bin.Insert(next)
		if bumped == nil {
			next = nil
			break
		}
		next = bumped
	}

	if next != nil {
		return next, ShelfCapacityExceededErr
	}

	return nil, nil
}

// interface implementations

// needed to implement the list.Model bubble tea componenet interface
func (s *Shelf) FilterValue() string {
	return s.Name
}

// constructor options

func WithShape(s Shape) shelfOpt {
	return func(shelf *Shelf) error {
		shelf.Shape = s
		return nil
	}
}

func WithShelfSortFunc(f sortFunc) shelfOpt {
	return func(s *Shelf) error {
		s.sortFunc = f
		return nil
	}
}

// generateBinIDs returns labels "A", "B", ..., "Z", "AA", etc.
func generateBinIDs(n int) []string {
	ids := make([]string, n)
	for i := 0; i < n; i++ {
		ids[i] = indexToLabel(i)
	}
	return ids
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
