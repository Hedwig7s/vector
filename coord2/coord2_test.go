package coord2_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math"
	"testing"

	"github.com/Hedwig7s/vector"
	"github.com/Hedwig7s/vector/coord2"
)

func TestSubWrapsAndSafeAlternatives(t *testing.T) {
	a := coord2.Fill[uint8](3)
	b := coord2.Fill[uint8](5)

	if got := a.Sub(b); got != coord2.Fill[uint8](254) {
		t.Errorf("Sub wrap: got %v", got)
	}
	if got := a.SubSaturating(b); got != coord2.Zero[uint8]() {
		t.Errorf("SubSaturating: got %v", got)
	}
	if got := b.SubSaturating(a); got != coord2.Fill[uint8](2) {
		t.Errorf("SubSaturating: got %v", got)
	}
	if got := a.AbsDiff(b); got != coord2.Fill[uint8](2) {
		t.Errorf("AbsDiff: got %v", got)
	}
	if got := b.AbsDiff(a); got != coord2.Fill[uint8](2) {
		t.Errorf("AbsDiff reversed: got %v", got)
	}
}

func TestAddSaturating(t *testing.T) {
	a := coord2.Fill[uint8](250)
	if got := a.AddSaturating(coord2.Fill[uint8](10)); got != coord2.Fill[uint8](255) {
		t.Errorf("got %v", got)
	}
	if got := a.Add(coord2.Fill[uint8](10)); got != coord2.Fill[uint8](4) {
		t.Errorf("Add should wrap, got %v", got)
	}
}

func TestDistanceEitherOrder(t *testing.T) {
	a := coord2.Fill[uint8](10)
	b := coord2.Zero[uint8]()
	want := math.Sqrt(200)
	if got := a.Distance(b); math.Abs(got-want) > 1e-9 {
		t.Errorf("a to b: got %v want %v", got, want)
	}
	if got := b.Distance(a); math.Abs(got-want) > 1e-9 {
		t.Errorf("b to a: got %v want %v", got, want)
	}
	if got := a.DistanceSquared(b); got != 200 {
		t.Errorf("squared: got %v", got)
	}
}

func TestLerp(t *testing.T) {
	a := coord2.Fill[uint8](200)
	b := coord2.Fill[uint8](100)
	if got := coord2.Lerp(a, b, 0.5); got != coord2.Fill[uint8](150) {
		t.Errorf("descending lerp: got %v", got)
	}
	if got := coord2.Lerp(a, b, 2); got != coord2.Zero[uint8]() {
		t.Errorf("extrapolation below 0 should clamp, got %v", got)
	}
	if got := coord2.LerpClamped(a, b, 2); got != b {
		t.Errorf("LerpClamped: got %v", got)
	}
}

func TestScaleClamps(t *testing.T) {
	v := coord2.Fill[uint8](100)
	if got := v.Scale(3); got != coord2.Fill[uint8](255) {
		t.Errorf("overflow: got %v", got)
	}
	if got := v.Scale(-1); got != coord2.Zero[uint8]() {
		t.Errorf("negative scale: got %v", got)
	}
	if got := v.Scale(0.5); got != coord2.Fill[uint8](50) {
		t.Errorf("half: got %v", got)
	}
	if got := v.DivByConstant(4); got != coord2.Fill[uint8](25) {
		t.Errorf("div: got %v", got)
	}
	if got := v.DivByConstant(0); got != coord2.Fill[uint8](255) {
		t.Errorf("div by zero: got %v", got)
	}
}

func TestMidpointNoOverflow(t *testing.T) {
	got := coord2.Midpoint(coord2.Fill[uint8](250), coord2.Fill[uint8](254))
	if got != coord2.Fill[uint8](252) {
		t.Errorf("got %v", got)
	}
}

func TestLengthNoOverflow(t *testing.T) {
	if got := coord2.Fill[uint8](255).LengthSquared(); got != 130050 {
		t.Errorf("got %v", got)
	}
}

func TestMinMaxClamp(t *testing.T) {
	a := coord2.New[uint16](1, 9)
	b := coord2.New[uint16](5, 2)
	if got, want := coord2.Min(a, b), coord2.New[uint16](1, 2); got != want {
		t.Errorf("Min: got %v want %v", got, want)
	}
	if got, want := coord2.Max(a, b), coord2.New[uint16](5, 9); got != want {
		t.Errorf("Max: got %v want %v", got, want)
	}
	if got, want := a.Clamp(2, 6), coord2.New[uint16](2, 6); got != want {
		t.Errorf("Clamp: got %v want %v", got, want)
	}
	if a.MinComponent() != 1 || a.MaxComponent() != 9 {
		t.Errorf("min/max component: %v %v", a.MinComponent(), a.MaxComponent())
	}
}

func TestComponentsAndSetters(t *testing.T) {
	v := coord2.New[uint32](1, 2)
	for i, want := range []uint32{1, 2} {
		if got := v.Component(i); got != want {
			t.Errorf("Component(%d): got %v want %v", i, got, want)
		}
	}
	if got := v.SetX(9).X(); got != 9 {
		t.Errorf("SetX: got %v", got)
	}
	if v.X() != 1 {
		t.Errorf("SetX mutated the original")
	}
	if got := coord2.FromArray([]uint32{7}); got.X() != 7 || got.Y() != 0 {
		t.Errorf("FromArray: got %v", got)
	}
	if !coord2.Zero[uint32]().IsZero() || v.IsZero() {
		t.Errorf("IsZero wrong")
	}
}

func TestAverage(t *testing.T) {
	got := coord2.Average([]coord2.Uint8{coord2.Fill[uint8](250), coord2.Fill[uint8](254), coord2.Fill[uint8](255)})
	if got != coord2.Fill[uint8](253) {
		t.Errorf("got %v", got)
	}
	if got := coord2.Average[uint8](nil); !got.IsZero() {
		t.Errorf("empty average: got %v", got)
	}
}

func TestJSON(t *testing.T) {
	want := coord2.Fill[uint64](math.MaxUint64)
	data, err := json.Marshal(want)
	if err != nil {
		t.Fatal(err)
	}
	var got coord2.Uint64
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("round trip lost precision: got %v want %v (json %s)", got, want, data)
	}

	var neg coord2.Uint8
	if err := json.Unmarshal([]byte(`{"x": -1}`), &neg); err == nil {
		t.Errorf("negative component should be rejected")
	}
}

func TestWriteRead(t *testing.T) {
	for _, endian := range []binary.ByteOrder{binary.BigEndian, binary.LittleEndian} {
		t.Run(endian.String(), func(t *testing.T) {
			roundTrip(t, coord2.New[uint8](1, 2), endian)
			roundTrip(t, coord2.New[uint16](1, 2), endian)
			roundTrip(t, coord2.New[uint32](1, 2), endian)
			roundTrip(t, coord2.New[uint64](1, 2), endian)
		})
	}
}

func roundTrip[T vector.Unsigned](t *testing.T, v coord2.Coord[T], endian binary.ByteOrder) {
	t.Helper()
	var buf bytes.Buffer
	if err := v.Write(&buf, endian); err != nil {
		t.Fatal(err)
	}
	got, err := coord2.Read[T](&buf, endian)
	if err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("got %v want %v", got, v)
	}
}

func TestArray(t *testing.T) {
	arr := coord2.Array[uint8]{coord2.New[uint8](1, 5), coord2.New[uint8](4, 2), coord2.New[uint8](2, 7)}

	if got, want := arr.Sum(), coord2.New[uint8](7, 14); got != want {
		t.Errorf("Sum: got %v want %v", got, want)
	}
	lo, hi := arr.Bounds()
	if lo != coord2.New[uint8](1, 2) || hi != coord2.New[uint8](4, 7) {
		t.Errorf("Bounds: got %v %v", lo, hi)
	}
	if got := arr.Add(coord2.Fill[uint8](1)); got[0] != coord2.New[uint8](2, 6) || arr[0] != coord2.New[uint8](1, 5) {
		t.Errorf("Add should not mutate: %v %v", got, arr)
	}
	arr.ScaleInplace(2)
	if arr[0] != coord2.New[uint8](2, 10) {
		t.Errorf("ScaleInplace: got %v", arr[0])
	}

	defer func() {
		if recover() == nil {
			t.Errorf("Bounds of empty array should panic")
		}
	}()
	coord2.Array[uint8]{}.Bounds()
}

func TestSerializable(t *testing.T) {
	s := coord2.Serializable[uint8]{X: 4}
	if s.Immutable().X() != 4 {
		t.Errorf("got %v", s.Immutable())
	}
}
