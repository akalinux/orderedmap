package omap

import (
	"cmp"
	"testing"
)

func TestCenterTreePut(t *testing.T) {
	nt := NewCenterTree[int, int](2, cmp.Compare)
	overwrite := 0
	nt.OnOverWrite = func(key, oldValue, newValue int) {
		overwrite++
		t.Logf("Overwritting Key: %d, Old: %d, New: %d", key, oldValue, newValue)
	}

	check := []int{}
	Sane := func(k int) {
		t.Logf("Put test of nt.Put(%d,%d)", k, k)
		idx, offset := GetIndex(k, cmp.Compare, nt.Slices)
		t.Logf("Indx is: %d, offset is: %d", idx, offset)
		nt.Put(k, k)

		if v, ok := nt.Get(k); !ok || v != k {
			t.Log("*** Error: Dumping out state of the array")
			for key, value := range nt.All() {
				t.Logf("key: %d, value %d", key, value)
			}
			t.Fatalf("Expected true for ok, got: %v, expected %d for value got: %v", ok, k, v)
			return
		}
		check = append(check, k)
		if size := len(check); size != nt.Size() {
			t.Fatalf("Expected size: %d, got: %d", size, nt.Size())
			return
		}
		for _, id := range check {
			if _, ok := nt.Get(id); !ok {
				t.Log("*** Error: Dumping out state of the array")
				for key, value := range nt.All() {
					t.Logf("key: %d, value %d", key, value)
				}
				t.Fatalf("nt.Get(%d) should return true for key: %d", id, id)
				return
			}
		}
		t.Logf("Begin: %d, end: %d, cap: %d, Growth: %d", nt.Begin, nt.End, cap(nt.CenteredSlice), nt.Growth)
	}
	t.Logf("Initalization test")
	Sane(3)
	t.Logf("Ovewrwrite test")
	nt.Put(3, -3)
	if v, ok := nt.Get(3); !ok || v != -3 || overwrite != 1 {
		t.Fatalf("Expected true for ok, got: %v, expected -3 for value got: %v, overwrite should be 1, got %d", ok, v, overwrite)
	}

	t.Logf("Preprend tests")
	Sane(2)
	Sane(1)
	t.Logf("Append tests")
	Sane(6)
	Sane(8)

	t.Logf("between tests")
	Sane(4)

	Sane(7)
	for _, k := range []int{5, 9} {
		idx, offset := GetIndex(k, cmp.Compare, nt.Slices)
		t.Logf("Key: %d, idx: %d, offset: %d", k, idx, offset)

	}
	Sane(9)
	Sane(10)
	Sane(11)
	Sane(0)
	Sane(-1)
	Sane(-2)
	Sane(12)
	Sane(13)
	Sane(14)
	Sane(15)
	Sane(16)
	Sane(17)
	Sane(20)
}

func TestCenterTreeRemove(t *testing.T) {
	s := NewCenterTree[int, int](2, cmp.Compare)
	s.Put(1, 1)
	s.Put(2, 2)
	if value, ok := s.Remove(2); value != 2 || !ok || s.Size() != 1 {
		t.Fatalf("expected: value: 2, got: %d, ok true, got: %v, Size: 1, got: %d", value, ok, s.Size())
	}
	if value, ok := s.Remove(1); value != 1 || !ok || s.Size() != 0 {
		t.Fatalf("expected: value: 1, got: %d, ok true, got: %v, Size: 0, got: %d", value, ok, s.Size())
	}
	s.Put(1, 1)
	s.Put(2, 2)
	if total := s.RemoveAll(); total != 2 || s.Size() != 0 {
		t.Fatalf("Expected the set to be empty, tota: %d, size: %d", total, s.Size())
	}
	s.Put(1, 1)
	s.Put(2, 2)
	if total := s.MassRemove(1, 2); total != 2 || s.Size() != 0 {
		t.Fatalf("Expected the set to be empty, tota: %d, size: %d", total, s.Size())
	}
	s.Put(1, 1)
	s.Put(2, 2)
	total := 0
	for k, v := range s.MassRemoveKV(1, 2) {
		total += k + v
	}
	if total != 6 || s.Size() != 0 {
		t.Fatalf("Expected the set to be empty, tota: %d, size: %d", total, s.Size())
	}
}

func TestCenterTreeRemoveBetween(t *testing.T) {
	var s *CenterTree[int, int]

	Reset := func() *CenterTree[int, int] {
		s = NewCenterTree[int, int](2, cmp.Compare)
		s.Put(1, 1)
		s.Put(2, 2)
		s.Put(3, 3)
		return s
	}
	s = Reset()
	if total := s.RemoveBetween(-1, 2); total != 2 && s.Size() == 1 {
		t.Fatalf("Expected total of: 2, got: %d, expected size of: 1, got: %d ", total, s.Size())
	}
	Check := func(test string, expected, size, find, start, end int, opts ...int) {
		t.Log(test)
		s = Reset()
		t.Logf("# -- Starting Begin: %d, End: %d -- $", s.Begin, s.End)
		total := 0
		for k, v := range s.RemoveBetweenKV(start, end, opts...) {
			total += k + v
		}
		t.Logf("Exepcting all keys + values to provide a total value of %d", expected)
		if total != expected && s.Size() == size {
			t.Fatalf("Expected total of: %d, got: %d, expected size of: %d, got: %d ", expected, total, size, s.Size())
		}
		if _, ok := s.Get(find); !ok {
			t.Fatalf("Expected key: %d to exist, but it does not", find)
		}
		for _, k := range []int{start, end} {
			if _, ok := s.Get(k); ok {
				t.Fatalf("%d should no longer exist in our data set", k)
			}
		}
		t.Logf("# -- Ending Begin: %d, End: %d -- #", s.Begin, s.End)

	}
	Check("remove between -1,2", 6, 1, 3, -1, 2)

	Check("remove between FIRST_KEY and 2", 6, 1, 3, 100, 2, FIRST_KEY)
	Check("remove between, 2, 100", 10, 1, 1, 2, 100)
	Check("Remove between 2 and LAST_KEY", 10, 1, 1, 2, -1, LAST_KEY)

	Check("Remove mid: 2,2", 4, 2, 1, 2, 2)

}
