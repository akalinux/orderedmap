package omap

import (
	"cmp"
	"iter"
	"testing"
)

func seqToSeq2(s iter.Seq[int]) func() iter.Seq2[int, int] {
	return func() iter.Seq2[int, int] {
		return func(yield func(int, int) bool) {
			count := -1
			for k := range s {
				count++
				if !yield(count, k) {
					return
				}
			}

		}
	}
}

func TestLockIters(t *testing.T) {
	s := New[int, int](cmp.Compare)
	for i := range 3 {
		s.Put(i, i)
	}

	for range s.keys() {
		break
	}

	for _, f := range []*KvSet[string, func() iter.Seq2[int, int]]{
		{"All", s.ToTs().All},
		{"Keys", seqToSeq2(s.ToTs().Keys())},
		{"Values", seqToSeq2(s.ToTs().Values())},
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

func TestMerge(t *testing.T) {
	dst := NewTs[int, int](cmp.Compare)
	src := NewTs[int, int](cmp.Compare)
	for i := range 3 {
		src.Put(i, i+3)
	}

	if check := dst.Merge(src); check != src.Size() {
		t.Fatalf("Expected size to be: %d, got %d", src.Size(), check)
	}

	if check := src.Merge(src); check != 0 {
		t.Fatalf("Should not be able to merge onto an ourself")
	}

	if src.Size() != dst.Size() {
		t.Fatalf("Source and dst do not match")
	}
	next, stop := iter.Pull2(dst.All())
	defer stop()
	for k, v := range dst.All() {
		x, y, ok := next()
		if x != k || y != v || !ok {
			t.Fail()
		}
	}
	stop()
	dst = NewTs[int, int](cmp.Compare)

	m := map[int]int{1: 2}
	if count := Merge(dst, m); count != 1 {
		t.Fail()
	}
	if v, ok := dst.Get(1); !ok || v != 2 {
		t.Fail()
	}

	dst = New[int, int](cmp.Compare)
	dst.Merge(dst)

}

func TestBug(t *testing.T) {
	s := NewSliceTree[int64, any](100, cmp.Compare)
	s.Put(1769664118175, nil)
	for range 1 {

		iter := s.RemoveBetweenKV(-1, 1769664118175, FIRST_KEY)
		for range iter {

		}

	}
}
