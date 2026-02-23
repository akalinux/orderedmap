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
		nt.Put(k, k)

		if v, ok := nt.Get(k); !ok || v != k {
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
	for _, k := range check {
		idx, offset := GetIndex(k, cmp.Compare, nt.Slices)
		if _, ok := nt.Get(k); !ok {
			t.Logf("Failed to fetchl %d", k)
			return
		}
		t.Logf("Key: %d, idx: %d, offset: %d", k, idx, offset)
	}

	t.Logf("between tests")
	Sane(7)

}
