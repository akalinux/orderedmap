package omap

import (
	"cmp"
	"fmt"
	"testing"
)

const FORMAT = "%05d"

func BenchmarkNew(b *testing.B) {
	for _, size := range []int{10000, 5000, 100} {
		var m map[string]any
		var s OrderedMap[string, any]
		BuildM := func() {
			m = make(map[string]any, size)
			for i := range size {
				k := fmt.Sprintf(FORMAT, i)
				m[k] = nil
			}
		}
		BuildCt := func() {
			s = NewCenterTree[string, any](size, cmp.Compare)
			for i := range size {
				k := fmt.Sprintf(FORMAT, i)
				s.Put(k, nil)
			}
		}

		b.Run(fmt.Sprintf("Build Native: %d", size), func(b *testing.B) {
			for range b.N {
				BuildM()
			}
		})
		b.Run(fmt.Sprintf("Build CenterTree: %d", size), func(b *testing.B) {
			for range b.N {
				BuildCt()
			}
		})
		b.Run(fmt.Sprintf("Get Native: %d", size), func(b *testing.B) {
			for range b.N {
				for i := range size {
					k := fmt.Sprintf(FORMAT, i)
					_, _ = m[k]
				}
			}
		})
		b.Run(fmt.Sprintf("Get CenterTree: %d", size), func(b *testing.B) {
			for range b.N {
				for i := range size {
					k := fmt.Sprintf(FORMAT, i)
					s.Get(k)
				}
			}
		})

		set := []KvSet[int, any]{}
		for i := range size {
			set = append(set, KvSet[int, any]{i, i})
		}
		b.Run("Test GetIndex", func(b *testing.B) {
			for range b.N {
				for i := range size {
					GetIndex(i, cmp.Compare, set)
				}
			}
		})
		/*
			b.Run("Test LessIndex", func(b *testing.B) {
				for range b.N {
					for i := range size {
						LessIndex(i, cmp.Less, set)
					}
				}
			})
		*/
	}
}
