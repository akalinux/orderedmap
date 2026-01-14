package omap

import (
	"cmp"
	"testing"
)

const BENCHMARK_KEY_COUNT = 8192

func BenchmarkPutBENCHMARK_KEY_COUNT(b *testing.B) {
	for range b.N {
		m := NewSliceTree[uint64, uint64](BENCHMARK_KEY_COUNT, cmp.Compare)
		//m := New[uint64, uint64](cmp.Compare)
		for i := range BENCHMARK_KEY_COUNT {
			v := uint64(i)
			m.Put(v, v)
		}
		// overwrite
		for i := range BENCHMARK_KEY_COUNT {
			v := uint64(i)
			m.Put(v, v)
		}
		if len(m.Slices) != BENCHMARK_KEY_COUNT {
			b.Fail()
		}
	}
}

func BenchmarkMapPutBENCHMARK_KEY_COUNT(b *testing.B) {
	for range b.N {
		m := make(map[uint64]uint64, BENCHMARK_KEY_COUNT)
		for i := range BENCHMARK_KEY_COUNT {
			v := uint64(i)
			m[v] = v
		}
		// overwrite
		for i := range BENCHMARK_KEY_COUNT {
			v := uint64(i)
			m[v] = v
		}
		if len(m) != BENCHMARK_KEY_COUNT {
			b.Fail()
		}
	}
}

func BenchmarkGetBENCHMARK_KEY_COUNT(b *testing.B) {
	m := New[uint64, uint64](cmp.Compare)
	for i := range BENCHMARK_KEY_COUNT {
		v := uint64(i)
		m.Put(v, v)
	}
	for range b.N {
		//m := NewSliceTree[uint64, uint64](BENCHMARK_KEY_COUNT, cmp.Compare)
		for i := range BENCHMARK_KEY_COUNT {
			c := uint64(i)
			if v, ok := m.Get(c); !ok || v != uint64(c) {
				b.Fail()
			}
		}
	}
}

func BenchmarkGetMapBENCHMARK_KEY_COUNT(b *testing.B) {
	m := make(map[uint64]uint64)
	for i := range BENCHMARK_KEY_COUNT {
		v := uint64(i)
		m[v] = v
	}
	for range b.N {
		//m := NewSliceTree[uint64, uint64](BENCHMARK_KEY_COUNT, cmp.Compare)
		for i := range BENCHMARK_KEY_COUNT {
			c := uint64(i)
			if v, ok := m[c]; !ok || v != uint64(c) {
				b.Fail()
			}
		}
	}
}
