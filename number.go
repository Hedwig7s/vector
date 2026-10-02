package vector

// Number is every type the signed vector packages (vector1 to vector4)
// support. Unsigned integers are deliberately not part of it: the vector
// packages need negation (Flip, Left, Down, Perpendicular, ...), which can not
// be expressed for unsigned types. See Unsigned and the coord packages.
type Number interface {
	int8 | int16 | int | int32 | int64 | float32 | float64
}

type Vector interface {
	Length() float64
	LengthSquared() float64
}
