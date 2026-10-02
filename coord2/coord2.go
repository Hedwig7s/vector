package coord2

import (
	"encoding/json"
	"fmt"
	"math"

	"github.com/Hedwig7s/vector"
	"github.com/Hedwig7s/vector/vector2"
)

// Coord is an immutable, unsigned, 2-component tuple meant for positions,
// sizes, scales and grid coordinates. Unlike vector2.Vector it has no concept
// of direction, so there is no negation (Flip), Normalized, Dot, Cross, Angle,
// Perpendicular, or directional constants.
//
// Arithmetic that can leave the range of T follows Go's integer rules
// (Add, Sub, MultByVector wrap around). Use the Saturating variants or the
// float64-based operations (Scale, Lerp, Midpoint, Sqrt, ...) when you need
// clamping instead: those convert back to T by clamping to [0, max(T)].
type Coord[T vector.Unsigned] struct {
	x T
	y T
}

type (
	Uint8  = Coord[uint8]
	Uint16 = Coord[uint16]
	Uint32 = Coord[uint32]
	Uint64 = Coord[uint64]
	Uint   = Coord[uint]
)

// New creates a new coord from its components
func New[T vector.Unsigned](x, y T) Coord[T] {
	return Coord[T]{
		x: x,
		y: y,
	}
}

// Fill creates a coord where each component is equal to v
func Fill[T vector.Unsigned](v T) Coord[T] {
	return Coord[T]{
		x: v,
		y: v,
	}
}

func Zero[T vector.Unsigned]() Coord[T] {
	return Fill[T](0)
}

func One[T vector.Unsigned]() Coord[T] {
	return Fill[T](1)
}

// Average returns the component wise mean of coords (rounded down). It is
// computed in float64, so it can not overflow. The average of no coords is zero.
func Average[T vector.Unsigned](coords []Coord[T]) Coord[T] {
	if len(coords) == 0 {
		return Zero[T]()
	}
	var sx float64
	var sy float64
	for _, c := range coords {
		sx += float64(c.x)
		sy += float64(c.y)
	}
	n := float64(len(coords))
	return Coord[T]{
		x: vector.FromFloatSaturating[T](sx / n),
		y: vector.FromFloatSaturating[T](sy / n),
	}
}

// Lerp linearly interpolates between a and b by t. Values outside of the range
// of T (possible when t is outside of 0 to 1) are clamped.
func Lerp[T vector.Unsigned](a, b Coord[T], t float64) Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](float64(a.x) + (float64(b.x)-float64(a.x))*t),
		y: vector.FromFloatSaturating[T](float64(a.y) + (float64(b.y)-float64(a.y))*t),
	}
}

// LerpClamped is Lerp with t clamped from 0 to 1
func LerpClamped[T vector.Unsigned](a, b Coord[T], t float64) Coord[T] {
	return Lerp(a, b, vector.Clamp(t, 0, 1))
}

// Min returns a coord where each component is the smallest of the two coords
func Min[T vector.Unsigned](a, b Coord[T]) Coord[T] {
	return Coord[T]{
		x: min(a.x, b.x),
		y: min(a.y, b.y),
	}
}

// Max returns a coord where each component is the largest of the two coords
func Max[T vector.Unsigned](a, b Coord[T]) Coord[T] {
	return Coord[T]{
		x: max(a.x, b.x),
		y: max(a.y, b.y),
	}
}

func MaxX[T vector.Unsigned](a, b Coord[T]) T {
	return max(a.x, b.x)
}

func MaxY[T vector.Unsigned](a, b Coord[T]) T {
	return max(a.y, b.y)
}

func MinX[T vector.Unsigned](a, b Coord[T]) T {
	return min(a.x, b.x)
}

func MinY[T vector.Unsigned](a, b Coord[T]) T {
	return min(a.y, b.y)
}

// Midpoint returns the point halfway between a and b (rounded down). It can
// not overflow.
func Midpoint[T vector.Unsigned](a, b Coord[T]) Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T]((float64(a.x) + float64(b.x)) * 0.5),
		y: vector.FromFloatSaturating[T]((float64(a.y) + float64(b.y)) * 0.5),
	}
}

// FromArray builds a coord from data to the best of its ability. Missing
// components are left as 0, extra values are ignored.
func FromArray[T vector.Unsigned](data []T) Coord[T] {
	v := Coord[T]{}
	if len(data) > 0 {
		v.x = data[0]
	}

	if len(data) > 1 {
		v.y = data[1]
	}
	return v
}

func (v Coord[T]) ToArr() []T {
	return []T{v.x, v.y}
}

func (v Coord[T]) ToFixedArr() [2]T {
	return [2]T{v.x, v.y}
}

// MarshalJSON encodes the components as numbers without going through float64,
// so uint64 values keep full precision.
func (v Coord[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		X T `json:"x"`
		Y T `json:"y"`
	}{
		X: v.x,
		Y: v.y,
	})
}

// UnmarshalJSON decodes the components. Negative or fractional numbers are
// rejected with an error.
func (v *Coord[T]) UnmarshalJSON(data []byte) error {
	aux := &struct {
		X T `json:"x"`
		Y T `json:"y"`
	}{}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}
	v.x = aux.X
	v.y = aux.Y
	return nil
}

func (v Coord[T]) Format(format string) string {
	return fmt.Sprintf(format, v.x, v.y)
}

func (v Coord[T]) MinComponent() T {
	return min(v.x, v.y)
}

func (v Coord[T]) MaxComponent() T {
	return max(v.x, v.y)
}

// IsZero returns true if all components are 0
func (v Coord[T]) IsZero() bool {
	return v.x == 0 && v.y == 0
}

// ToInt converts to a signed vector. Components above the range of int wrap.
func (v Coord[T]) ToInt() vector2.Vector[int] {
	return vector2.New(int(v.x), int(v.y))
}

// ToInt64 converts to a signed vector. Components above the range of int64 wrap.
func (v Coord[T]) ToInt64() vector2.Vector[int64] {
	return vector2.New(int64(v.x), int64(v.y))
}

func (v Coord[T]) ToFloat64() vector2.Vector[float64] {
	return vector2.New(float64(v.x), float64(v.y))
}

func (v Coord[T]) ToFloat32() vector2.Vector[float32] {
	return vector2.New(float32(v.x), float32(v.y))
}

func (v Coord[T]) X() T {
	return v.x
}

// SetX changes the x component of the coord
func (v Coord[T]) SetX(newX T) Coord[T] {
	return Coord[T]{
		x: newX,
		y: v.y,
	}
}

func (v Coord[T]) Y() T {
	return v.y
}

// SetY changes the y component of the coord
func (v Coord[T]) SetY(newY T) Coord[T] {
	return Coord[T]{
		x: v.x,
		y: newY,
	}
}

func (v Coord[T]) YX() Coord[T] {
	return Coord[T]{x: v.y, y: v.x}
}

// Add is component wise addition. Wraps on overflow.
func (v Coord[T]) Add(other Coord[T]) Coord[T] {
	return Coord[T]{
		x: v.x + other.x,
		y: v.y + other.y,
	}
}

// AddSaturating is component wise addition clamped to the max of T
func (v Coord[T]) AddSaturating(other Coord[T]) Coord[T] {
	return Coord[T]{
		x: addSat(v.x, other.x),
		y: addSat(v.y, other.y),
	}
}

func (v Coord[T]) AddX(other T) Coord[T] {
	return Coord[T]{
		x: v.x + other,
		y: v.y,
	}
}

func (v Coord[T]) AddY(other T) Coord[T] {
	return Coord[T]{
		x: v.x,
		y: v.y + other,
	}
}

// Sub is component wise subtraction. Wraps around (like any Go unsigned
// integer) when a component of other is larger. Prefer SubSaturating or AbsDiff
// when that is not what you want.
func (v Coord[T]) Sub(other Coord[T]) Coord[T] {
	return Coord[T]{
		x: v.x - other.x,
		y: v.y - other.y,
	}
}

// SubSaturating is component wise subtraction clamped to 0
func (v Coord[T]) SubSaturating(other Coord[T]) Coord[T] {
	return Coord[T]{
		x: subSat(v.x, other.x),
		y: subSat(v.y, other.y),
	}
}

// AbsDiff returns the component wise absolute difference |v - other|
func (v Coord[T]) AbsDiff(other Coord[T]) Coord[T] {
	return Coord[T]{
		x: absDiff(v.x, other.x),
		y: absDiff(v.y, other.y),
	}
}

// Midpoint returns the point halfway between v and o (rounded down)
func (v Coord[T]) Midpoint(o Coord[T]) Coord[T] {
	return Midpoint(v, o)
}

// Clamp clamps each component between lo and hi
func (v Coord[T]) Clamp(lo, hi T) Coord[T] {
	return Coord[T]{
		x: min(max(v.x, lo), hi),
		y: min(max(v.y, lo), hi),
	}
}

// Mod is the floating point remainder of each component divided by t
// (math.Mod). A t of 0 gives 0.
func (v Coord[T]) Mod(t float64) Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](math.Mod(float64(v.x), t)),
		y: vector.FromFloatSaturating[T](math.Mod(float64(v.y), t)),
	}
}

// Scale multiplies each component by t. The result is clamped to the range of
// T, so a negative t gives 0.
func (v Coord[T]) Scale(t float64) Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](float64(v.x) * t),
		y: vector.FromFloatSaturating[T](float64(v.y) * t),
	}
}

// DivByConstant divides each component by t. The result is clamped to the
// range of T, so dividing by 0 gives the max of T (or 0 for a 0 component).
func (v Coord[T]) DivByConstant(t float64) Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](float64(v.x) / t),
		y: vector.FromFloatSaturating[T](float64(v.y) / t),
	}
}

// MultByVector is component wise multiplication (Hadamard product). Wraps on
// overflow.
func (v Coord[T]) MultByVector(o Coord[T]) Coord[T] {
	return Coord[T]{
		x: v.x * o.x,
		y: v.y * o.y,
	}
}

// DivByVector is component wise integer division. Panics if a component of o
// is 0.
func (v Coord[T]) DivByVector(o Coord[T]) Coord[T] {
	return Coord[T]{
		x: v.x / o.x,
		y: v.y / o.y,
	}
}

// LengthSquared is the squared distance from the origin, computed in float64
func (v Coord[T]) LengthSquared() float64 {
	return float64(v.x)*float64(v.x) + float64(v.y)*float64(v.y)
}

// Length is the distance from the origin
func (v Coord[T]) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// DistanceSquared is the squared euclidean distance between two points. Safe
// for either argument being larger.
func (v Coord[T]) DistanceSquared(other Coord[T]) float64 {
	dx := float64(other.x) - float64(v.x)
	dy := float64(other.y) - float64(v.y)
	return dx*dx + dy*dy
}

// Distance is the euclidean distance between two points
func (v Coord[T]) Distance(other Coord[T]) float64 {
	return math.Sqrt(v.DistanceSquared(other))
}

// Sqrt returns the square root (rounded down) of each component
func (v Coord[T]) Sqrt() Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](math.Sqrt(float64(v.x))),
		y: vector.FromFloatSaturating[T](math.Sqrt(float64(v.y))),
	}
}

// Log returns the natural logarithm (0 for components below 1) of each component
func (v Coord[T]) Log() Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](math.Log(float64(v.x))),
		y: vector.FromFloatSaturating[T](math.Log(float64(v.y))),
	}
}

// Log10 returns the decimal logarithm (0 for components below 1) of each component
func (v Coord[T]) Log10() Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](math.Log10(float64(v.x))),
		y: vector.FromFloatSaturating[T](math.Log10(float64(v.y))),
	}
}

// Log2 returns the binary logarithm (0 for components below 1) of each component
func (v Coord[T]) Log2() Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](math.Log2(float64(v.x))),
		y: vector.FromFloatSaturating[T](math.Log2(float64(v.y))),
	}
}

// Exp returns the e**x (clamped to the max of T) of each component
func (v Coord[T]) Exp() Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](math.Exp(float64(v.x))),
		y: vector.FromFloatSaturating[T](math.Exp(float64(v.y))),
	}
}

// Exp2 returns the 2**x (clamped to the max of T) of each component
func (v Coord[T]) Exp2() Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](math.Exp2(float64(v.x))),
		y: vector.FromFloatSaturating[T](math.Exp2(float64(v.y))),
	}
}

// Expm1 returns the e**x - 1 (clamped to the max of T) of each component
func (v Coord[T]) Expm1() Coord[T] {
	return Coord[T]{
		x: vector.FromFloatSaturating[T](math.Expm1(float64(v.x))),
		y: vector.FromFloatSaturating[T](math.Expm1(float64(v.y))),
	}
}

func (v Coord[T]) Values() (T, T) {
	return v.x, v.y
}

func (v Coord[T]) Component(index int) T {
	switch index {
	case 0:
		return v.x

	case 1:
		return v.y

	default:
		panic(fmt.Errorf("invalid index: %d", index))
	}
}

func addSat[T vector.Unsigned](a, b T) T {
	if sum := a + b; sum >= a {
		return sum
	}
	return ^T(0)
}

func subSat[T vector.Unsigned](a, b T) T {
	if a > b {
		return a - b
	}
	return 0
}

func absDiff[T vector.Unsigned](a, b T) T {
	if a > b {
		return a - b
	}
	return b - a
}
