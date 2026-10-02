package coord2

import (
	"errors"

	"github.com/Hedwig7s/vector"
)

type Array[T vector.Unsigned] []Coord[T]

type (
	Uint8Array  = Array[uint8]
	Uint16Array = Array[uint16]
	Uint32Array = Array[uint32]
	Uint64Array = Array[uint64]
	UintArray   = Array[uint]
)

// Add returns a new array with other added to each element (wraps on overflow)
func (a Array[T]) Add(other Coord[T]) (out Array[T]) {
	out = make(Array[T], len(a))
	for i, v := range a {
		out[i] = v.Add(other)
	}
	return
}

func (a Array[T]) AddInplace(other Coord[T]) Array[T] {
	for i, v := range a {
		a[i] = v.Add(other)
	}
	return a
}

// Scale returns a new array with each element scaled by t (see Coord.Scale)
func (a Array[T]) Scale(t float64) (out Array[T]) {
	out = make(Array[T], len(a))
	for i, v := range a {
		out[i] = v.Scale(t)
	}
	return
}

func (a Array[T]) ScaleInplace(t float64) Array[T] {
	for i, v := range a {
		a[i] = v.Scale(t)
	}
	return a
}

// Sum adds every element together (wraps on overflow)
func (a Array[T]) Sum() (sum Coord[T]) {
	for _, v := range a {
		sum = sum.Add(v)
	}
	return
}

// Bounds returns the min and max points of an AABB encompassing every element
func (a Array[T]) Bounds() (Coord[T], Coord[T]) {
	if len(a) == 0 {
		panic(errors.New("can not compute bounds from 0 coord elements"))
	}

	minV := a[0]
	maxV := a[0]

	for _, v := range a[1:] {
		minV = Min(minV, v)
		maxV = Max(maxV, v)
	}

	return minV, maxV
}
