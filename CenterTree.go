package omap

import (
	"iter"
	"slices"
)

type CenterTree[K any, V any] struct {
	*SliceTree[K, V]
	CenteredSlice []KvSet[K, V]
	Begin         int
	End           int
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
		Begin:         begin,
		End:           begin,
		CenteredSlice: slices,
	}

}

// Remove implemenets [OrderedMap]
func (s CenterTree[K, V]) Remove(key K) (value V, ok bool) {
	value, ok = s.SliceTree.Remove(key)
	s.End = s.Begin + s.Size() - 1
	return
}

// RemoveAll implemenets [OrderedMap]
func (s *CenterTree[K, V]) RemoveAll() (size int) {
	size = s.Size() - 1
	s.Begin = 0
	s.End = 0
	s.Slices = s.Slices[:0]

	return
}

// MassRemove implemenets [OrderedMap]
func (s *CenterTree[K, V]) MassRemove(keys ...K) (total int) {
	total = s.SliceTree.MassRemove(keys...)
	s.End = s.Begin + s.Size()
	return
}

// MassRemoveKV implemenets [OrderedMap]
func (s *CenterTree[K, V]) MassRemoveKV(keys ...K) iter.Seq2[K, V] {
	seq := s.SliceTree.MassRemoveKV(keys...)
	s.End = s.Begin + s.Size() - 1
	return seq
}

func (s *CenterTree[K, V]) reballance(offset, idx int) (pos int) {

	pos = -1
	limit := cap(s.CenteredSlice) - 1
	os := 0
	begin := 0
	if offset == -1 {
		diff := limit - s.End
		if diff < s.Growth {
			return -1
		}
		begin = (diff >> 1) - 1
	} else {
		if s.Begin < s.Growth {
			return -1
		}
		os = 1
		begin = s.Begin >> 1

	}
	ns := make([]KvSet[K, V], cap(s.CenteredSlice))
	copy(ns[begin:begin+idx+os], s.Slices[0:idx+os])
	size := len(s.Slices)
	copy(ns[begin+idx+1+os:], s.Slices[idx:size])
	s.CenteredSlice = ns
	pos = begin + idx + os
	s.Begin = begin
	s.End = begin + size

	return
}

// Put implements [OrderedMap]
func (s *CenterTree[K, V]) Put(k K, v V) {
	size := len(s.Slices)
	limit := cap(s.CenteredSlice) - 1
	if size == 0 {
		s.CenteredSlice[s.Begin] = KvSet[K, V]{k, v}
		s.Slices = s.CenteredSlice[s.Begin : s.Begin+1]
		return
	}

	var idx int
	var offset int
	Cmp := s.Cmp
	Slices := s.Slices
	if size > 10 {
		if Cmp(Slices[0].Key, k) == 1 {
			offset = -1
		} else if idx = size - 1; Cmp(Slices[idx].Key, k) == -1 {
			offset = 1
		} else {
<<<<<<< Updated upstream
			idx, offset = GetIndex(k, Cmp, Slices[1:idx])
			idx++
=======
			mid := getMid(size)
			offset = Cmp(k, Slices[mid].Key)
			if offset == 0 {
				idx = mid
			} else if offset < 0 {
				idx, offset = GetIndex(k, Cmp, Slices[1:mid-1])
				idx++

			} else {
				idx, offset = GetIndex(k, Cmp, Slices[mid+1:size-1])
				idx += mid + 1
			}
>>>>>>> Stashed changes
		}
	} else {
		idx, offset = GetIndex(k, Cmp, Slices)
	}
	var pos int
	switch offset {

	case -1:
		begin := s.Begin - 1
		if begin >= 0 {
			// copy left
			pos = idx + begin
			copy(s.CenteredSlice[begin:idx+s.Begin], Slices[0:idx])
			s.Begin = begin
		} else {
			pos = s.reballance(offset, idx)
			if pos == -1 {
				// grow
				ns := make([]KvSet[K, V], len(s.CenteredSlice)+s.Growth)
				end := s.Growth + s.End
				begin := s.Growth - 1
				copy(ns[begin:idx+begin], Slices[0:idx])

				copy(ns[begin+idx+1:end+1], Slices[idx:])
				s.CenteredSlice = ns
				s.Begin = begin
				pos = s.Begin + idx
				s.End = end
			}
		}
	case 1:
		end := s.End + 1
		if end > limit {
			pos = s.reballance(offset, idx)
			if pos == -1 {
				s.CenteredSlice = slices.Grow(s.CenteredSlice, s.Growth)
				pos = s.Begin + idx + 1
				s.CenteredSlice = s.CenteredSlice[0:cap(s.CenteredSlice)]
				copy(s.CenteredSlice[pos+1:end+1], Slices[idx+1:size])
				s.End = end
				pos = s.Begin + idx + 1
			}

		} else {

			pos = s.Begin + idx + 1
			copy(s.CenteredSlice[s.Begin+idx+1:end], Slices[idx:size])
			s.End = end
		}
	case 0:
		if s.OnOverWrite != nil {
			s.OnOverWrite(k, Slices[idx].Value, v)
		}
		s.Slices[idx].Value = v
		return
	}
	s.CenteredSlice[pos] = KvSet[K, V]{k, v}
	s.Slices = s.CenteredSlice[s.Begin : s.End+1]
}

// SetGrowth implements [OrderedMap]
func (s *CenterTree[K, V]) SetGrowth(grow int) {
	if grow <= MIN_GROWTH {
		s.Growth = MIN_GROWTH
		return
	}
	s.Growth = grow
}

func (s *CenterTree[K, V]) ToTs() OrderedMap[K, V] {
	return &ThreadSafeOrderedMap[K, V]{Tree: s}
}

// Merge implements [OrderedMap]
func (s *CenterTree[K, V]) Merge(set OrderedMap[K, V]) int {
	count := 0
	// do not add to ourself!
	if s == set {
		return 0
	}
	// is this worth trying to optimize?
	for k, v := range set.All() {
		count++
		s.Put(k, v)
	}
	return count
}

// RemoveBetween implements [OrderedMap]
func (s *CenterTree[K, V]) RemoveBetween(a, b K, opt ...int) (total int) {
	s.clearBetween(a, b, func(x, y, t int, res bool) {
		if res {
			total = 1 + y - x
		}
	}, opt...)
	return
}

// RemoveBetween implements [OrderedMap]
func (s *CenterTree[K, V]) RemoveBetweenKV(a, b K, opt ...int) (removed iter.Seq2[K, V]) {
	s.clearBetween(a, b, func(x, y, t int, res bool) {
		if res {
			o := y + 1
			set := slices.Clone(s.Slices[x:o])
			removed = KvIter(set)
		} else {
			removed = KvIter([]KvSet[K, V]{})
		}
	}, opt...)
	return
}

func (s *CenterTree[K, V]) clearBetween(a, b K, cb func(x, y, t int, ok bool), opt ...int) {
	begin, end, total, ok := s.betweenChecks(a, b, opt...)
	cb(begin, end, total, ok)
	if ok {
		if begin == s.Begin && end == s.End {
			s.RemoveAll()
			return
		} else if begin == 0 {
			s.Begin = s.Begin + end + 1
		} else if s.Begin+end == s.End {
			s.End = s.Begin + begin - 1
		} else {
			s.Slices = slices.Delete(s.Slices, begin, 1+end)
			s.End = s.Begin + s.Size() - 1
			return
		}
		s.Slices = s.CenteredSlice[s.Begin : s.End+1]
	}
}
