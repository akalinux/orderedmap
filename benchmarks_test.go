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
			fmt.Sprintf("Native map Put keys: [%d]", bs),
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
			fmt.Sprintf("Slicetree Put keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {

					s = NewSliceTree[*string, any](bs, Cmp)
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
			fmt.Sprintf("CenterTree Put keys: [%d]", bs),
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
			fmt.Sprintf("Native map Get keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					for i := range bs {
						_, _ = m[*keys[i]]
					}
				}
			},
		)
		b.Run(
			fmt.Sprintf("SliceTree Get keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					for i := range bs {
						s.Get(keys[i])
					}
				}
			},
		)
		b.Run(
			fmt.Sprintf("CenterTree Get keys: [%d]", bs),
			func(b *testing.B) {
				for range b.N {
					for i := range bs {
						ct.Get(keys[i])
					}
				}
			},
		)
		b.Run(
			fmt.Sprintf("Native map Count between: [%d]", bs),
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
			fmt.Sprintf("SliceTree Count Between: [%d]", bs),
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
			fmt.Sprintf("CenterTree Count Between: [%d]", bs),
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

	for _, size := range []uint64{100, 1000, 10000} {

		var m map[uint64]any
		b.Run(
			fmt.Sprintf("Native Map Uint64 Put keys: %d", size),
			func(b *testing.B) {
				m = make(map[uint64]any, size)
				for k := range size {
					m[uint64(k)] = nil
				}

			},
		)
		var ct OrderedMap[uint64, any]
		b.Run(
			fmt.Sprintf("CenterTree Map Uint64 Put keys: %d", size),
			func(b *testing.B) {
				ct = NewCenterTree[uint64, any](int(size), cmp.Compare)
				for k := range size {
					ct.Put(k, nil)
				}

			},
		)
		var st OrderedMap[uint64, any]
		b.Run(
			fmt.Sprintf("SliceTree Map Uint64 Put keys: %d", size),
			func(b *testing.B) {
				st = NewSliceTree[uint64, any](int(size), cmp.Compare)
				for k := range size {
					st.Put(k, nil)
				}

			},
		)
		b.Run(
			fmt.Sprintf("Native Map Uint64 Get keys: %d", size),
			func(b *testing.B) {
				for k := range size {
					_, _ = m[uint64(k)]
				}

			},
		)
		sets := []OrderedMap[uint64, any]{ct, st}
		names := []string{"CenterTree", "SliceTree"}
		for idx := range 2 {
			s := sets[idx]
			b.Run(
				fmt.Sprintf("%s Map Uint64 int Get keys: %d", names[idx], size),
				func(b *testing.B) {
					for k := range size {
						//k <<= 16
						s.Get(k)
					}

				},
			)
		}
	}
	for _, size := range []uint8{10, 100, 255} {

		m := make(map[uint8]any, size)
		b.Run(
			fmt.Sprintf("Native Map Uint8 Put keys: %d", size),
			func(b *testing.B) {
				for k := range size {
					m[uint8(k)] = nil
				}

			},
		)
		ct := NewCenterTree[uint8, any](int(size), cmp.Compare)
		b.Run(
			fmt.Sprintf("CenterTree Map Uint8 int Put keys: %d", size),
			func(b *testing.B) {
				for k := range size {
					ct.Put(k, nil)
				}

			},
		)
		st := NewSliceTree[uint8, any](int(size), cmp.Compare)
		b.Run(
			fmt.Sprintf("SliceTree Map Uint8 int Put keys: %d", size),
			func(b *testing.B) {
				for k := range size {
					//k <<= 16
					st.Put(k, nil)
				}

			},
		)
		b.Run(
			fmt.Sprintf("Native Map Uint8 Get keys: %d", size),
			func(b *testing.B) {
				for k := range size {
					_, _ = m[uint8(k)]
				}

			},
		)
		sets := []OrderedMap[uint8, any]{ct, st}
		names := []string{"CenterTree", "SliceTree"}
		for idx := range 2 {
			s := sets[idx]
			b.Run(
				fmt.Sprintf("%s Map Uint8 int Get keys: %d", names[idx], size),
				func(b *testing.B) {
					for k := range size {
						//k <<= 16
						s.Get(k)
					}

				},
			)
		}
	}
}
