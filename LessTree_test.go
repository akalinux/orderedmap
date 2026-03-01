package omap

import (
	"cmp"
	"testing"
)

func TestLessIndex(t *testing.T) {
	Slice := []KvSet[int, int]{{1, 1}}
	Expected := func(k, cidx, coffset int) {
		idx, offset := LessIndex(k, cmp.Less, Slice)

		if idx != cidx || offset != coffset {
			t.Fatalf("Expected index: %d, got: %d, expected offset: %d, got: %d", cidx, idx, coffset, offset)
		}
	}
	Expected(0, 0, -1)
	Expected(1, 0, 0)
	Expected(2, 0, 1)

	t.Logf("%d", getMid(3))
}
