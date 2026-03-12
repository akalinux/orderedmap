package omap

import (
	"cmp"
	"iter"
	"slices"
)

func KvIter[K any, V any](set []KvSet[K, V]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, row := range set {
			if !yield(row.Key, row.Value) {
				return
			}
		}
	}
}

func (s *SliceTree[K, V]) rangedel(a, b int) {
	s.Slices = slices.Delete(s.Slices, a, b+1)
}

func reverse(a, b int) int {
	return cmp.Compare(b, a)
}

// Returns the index and offset of a given key.
//
// The index is the current relative position in the slice.
//
// The offset represents where the item would be placed:
//   - offset of 0, at index value.
//   - offset of 1, expand the slice after the inddex and put the value to the right of the index
//   - offset of -1, expand the slice before the index and put the value to left of the current position
//
// Complexity: o(log n)
func GetIndex[K any, V any](k K, Cmp func(a, b K) int, Slices []KvSet[K, V]) (index, offset int) {
	end := len(Slices)
	switch end {
	case 0:
		return 0, 0
	case 1:
		return 0, Cmp(k, Slices[0].Key)
	}
	mid := getMid(end)
	offset = Cmp(k, Slices[mid].Key)
	end--
	begin, diff := 0, 0
	for {
		switch offset {
		case -1:
			// moving left
			end = mid - 1
			diff = end - begin
		case 1:
			begin = mid + 1
			diff = end - begin
		default:
			return mid, 0
		}
		switch diff {
		case -1:
			return begin, -1
		case 0:
			return begin, Cmp(k, Slices[begin].Key)
		}
		diff += 1
		mid = begin + getMid(diff)

		offset = Cmp(k, Slices[mid].Key)
	}
}

func getMid(size int) int {
	// shift right 1 same as divide  by 2, bitwise is always faster
	// size&1 is the same as size%2, but biwise is faster
	return (size-2)>>1 + size&1
}

// Merges a map into an OrderedMap instance.
func Merge[K comparable, V any](dst OrderedMap[K, V], src map[K]V) int {
	count := 0
	// is this worth trying to optimize?
	for k, v := range src {
		count++
		dst.Put(k, v)
	}
	return count
}

func MergeKvSet[K any, V any](left, right, dst []KvSet[K, V], start, finish int, Cmp func(a, b K) int, OnOverwrite func(K, V, V)) (res []KvSet[K, V], end int) {

	res = dst
	pos := 0
	f := len(right)
	dp := start
	i := start
	cmp := 0
	for i <= finish && pos < f {

		cmp = Cmp(left[i].Key, right[pos].Key)
		switch cmp {
		case -1:
			dst[dp] = left[i]
			i++
		case 1:
			dst[dp] = right[pos]
			pos++
		default:
			if OnOverwrite != nil {
				OnOverwrite(left[i].Key, left[i].Value, right[pos].Value)
			}
			dst[dp] = right[pos]
			i++
			pos++
		}
		end = dp
		dp++
	}
	if i <= finish {
		res = append(res[0:end+1], left[i:finish+1]...)
		end += finish - i + 1
	} else if pos < f {
		res = append(res[0:end+1], right[pos:f]...)
		end += f - pos
	}
	return
}
