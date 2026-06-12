package vector3

import "github.com/Hedwig7s/vector"

type Serializable[T vector.Number] struct {
	X T
	Y T
	Z T
}

func (m Serializable[T]) Immutable() Vector[T] {
	return Vector[T]{
		x: m.X,
		y: m.Y,
		z: m.Z,
	}
}
