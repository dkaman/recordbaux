package physical

type Shape interface {
	NumBins() int
	BinSize() int
}

type Rectangular struct {
	X    int
	Y    int
	Size int
}

func (s *Rectangular) NumBins() int {
	return s.X * s.Y
}

func (s *Rectangular) BinSize() int {
	return s.Size
}

type Irregular struct {
	N    int
	Size int
}

func (s *Irregular) NumBins() int {
	return s.N
}

func (s *Irregular) BinSize() int {
	return s.Size
}
