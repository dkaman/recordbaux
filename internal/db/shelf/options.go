package shelf

import (
	"errors"
)

var (
	InvalidShapeDimensionsErr = errors.New("a shape cannot have negative dimension")
)

type option func(*Entity) error

func WithShapeRect(x, y, size int, sort string) option {
	return func(e *Entity) error {
		if x < 0 || y < 0 {
			return InvalidShapeDimensionsErr
		}

		s := Shape{
			X: x,
			Y: y,
		}

		err := e.SetShape(s)
		if err != nil {
			return err
		}

		e.AddBins(x * y, size, sort)

		return nil
	}
}
