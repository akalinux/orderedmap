package omap

import (
	"cmp"
	"fmt"
	"testing"
)

const BENCHMARK_STRING_SIZE = 1024
const BENCHMARK_STRINGS = 100
const BENCHMARK_START = 4
const BENCHMARK_END = 7

func BenchmarkNew(b *testing.B) {

	for tx := BENCHMARK_START; tx <= BENCHMARK_END; tx++ {
		ss := tx * BENCHMARK_STRING_SIZE
		bs := tx * tx * BENCHMARK_STRINGS
		keys := make([]*string, 0, bs)
		f := fmt.Sprintf("%% %dd", ss)

		size := len(fmt.Sprintf("%d", ss))
		seek_f := fmt.Sprintf("%%%dd", size)
		b.Logf("Scan Key format: [%s]", seek_f)
		b.Logf("Key size: %d", size)
		b.Logf("Example using 1,Scan Key format: [%s]", fmt.Sprintf(seek_f, 1))
		b.Logf("Example using %d,Scan Key format: [%s]", bs, fmt.Sprintf(seek_f, bs))

		for i := range bs {
			v := fmt.Sprintf(f, i)
			keys = append(keys, &v)
		}
		var m map[string]*KvSet[*string, any]
		b.Run(
			fmt.Sprintf("Native map Put, keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					m = make(map[string]*KvSet[*string, any], bs)
					for i := range bs {
						m[(*keys[i])[len(*keys[i])-size:]] = &KvSet[*string, any]{keys[i], nil}
					}
				}

			},
		)
		if len(m) != bs {
			b.Fatalf("Go map should contain: %d elements, got: %d", bs, len(m))
		}

		Cmp := func(a, b *string) int {
			return cmp.Compare(
				(*a)[len(*a)-size:],
				(*b)[len(*b)-size:],
			)
		}

		var s *SliceTree[*string, any]
		b.Run(
			fmt.Sprintf("Slicetree, Put, keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {

					s = NewSliceTree[*string, any](bs>>1, Cmp)
					for i := range bs {
						s.Put(keys[i], nil)
					}
				}

			},
		)
		if s.Size() != bs {
			b.Fatalf("Go map should contain: %d elements, got: %d", bs, len(m))
		}
		var ct *CenterTree[*string, any]
		b.Run(
			fmt.Sprintf("CenterTree, Put, keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {

					ct = NewCenterTree[*string, any](bs, Cmp)
					for i := range bs {
						ct.Put(keys[i], nil)
					}
				}

			},
		)
		if ct.Size() != bs {
			b.Fatalf("Go map should contain: %d elements, got: %d", bs, ct.Size())
		}
		b.Run(
			fmt.Sprintf("Native map, Get, keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					for i := range bs {
						_, _ = m[*keys[i]]
					}
				}
			},
		)
		b.Run(
			fmt.Sprintf("SliceTree, Get, keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					for i := range bs {
						s.Get(keys[i])
					}
				}
			},
		)
		b.Run(
			fmt.Sprintf("CenterTree , Get, keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					for i := range bs {
						ct.Get(keys[i])
					}
				}
			},
		)
		b.Run(
			fmt.Sprintf("Native map, Count between, keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					count := 0
					for i := range bs {
						key := fmt.Sprintf(seek_f, i)
						for check := range m {
							if key >= check && key <= check {
								count++
							}
						}
					}
					if count != bs {
						b.Fatalf("Expected, %d, got %d", bs, count)
					}
				}
			},
		)
		b.Run(
			fmt.Sprintf("SliceTree, Count Between nodes keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					count := 0
					for i := range bs {
						count += s.Between(keys[i], keys[i])
					}
					if count != bs {
						b.Fail()
					}

				}
			},
		)
		b.Run(
			fmt.Sprintf("CenterTree, Count Between nodes keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					count := 0
					for i := range bs {
						count += ct.Between(keys[i], keys[i])
					}
					if count != bs {
						b.Fail()
					}

				}
			},
		)
	}
}
