package omap

import (
	"cmp"
	"fmt"
	"testing"
)

const BENCHMARK_STRING_SIZE = 1024
const BENCHMARK_STRINGS = 100
const BENCHMARK_PASS = 4

func BenchmarkNew(b *testing.B) {

	for tx := 1; tx <= BENCHMARK_PASS; tx++ {
		ss := tx * BENCHMARK_STRING_SIZE
		bs := tx * tx * BENCHMARK_STRINGS
		keys := make([]*string, 0, bs)
		f := fmt.Sprintf("%% %dd", ss)
		size := len(fmt.Sprintf("%d", ss))
		for i := range bs {
			v := fmt.Sprintf(f, i)
			keys = append(keys, &v)
		}
		var m map[string]KvSet[*string, any]
		b.Run(
			fmt.Sprintf("Native map, size: [%d], keys: [%d]", size, bs),
			func(b *testing.B) {
				for range b.N {
					m = make(map[string]KvSet[*string, any], bs)
					for i := range bs {
						m[(*keys[i])[len(*keys[i])-size:]] = KvSet[*string, any]{keys[i], nil}
					}
				}

			},
		)

		Cmp := func(a, b *string) int {
			return cmp.Compare(
				(*a)[len(*a)-size:],
				(*b)[len(*b)-size:],
			)
		}

		var s *SliceTree[*string, any]
		b.Run(
			fmt.Sprintf("Slicetree, size [%d], count [%d]", size, bs),
			func(b *testing.B) {
				for range b.N {

					s = NewSliceTree[*string, any](bs, Cmp)
					for i := range bs {
						s.Put(keys[i], nil)
					}
				}

			},
		)
		b.Run(
			fmt.Sprintf("Native map, Get size: [%d], keys: [%d]", size, bs),
			func(b *testing.B) {
				for range b.N {
					for i := range bs {
						_, _ = m[*keys[i]]
					}
				}
			},
		)
		b.Run(
			fmt.Sprintf("SliceTree, Get size: [%d], keys: [%d]", size, bs),
			func(b *testing.B) {
				for range b.N {
					for i := range bs {
						s.Get(keys[i])
					}
				}
			},
		)
	}
}
