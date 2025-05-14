package physical

import (
	"fmt"
	"sort"
)

type SortField int

const (
	ByCatalogNumber SortField = iota
	ByTitle
)

type Shelf struct {
	Name       string
	Bins       []Bin
	BinSize    int
	SortBy     SortField
	BinLetters []string // e.g. ["A", "B", "C"]
}

func NewShelf(name string, numBins, binSize int) *Shelf {
	letters := make([]string, numBins)
	for i := 0; i < numBins; i++ {
		letters[i] = string('A' + i)
	}
	bins := make([]Bin, numBins)
	for i, letter := range letters {
		bins[i] = Bin{ID: letter}
	}
	return &Shelf{
		Name:       name,
		Bins:       bins,
		BinSize:    binSize,
		SortBy:     ByCatalogNumber,
		BinLetters: letters,
	}
}

func (s *Shelf) Insert(r Record) {
	all := s.flatten()
	all = append(all, r)

	// Sort
	sort.Slice(all, func(i, j int) bool {
		switch s.SortBy {
		case ByTitle:
			return all[i].Title < all[j].Title
		default:
			return all[i].CatalogNumber < all[j].CatalogNumber
		}
	})

	// Redistribute
	idx := 0
	for i := range s.Bins {
		s.Bins[i].Records = nil
		for j := 0; j < s.BinSize && idx < len(all); j++ {
			s.Bins[i].Records = append(s.Bins[i].Records, all[idx])
			idx++
		}
	}
}

func (s *Shelf) flatten() []Record {
	var all []Record
	for _, bin := range s.Bins {
		all = append(all, bin.Records...)
	}
	return all
}

func (s *Shelf) GetCoordinates(catalog string) (binID string, index int, found bool) {
	for _, bin := range s.Bins {
		for idx, r := range bin.Records {
			if r.CatalogNumber == catalog {
				return bin.ID, idx, true
			}
		}
	}
	return "", -1, false
}

func (s *Shelf) DebugPrint() {
	for _, bin := range s.Bins {
		fmt.Printf("Bin %s:\n", bin.ID)
		for _, r := range bin.Records {
			fmt.Printf("  - [%s] %s\n", r.CatalogNumber, r.Title)
		}
	}
}

func (s *Shelf) FilterValue() string {
	return s.Name
}
