package vector2

import "github.com/Hedwig7s/vector"

type Serializable[T vector.Number] struct {
	X T
	Y T
}

func (m Serializable[T]) Immutable() Vector[T] {
	return Vector[T]{
		x: m.X,
		y: m.Y,
	}
}
