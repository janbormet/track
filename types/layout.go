package types

import (
	"encoding"
)

type Layout interface {
	ToIntervals(encoding.BinaryMarshaler) (Intervals, error)
	FromIntervals(Intervals) (encoding.BinaryMarshaler, error)
}
