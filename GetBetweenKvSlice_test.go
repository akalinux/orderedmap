package omap

import (
	"cmp"
	"testing"
)

func TestGetBetweenKvSlice(t *testing.T) {
	s := NewTs[int, int](cmp.Compare)
	for k := range 5 {
		s.Put(k, k)
	}

	Check := func(name string, a, b, size, sum int) {
		t.Log(name)
		list := s.GetBetweenKvSlice(a, b)
		total := 0
		for _, kv := range list {
			t.Logf("GOt KvSet: %v", kv)
			total += kv.Key + kv.Value
		}
		if len(list) != size || total != sum {
			t.Fatalf("Expected a size of: %d, got: %d, expected sum: %d, got: %d", size, len(list), sum, total)
		}
	}
	Check("Begin out of bounds", -1, -1, 0, 0)
	Check("End out of bounds", 5, 5, 0, 0)
	Check("Before Start -1 to 1 ", -1, 1, 2, 2)
	Check("Afer End 3 to 5 ", 3, 4, 2, 14)
	Check("Inside End 2 to 3 ", 2, 3, 2, 10)
}
