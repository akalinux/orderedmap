package omap

import (
	"cmp"
	"testing"
)

const BENCHMARK_KEY_COUNT = 18

func BenchmarkNew(b *testing.B) {
	for range b.N {
		NewSliceTree[int, int](0, cmp.Compare)
	}
}
func BenchmarkNewFromMap(b *testing.B) {
	m := make(map[int]int, BENCHMARK_KEY_COUNT)
	for i := range BENCHMARK_KEY_COUNT {
		m[i] = i
	}
	for range b.N {
		NewFromMap(m, cmp.Compare)
	}
}
func BenchmarkToMap(b *testing.B) {
	m := FreshOm()
	for range b.N {
		ToMap(m)
	}
}

func BenchmarkPut(b *testing.B) {
	om := FreshOm()
	for range b.N {
		om.Put(8, 10)
	}
}

func FreshOm() OrderedMap[int, int] {

	om := NewSliceTree[int, int](BENCHMARK_KEY_COUNT, cmp.Compare)
	for i := range BENCHMARK_KEY_COUNT {
		v := int(i)
		om.Put(v, v)
	}
	return om
}
func BenchmarkGet(b *testing.B) {
	om := FreshOm()
	for range b.N {
		om.Get(8)
	}
}

func BenchmarkKeys(b *testing.B) {
	om := FreshOm()
	for range b.N {
		om.Keys()
	}
}

func BenchmarkValues(b *testing.B) {
	om := FreshOm()
	for range b.N {
		om.Values()
	}
}
func BenchmarkAll(b *testing.B) {
	om := FreshOm()
	for range b.N {
		om.All()
	}
}

func BenchmarkSize(b *testing.B) {
	om := FreshOm()
	for range b.N {
		om.Size()
	}
}

func BenchmarkMerge(b *testing.B) {
	om := FreshOm()
	om2 := FreshOm()
	for range b.N {
		om.Merge(om2)
	}
}

func BenchmarkBetween(b *testing.B) {
	om := FreshOm()
	for range b.N {
		om.Between(8, 9)
	}
}

func BenchmarkBetweenKV(b *testing.B) {
	om := FreshOm()
	for range b.N {
		om.BetweenKV(8, 9)
	}
}

func BenchmarkMergeMap(b *testing.B) {
	om := FreshOm()
	m := make(map[int]int, BENCHMARK_KEY_COUNT)
	for i := range BENCHMARK_KEY_COUNT {
		m[i] = i
	}
	for range b.N {
		Merge(om, m)
	}
}

func BenchmarkBetweenFirst(b *testing.B) {
	o := FreshOm()
	for range b.N {
		o.Between(-1, 8, FIRST_KEY)
	}
}

func BenchmarkBetweenLast(b *testing.B) {
	o := FreshOm()
	for range b.N {
		o.Between(8, -1, LAST_KEY)
	}
}
func BenchmarkBetweenFirstKV(b *testing.B) {
	o := FreshOm()
	for range b.N {
		o.BetweenKV(-1, 8, FIRST_KEY)
	}
}

func BenchmarkBetweenLastKV(b *testing.B) {
	o := FreshOm()
	for range b.N {
		o.BetweenKV(8, -1, LAST_KEY)
	}
}

func BenchmarkMassRemove(b *testing.B) {
	o := FreshOm()
	for range b.N {
		o.MassRemove(1, 3, 5)
	}
}

func BenchmarkRemoveAll(b *testing.B) {
	o := FreshOm()
	for range b.N {
		o.RemoveAll()
	}
}

func BenchmarkMassRemoveKV(b *testing.B) {
	o := FreshOm()
	for range b.N {
		o.MassRemoveKV(1, 3, 5)
	}
}

func BenchmarkRemoveBetween(b *testing.B) {
	o := FreshOm()
	for range b.N {
		o.RemoveBetween(1, 3)
	}
}

func BenchmarkRemoveBetweenKV(b *testing.B) {
	o := FreshOm()
	for range b.N {
		o.RemoveBetweenKV(1, 3)
	}
}
