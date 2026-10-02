package vector_test

import (
	"math"
	"testing"

	"github.com/Hedwig7s/vector"
)

func TestFromFloatSaturating(t *testing.T) {
	if got := vector.FromFloatSaturating[uint8](math.NaN()); got != 0 {
		t.Errorf("NaN: got %v", got)
	}
	if got := vector.FromFloatSaturating[uint8](-3.5); got != 0 {
		t.Errorf("negative: got %v", got)
	}
	if got := vector.FromFloatSaturating[uint8](math.Inf(-1)); got != 0 {
		t.Errorf("-Inf: got %v", got)
	}
	if got := vector.FromFloatSaturating[uint8](12.9); got != 12 {
		t.Errorf("truncate: got %v", got)
	}
	if got := vector.FromFloatSaturating[uint8](255); got != 255 {
		t.Errorf("exact max: got %v", got)
	}
	if got := vector.FromFloatSaturating[uint8](300); got != math.MaxUint8 {
		t.Errorf("overflow: got %v", got)
	}
	if got := vector.FromFloatSaturating[uint32](math.Inf(1)); got != math.MaxUint32 {
		t.Errorf("+Inf: got %v", got)
	}
	if got := vector.FromFloatSaturating[uint64](math.Pow(2, 64)); got != math.MaxUint64 {
		t.Errorf("2^64: got %v", got)
	}
	if got := vector.FromFloatSaturating[uint64](math.Pow(2, 63)); got != 1<<63 {
		t.Errorf("2^63: got %v", got)
	}
}
