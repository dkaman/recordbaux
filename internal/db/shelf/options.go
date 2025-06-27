package shelf

type option func(*Entity) error

func WithShapeRect(x, y, size int, sort string) option {
	return func(e *Entity) error {
		for range y {
			e.AddRow().AddBins(x, size, sort)
		}
		return nil
	}
}

func WithShapeIrregular(n, size int, sort string) option {
	return func(e *Entity) error {
		e.AddRow().AddBins(n, size, sort)
		return nil
	}
}
