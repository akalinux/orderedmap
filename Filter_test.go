package omap

import (
	"cmp"
	"testing"
)

func TestFilter(t *testing.T) {

	st := NewSliceTree[int, int](5, cmp.Compare).ToTs()
	ct := NewCenterTree[int, int](5, cmp.Compare).ToTs()
	names := []string{"SliceTree", "CenterTree"}
	for test, s := range []OrderedMap[int, int]{st, ct} {
		name := names[test]
		t.Logf("Starting remove test of: %s", name)
		for k := range 5 {
			s.Put(k, k)
		}

		s.Filter(func(k, v int) bool {
			remove := k == 1 || k == 3
			if remove {

				t.Logf("Removing Key: %d", k)
			}
			return remove
		})
		if s.Size() != 3 {
			t.Log("Test failed, dumping array contents")
			for k, v := range s.All() {
				t.Logf("Have key: %d, value: %d", k, v)
			}
			t.Fatalf("Expected a size of: 3, got: %d", s.Size())
		}
		s.Filter(func(k, v int) bool { return true })
		if s.Size() != 0 {
			t.Fatalf("Expected an empty set!")
		}
	}
}
