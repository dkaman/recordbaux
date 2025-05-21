package physical

import "sort"

type binOpt func(*Bin) error

type Bin struct {
	sortFunc sortFunc

	ID      string
	Records []*Record
	Size    int
}

func newBin(id string, size int, opts ...binOpt) (*Bin, error) {
	b := &Bin{
		ID:   id,
		Size: size,
	}

	for _, o := range opts {
		err := o(b)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (b *Bin) Insert(r *Record) *Record {
	b.Records = append(b.Records, r)

	sort.Sort(b)

	if len(b.Records) > b.Size {
		fullShelf := b.Records[0:b.Size]

		last := b.Records[len(b.Records)-1]

		b.Records = fullShelf

		return last
	}

	return nil
}

func (b *Bin) Len() int           { return len(b.Records) }
func (b *Bin) Less(i, j int) bool { return b.sortFunc(b.Records[i], b.Records[j]) }
func (b *Bin) Swap(i, j int)      { b.Records[i], b.Records[j] = b.Records[j], b.Records[i] }

func WithSortFunc(f sortFunc) binOpt {
	return func(b *Bin) error {
		b.sortFunc = f
		return nil
	}
}
