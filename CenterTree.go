package omap

import (
	"fmt"
)

type CenterTree[K any, V any] struct {
	*SliceTree[K, V]
	slices []KvSet[K, V]
	begin  int
	end    int
}

const MIN_GROWTH = 1

func NewCenterTree[K any, V any](growth int, cmp func(a, b K) int) *CenterTree[K, V] {
	if growth < MIN_GROWTH {
		growth = MIN_GROWTH
	}

	slices := make([]KvSet[K, V], growth*2+1)
	mid := getMid(cap(slices))
	begin := mid

	return &CenterTree[K, V]{
		SliceTree: &SliceTree[K, V]{
			Cmp:    cmp,
			Growth: growth,
			// pass an empty slice
			Slices: slices[:0],
		},
		begin:  begin,
		end:    begin,
		slices: slices,
	}

}

func (s *CenterTree[K, V]) Put(k K, v V) {
	size := len(s.Slices)
	limit := cap(s.slices) - 1
	if size == 0 {
		s.slices[s.begin] = KvSet[K, V]{k, v}
		s.Slices = s.slices[s.begin : s.begin+1]
		return
	}

	idx, offset := GetIndex(k, s.Cmp, s.Slices)
	var pos int
	switch offset {
	case 0:
		if s.OnOverWrite != nil {
			s.OnOverWrite(k, s.Slices[idx].Value, v)
		}
		s.Slices[idx].Value = v
		return
	case -1:
		begin := s.begin - 1
		if begin >= 0 {
			// copy left
			fmt.Printf("Copy left block\n")
			pos = idx + begin
			copy(s.slices[begin:idx+s.begin], s.Slices[0:idx])
			s.begin = begin
		} else {
			diff := limit - s.end
			fmt.Printf("Diff: %d\n", diff)
			if diff >= s.Growth {
				// reballance
				o := diff >> 1
				begin := o - 1
				end := s.end - o
				fmt.Printf("reballance, begin: %d, end %d\n", begin, end)

				copy(s.slices[begin:idx+begin], s.Slices[0:idx])
				s.begin = begin
				s.end = end
				pos = begin + idx
			} else {
				// grow
				fmt.Printf("Grow left\n")
				ns := make([]KvSet[K, V], len(s.slices)+s.Growth)
				end := s.Growth + s.end
				begin := s.Growth - 1
				fmt.Printf("New begin: %d, New End: %d, new size: %d\n", begin, end, len(ns))
				copy(ns[begin:idx+begin], s.Slices[0:idx])

				copy(ns[begin+idx+1:], s.Slices[idx:])
				s.slices = ns
				s.begin = begin
				pos = s.begin + idx
				s.end = end
			}
		}
	default:
		end := s.end + 1
		if end > limit {
			fmt.Printf("Todo\n")
		} else {
			if idx+s.begin == size-1 {
				// simple append
				s.end = end
				pos = s.end
			} else {

				pos := s.begin + idx
				fmt.Printf("Limit: %d,start: %d,end: %d\n", limit, pos+1, end)
				copy(s.slices[pos+1:end], s.Slices[idx:size])
				s.end = end
			}
		}
	}
	s.slices[pos] = KvSet[K, V]{k, v}
	s.Slices = s.slices[s.begin : s.end+1]
}
