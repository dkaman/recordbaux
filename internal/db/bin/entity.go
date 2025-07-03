package bin

import (
	"errors"
	"fmt"
	"sort"

	"github.com/dkaman/recordbaux/internal/db/record"
)

var (
	LabelIsEmptyErr = errors.New("bin label cannot be empty")
)

type Entity struct {
	ID       uint `gorm:"primaryKey"`
	ShelfID  uint `gorm:"index"`
	Label    string
	Size     int
	Records  []*record.Entity `gorm:"foreignKey:BinID;references:ID"`
	SortType string
}

func New(l string, s int, sort string, opts ...option) (*Entity, error) {
	if l == "" {
		return nil, LabelIsEmptyErr
	}

	if s < 0 {
		return nil, fmt.Errorf("size '%d' for bin is invalid", s)
	}

	if _, ok := sortRegistry[sort]; !ok {
		return nil, fmt.Errorf("sort type '%s' not found in registry", sort)
	}

	e := &Entity{
		Label:    l,
		Size:     s,
		SortType: sort,
		Records:  make([]*record.Entity, 0),
	}

	for _, o := range opts {
		err := o(e)
		if err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *Entity) Insert(r *record.Entity) *record.Entity {
	e.Records = append(e.Records, r)

	sort.Sort(e)

	var bumped *record.Entity

	if len(e.Records) > e.Size {
		bumped = e.Records[e.Size]
		e.Records = e.Records[:e.Size]
	}

	for i, rec := range e.Records {
		rec.Position = i
	}

	if bumped != nil {
		bumped.BinID = 0
		bumped.Position = 0
	}

	return bumped
}

// implementing the tabler interface to change default name so it's not
// entitites
func (e *Entity) TableName() string {
	return "bins"
}

// sort interface defs
func (e *Entity) Len() int { return len(e.Records) }
func (e *Entity) Less(i, j int) bool {
	// we check this during construction so we can be confident it exists
	f := sortRegistry[e.SortType]
	return f(e.Records[i], e.Records[j])
}
func (e *Entity) Swap(i, j int) { e.Records[i], e.Records[j] = e.Records[j], e.Records[i] }
