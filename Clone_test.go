package omap

import (
	"cmp"
	"testing"
)

func TestClone(t *testing.T) {
	var s OrderedMap[int, int]

	s = NewSliceTree[int, int](5, cmp.Compare).ToTs()

	TestClone := func() {
		for i := range 50 {
			s.Put(i, i)
		}
		c := s.Clone()
		for k, v := range s.All() {
			if lv, ok := c.Get(k); !ok || lv != v {
				t.Fatalf("Expected ok: true for key: %d, got %v, expected Value: %d, got: %v, ", k, ok, v, lv)
			}
		}
		if s.Size() != c.Size() {
			t.Fatalf("Size missmatch, expected: %d, got %d", s.Size(), c.Size())
		}
	}
	t.Logf("Testing SliceTree")
	TestClone()
	t.Logf("Testing Centertree")
	s = NewCenterTree[int, int](5, cmp.Compare)
	TestClone()

}
