package coord4

import "github.com/Hedwig7s/vector"

type Serializable[T vector.Unsigned] struct {
	X T
	Y T
	Z T
	W T
}

func (m Serializable[T]) Immutable() Coord[T] {
	return Coord[T]{
		x: m.X,
		y: m.Y,
		z: m.Z,
		w: m.W,
	}
}
