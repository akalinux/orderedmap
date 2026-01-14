package omap

import (
	"cmp"
	"testing"
)

const BENCHMARK_KEY_COUNT = 8192

var om OrderedMap[uint64, uint64]

func BenchmarkPutBENCHMARK_KEY_COUNT(b *testing.B) {
	for range b.N {
		om = NewSliceTree[uint64, uint64](BENCHMARK_KEY_COUNT, cmp.Compare)
		for i := range BENCHMARK_KEY_COUNT {
			v := uint64(i)
			om.Put(v, v)
		}
		// overwrite
		for i := range BENCHMARK_KEY_COUNT {
			v := uint64(i)
			om.Put(v, v)
		}
		if om.Size() != BENCHMARK_KEY_COUNT {
			b.Fail()
		}
	}
}

func BenchmarkGetBENCHMARK_KEY_COUNT(b *testing.B) {

	for range b.N {
		//m := NewSliceTree[uint64, uint64](BENCHMARK_KEY_COUNT, cmp.Compare)
		for i := range BENCHMARK_KEY_COUNT {
			c := uint64(i)
			if v, ok := om.Get(c); !ok || v != uint64(c) {
				b.Fail()
			}
		}
	}
}
