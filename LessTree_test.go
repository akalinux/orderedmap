package omap

import (
	"cmp"
	"testing"
)

func TestLessIndex(t *testing.T) {
	Slice := []KvSet[int, int]{{1, 1}}
	Expected := func(k, cidx, coffset int) {
		idx, offset := LessIndex(k, cmp.Less, Slice)

		t.Logf("Testing Lookup of: %d, idx: %d, offset: %d", k, cidx, coffset)
		if idx != cidx || offset != coffset {
			t.Fatalf("Expected index: %d, got: %d, expected offset: %d, got: %d", cidx, idx, coffset, offset)
		}
	}
	Expected(0, 0, -1)
	Expected(1, 0, 0)
	Expected(2, 0, 1)
	Slice = []KvSet[int, int]{{0, 0}, {1, 1}, {2, 2}}
	Expected(1, 1, 0)
	Expected(2, 2, 0)
	Expected(0, 0, 0)
	Expected(-1, 0, -1)
	Expected(3, 2, 1)
	Slice = []KvSet[int, int]{{0, 0}, {3, 3}, {5, 5}}
	Expected(2, 0, 1)
	Expected(4, 2, -1)
	Slice = []KvSet[int, int]{{0, 0}, {3, 3}, {5, 5}, {7, 7}}
	Expected(4, 2, -1)
	Expected(6, 3, -1)
	Expected(7, 3, 0)

}

func TestUpgradedGetIndex(t *testing.T) {

	Slice := []KvSet[int, int]{{0, 0}, {1, 1}, {2, 2}}
	Expected := func(k, cidx, coffset int) {
		idx, offset := GetIndex(k, cmp.Compare, Slice)

		t.Logf("Checking key: %d", k)
		if idx != cidx || offset != coffset {
			t.Fatalf("Expected index: %d, got: %d, expected offset: %d, got: %d", cidx, idx, coffset, offset)
		}
	}
	Expected(0, 0, 0)
	Expected(1, 1, 0)
	Expected(2, 2, 0)

}
