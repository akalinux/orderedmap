package omap

import (
	"cmp"
	"iter"
	"testing"
)

func TestLockIters(t *testing.T) {
	s := NewTs[int, int](cmp.Compare)
	for i := range 3 {
		s.Put(i, i)
	}

	for _, f := range []*KvSet[string, func() iter.Seq2[int, int]]{
		{"All", s.All},
		{"Keys", s.Keys},
		{"Values", s.Values},
	} {
		count := 0
		t.Logf("Testing loop iterator from: %s", f.Key)
		for range f.Value() {
			count++
		}
		if count != 3 {
			t.Fatalf("Failed to get 3 elements from %s, got: %d", f.Key, count)
		}
		count = 0
		for range f.Value() {
			count++
			break
		}
		if count != 1 {
			t.Fatalf("Failed to break at element: 0 elements from %s, got: %d", f.Key, count)
		}

	}

}
