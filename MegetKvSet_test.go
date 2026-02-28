package omap

import (
	"cmp"
	"testing"
)

func TestMergeKvSet(t *testing.T) {
	res, end := MergeKvSet(
		[]KvSet[int, int]{{0, 1}},
		[]KvSet[int, int]{{1, 2}},
		make([]KvSet[int, int], 2),
		0,
		0,
		cmp.Compare,
		nil,
	)

	CheckKvSet := func(name string, expectedSize, expectedSum, expectedEnd int, res []KvSet[int, int], end int) {
		t.Logf("Testing: [%s]", name)
		sum := 0
		size := len(res)
		for _, set := range res {
			sum += set.Key + set.Value
		}
		if sum != expectedSum || end != expectedEnd || size != expectedSize {

			t.Logf("Error, outputing Resulting array: %v", res)
			t.Fatalf("Expected a sum of: %d, got %d, expected size of: %d, got %d, expected end: %d, got: %d",
				expectedSum,
				sum,
				expectedSize,
				size,
				expectedEnd,
				end,
			)
		}
	}
	CheckKvSet("basic a<b", 2, 4, 1, res, end)

	res, end = MergeKvSet(
		[]KvSet[int, int]{{1, 2}},
		[]KvSet[int, int]{{0, 1}},
		make([]KvSet[int, int], 2),
		0,
		0,
		cmp.Compare,
		nil,
	)
	CheckKvSet("basic a>b", 2, 4, 1, res, end)

	res, end = MergeKvSet(
		[]KvSet[int, int]{{1, 1}, {3, 3}},
		[]KvSet[int, int]{{2, 2}, {4, 4}},
		make([]KvSet[int, int], 4),
		0,
		1,
		cmp.Compare,
		nil,
	)
	CheckKvSet("Merge 2 non overlapping sets", 4, 20, 3, res, end)

	res, end = MergeKvSet(
		[]KvSet[int, int]{{}, {1, 1}, {3, 3}},
		[]KvSet[int, int]{{2, 2}, {4, 4}},
		make([]KvSet[int, int], 4),
		1,
		2,
		cmp.Compare,
		nil,
	)
	CheckKvSet("offset of 1 Merge 2 non overlapping sets", 5, 20, 4, res, end)
	res, end = MergeKvSet(
		[]KvSet[int, int]{{}, {1, 1}, {3, 3}},
		[]KvSet[int, int]{{2, 2}, {3, 3}, {4, 4}},
		make([]KvSet[int, int], 4),
		1,
		2,
		cmp.Compare,
		func(i1, i2, i3 int) {},
	)
	CheckKvSet("offset of 1 Merge 2 non overlapping sets, 1 overlapping set", 5, 20, 4, res, end)
	Build := func(slices []KvSet[int, int]) *SliceTree[int, int] {
		return &SliceTree[int, int]{
			Cmp:    cmp.Compare[int],
			Slices: slices,
			Growth: 2,
		}
	}
	var check OrderedMap[int, int] = Build([]KvSet[int, int]{{1, 1}, {3, 3}})
	total := check.FastMerge(Build([]KvSet[int, int]{{2, 2}, {4, 4}}))
	if total != 2 {
		t.Fatalf("Expected 2 keys to be added, got %d", total)
	}
	total = check.Merge(Build([]KvSet[int, int]{{2, 2}, {4, 4}}))
	if total != 0 {
		t.Fatalf("Should not have merged anyting!")
	}
	total = check.FastMerge(Build([]KvSet[int, int]{}))
	if total != 0 {
		t.Fatalf("Should not have merged anyting!")
	}

	check = NewCenterTree[int, int](2, cmp.Compare)
	check.Put(1, 1)
	check.Put(3, 3)

	total = check.FastMerge(Build([]KvSet[int, int]{{2, 2}, {4, 4}}))
	if total != 2 {
		t.Fatalf("Expected 2 keys to be added, got %d", total)
	}
	check = NewCenterTree[int, int](2, cmp.Compare).ToTs()
	check.Put(1, 1)
	check.Put(3, 3)
	src := Build([]KvSet[int, int]{})
	check.GetKvSlice()
	total = check.FastMerge(src)
	if total != 0 || check.Size() != 2 {
		t.Fatalf("Expected 0 keys to be added, got %d, Size should be: 2, got: %d", total, check.Size())
	}
	chk := NewCenterTree[int, int](2, cmp.Compare)
	total = chk.FastMerge(Build([]KvSet[int, int]{{2, 2}, {4, 4}}))
	if total != 2 || chk.Size() != 2 {
		t.Fatalf("Expected 2 keys to be added, got %d, Size should be: 2, got: %d", total, check.Size())
	}
	check = Build([]KvSet[int, int]{})
	total = check.FastMerge(Build([]KvSet[int, int]{{2, 2}, {4, 4}}))
	if total != 2 || chk.Size() != 2 {
		t.Fatalf("Expected 2 keys to be added, got %d, Size should be: 2, got: %d", total, check.Size())
	}

}
