package omap

import (
	"cmp"
	"testing"
)

func TestFilterBetween(t *testing.T) {
	NewFromSt := func() OrderedMap[int, int] {
		return NewSliceTree[int, int](2, cmp.Compare).ToTs()
	}
	NewFromCt := func() OrderedMap[int, int] {
		return NewCenterTree[int, int](2, cmp.Compare).ToTs()
	}
	var st OrderedMap[int, int]
	Build := func(cb func() OrderedMap[int, int]) {
		st = cb()
		for i := range 6 {
			st.Put(i, i)
		}
	}
	Rm3 := func(cb func() OrderedMap[int, int]) {
		Build(cb)
		st.FilterBetween(func(k, v int) bool {
			return k == 3
		}, 2, 4)
		if _, ok := st.Get(3); st.Size() != 5 || ok {
			t.Logf("%v", st.GetKvSlice())
			t.Fatalf("Expected size of 5, got: %v, expected false, got: %v", st.Size(), ok)
		}
	}
	Rm3(NewFromSt)
	Rm3(NewFromCt)

	RmFromZero := func(cb func() OrderedMap[int, int]) {
		Build(cb)
		st.FilterBetween(func(k, v int) bool {
			return k == 0 || k == 3
		}, 0, 4)
		for _, removed := range []int{0, 3} {

			if _, ok := st.Get(removed); st.Size() != 4 || ok {
				t.Logf("%v", st.GetKvSlice())
				t.Fatalf("Expected size of 4, got: %v, expected false, got: %v", st.Size(), ok)
			}
		}
	}
	RmFromZero(NewFromSt)
	RmFromZero(NewFromCt)
	for _, set := range []OrderedMap[int, int]{NewFromCt(), NewFromSt()} {
		total := 0
		set.FilterBetween(func(k, v int) bool { total++; return true }, -2, -1)
		if total != 0 {
			t.Fatalf("Null lookup should never call anything")
		}
		// odd bit of missed code coverage
		set.RemoveBetweenKV(-2, 3)
	}

}
