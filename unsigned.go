package vector

// Unsigned is the set of unsigned integer types usable by the coord packages
type Unsigned interface {
	uint8 | uint16 | uint | uint32 | uint64
}

// FromFloatSaturating converts f to T, clamping to the range of T instead of
// relying on Go's implementation specific out of range float conversion.
// NaN and anything at or below 0 give 0, anything at or above the max of T
// gives the max of T. Fractions are truncated.
func FromFloatSaturating[T Unsigned](f float64) T {
	if f != f || f <= 0 {
		return 0
	}
	if max := ^T(0); f >= float64(max) {
		return max
	}
	return T(f)
}
