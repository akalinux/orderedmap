package omap

import (
	"cmp"
	"iter"
	"slices"
)

type SliceTree[K any, V any] struct {

	// Internally managed keys slice
	Slices []*KvSet[K, V]

	// Compare function.
	Cmp func(a, b K) int

	// Required non 0 value, determins by what capacity we grow the internal
	// Slice.  Default is 1.
	Growth int

	// Required non nil value, called when ever a value is overwritten.
	// Seting this funtion saves on having to write a check when data is overwritten.
	OnOverWrite func(key K, oldValue V, newValue V) V
}

// Stub overwrite method, used by the constructor.  This is the default callback used when a value is overwritten.
func stubOnOverwrite[K any, V any](key K, oldValue, newValue V) V { return newValue }

// Creatss a new SliceTree with the internal Slice set to "size".
func NewSliceTree[K any, V any](size int, cb func(a, b K) int) *SliceTree[K, V] {
	return &SliceTree[K, V]{
		Slices:      make([]*KvSet[K, V], 0, size),
		Cmp:         cb,
		Growth:      100,
		OnOverWrite: stubOnOverwrite[K, V],
	}
}

// Creates a new SliceTee with the default Slices size of 100.  If you require more control over the starting size of the slice
// use the NewSliceTree function in stead.
func New[K any, V any](cb func(a, b K) int) *SliceTree[K, V] {
	return NewSliceTree[K, V](100, cb)
}

func (s *SliceTree[K, V]) getMid(size int) int {
	// shift right 1 same sa divide  by 2.. gotta love int maths
	return (size-2)>>1 + size&1
}

// Tries to remove the element of k, returns false if it fails.
// Complexity: o(log n)
func (s *SliceTree[K, V]) Remove(k K) bool {

	idx, offset := s.GetIndex(k)
	return s.clearIdx(idx, offset)
}

// Sets the key/vale pair and returns the index id.
// Comlexity: o(log n)
func (s *SliceTree[K, V]) Put(k K, v V) (index int) {
	total := len(s.Slices)

	if total == 0 {
		// 0 size.. just append
		s.Slices = append(s.Slices, &KvSet[K, V]{k, v})
		return 0
	}
	idx, offset := s.GetIndex(k)
	return s.SetIndex(idx, offset, k, v)
}

// Sets the value in the index to v.  The last index value returned from Put to update the last index point.
// This lets you bypass the o(log n) update complexity for writing to the same element over and over again.
// The internals still call s.OnOverWrite for you.
func (s *SliceTree[K, V]) Set(index int, v V) (status bool) {
	if index < 0 || len(s.Slices) <= index {
		return

	}
	el := s.Slices[index]
	el.Value = s.OnOverWrite(el.Key, el.Value, v)
	status = true
	return
}

// Tries to fetch value based on key of k, if k does not exist, found is false.
func (s *SliceTree[K, V]) Get(k K) (value V, found bool) {
	if len(s.Slices) == 0 {
		return
	}
	i, o := s.GetIndex(k)
	if o == 0 {
		return s.Slices[i].Value, true
	}
	return
}

// Returns true if the k exists in the slcie.
// Complexity: o(log n)
func (s *SliceTree[K, V]) Exists(k K) bool {
	size := len(s.Slices)
	switch size {
	case 0:
		return false
	case 1:
		return s.Cmp(s.Slices[0].Key, k) == 0
	}
	_, o := s.GetIndex(k)
	return o == 0
}

func (s *SliceTree[K, V]) clearIdx(idx, offset int) (result bool) {

	size := len(s.Slices)
	if offset != 0 || size == 0 || idx >= size || idx < 0 {
		result = false
	} else if size == 1 {
		// single element
		if idx == 0 {
			s.Slices = s.Slices[:0]
			return true
		}
	} else {
		s.Slices = slices.Delete(s.Slices, idx, idx+1)
		return true
	}

	return false
}

// Removes all elements in the slice, but keeps the memory allocated.
func (s *SliceTree[K, V]) RemoveAll() int {
	t := len(s.Slices)
	s.Slices = s.Slices[:0]
	return t
}

func (s *SliceTree[K, V]) clearTo(key K, x int, cb func(a, b int)) {
	i, o := s.GetIndex(key)
	index := i + o + x
	end := index + 1
	size := len(s.Slices)
	if index < 0 {
		cb(0, 0)
		return
	} else if end > size {
		cb(0, size)
		s.Slices = s.Slices[:0]
		return
	}

	cb(0, end)

	if end == size {
		s.Slices = s.Slices[:0]
	} else {
		s.Slices = s.Slices[end:size]
	}
}

func (s *SliceTree[K, V]) clearFrom(key K, x int, cb func(a, b int)) {
	i, o := s.GetIndex(key)
	index := i + o + x
	end := index + 2
	size := len(s.Slices)
	if index <= 0 {
		cb(0, size)
		s.Slices = s.Slices[:0]
		return
	} else if index >= size {
		cb(0, 0)
		return
	}
	if end > size {
		end = size
	}

	cb(index, end)
	s.Slices = s.Slices[0:index]
}

// RemoveFrom implements [OrderedMapExt].
func (s *SliceTree[K, V]) RemoveFrom(key K) (total int) {
	s.clearFrom(key, 0, func(a, b int) {
		total = b - a
	})
	return
}

// RemoveFromS implements [OrderedMapExt].
func (s *SliceTree[K, V]) RemoveFromS(key K) (result []*KvSet[K, V]) {
	s.clearFrom(key, 0, func(a, b int) {
		result = s.Slices[a:b]
	})
	return
}

// RemoveAfterS implements [OrderedMapExt].
func (s *SliceTree[K, V]) RemoveAfterS(key K) (result []*KvSet[K, V]) {
	s.clearFrom(key, 1, func(a, b int) {
		result = s.Slices[a:b]
	})
	return
}

// RemoveAfter implements [OrderedMapExt].
func (s *SliceTree[K, V]) RemoveAfter(key K) (total int) {
	s.clearFrom(key, 1, func(a, b int) {
		total = b - a
	})
	return
}

// RemoveTo implements [OrderedMapExt].
func (s *SliceTree[K, V]) RemoveTo(key K) (total int) {
	s.clearTo(key, 0, func(a, b int) {
		total = b - a
	})
	return
}

// RemoveBefore implements [OrderedMapExt].
func (s *SliceTree[K, V]) RemoveBefore(key K) (total int) {
	s.clearTo(key, -1, func(a, b int) {
		total = b - a
	})
	return
}

// RemoveToS implements [OrderedMapExt].
func (s *SliceTree[K, V]) RemoveToS(key K) (result []*KvSet[K, V]) {
	s.clearTo(key, 0, func(a, b int) {
		total := b - a
		if total == 0 {
			return
		}
		result = s.Slices[0:b]
	})
	return
}

// RemoveBeforeS implements [OrderedMapExt].
func (s *SliceTree[K, V]) RemoveBeforeS(key K) (result []*KvSet[K, V]) {
	s.clearTo(key, -1, func(a, b int) {
		total := b - a
		if total == 0 {
			return
		}
		result = s.Slices[0:b]
	})
	return
}

// Returns the total number key/value pairs in the slice.
func (s *SliceTree[K, V]) Size() int {
	return len(s.Slices)
}

// Sets the given k,v pair based on the index and offset provided by a call to GetIndex.
// Returns the resulting array index id.
//
// Using a combinaiton of GetIndex and SetIndex lets you bypass the o(log n) comlexity when wiring to the same node over and over again.
// The value reutrned from Put can be used to update the internals using SetIndex with the offset being 0.
func (s *SliceTree[K, V]) SetIndex(idx, offset int, k K, v V) (index int) {
	size := len(s.Slices)
	if offset != 0 {
		ns := size + 1
		s.grow(ns)
		s.Slices = append(s.Slices, nil)
		kv := &KvSet[K, V]{k, v}
		switch idx {
		case 0:
			if offset == 1 {
				copy(s.Slices[2:], s.Slices[1:size])
				s.Slices[1] = kv
				return 1
			} else {
				copy(s.Slices[1:], s.Slices[0:size])
				s.Slices[0] = kv
			}
			return 0
		default:
			index = idx + offset
			copy(s.Slices[index+1:], s.Slices[index:size])
			if offset < 0 {
				s.Slices[idx] = kv
				return idx
			} else {
				s.Slices[index] = kv
				return index
			}

		}
	} else {
		ns := idx + 1
		if size < ns {
			// empty slice!
			s.grow(ns)
			s.Slices = s.Slices[:ns]
			s.Slices[idx] = &KvSet[K, V]{k, v}
		} else {
			// overwrite
			s.Slices[idx].Value = s.OnOverWrite(k, s.Slices[idx].Value, v)
		}

		return idx
	}
}

func (s *SliceTree[K, V]) grow(size int) {
	if cap(s.Slices) < size {
		s.Slices = slices.Grow(s.Slices, s.Growth)
	}
}

// Returns the index and offset of a given key.
//
// The index is the current relative postion in the slice.
//
// The offset represents where the item would be placed:
//   - offset of 0, at index value.
//   - offset of 1, expand the slice after the inddex and put the value to the right of the index
//   - offset of -1, expand the slice before the index and put the value to left of the current postion
//
// Complexity: o(log n)
func (s *SliceTree[K, V]) GetIndex(k K) (index, offset int) {

	size := len(s.Slices)
	switch size {
	case 0:
		return 0, 0
	case 1:
		return 0, s.Cmp(k, s.Slices[0].Key)
	}
	nextMid := s.getMid(size)
	// well if we get here.. we need to walk the tree
	nextBegin := 0
	nextEnd := len(s.Slices) - 1
	var resolved bool

	for {
		nextBegin, nextEnd, nextMid, offset, resolved = s.resolveNext(nextBegin, nextEnd, nextMid, k)
		if resolved {
			index = offset + nextMid
			if index < 0 {
				return nextMid, offset
			}
			offset = s.Cmp(k, s.Slices[index].Key)
			break
		}
	}
	return
}

func (s *SliceTree[K, V]) lookAhead(end, mid int, k K) (nextBegin, nextEnd, nextMid, offset int, resolved bool) {
	nextBegin = mid + 1
	diff := end - nextBegin

	if diff <= 0 {
		resolved = true
		nextMid = mid
		offset = 1
		return
	}
	nextMid = nextBegin + s.getMid(diff+1)
	nextEnd = end
	offset = s.Cmp(s.Slices[nextMid].Key, k)
	resolved = offset == 0
	return
}

func (s *SliceTree[K, V]) lookBehind(begin, mid int, k K) (nextBegin, nextEnd, nextMid, offset int, resolved bool) {
	nextEnd = mid - 1
	diff := nextEnd - begin

	if diff <= 0 {
		resolved = true
		nextMid = mid
		offset = -1
		return
	}
	nextMid = begin + s.getMid(diff+1)
	nextBegin = begin
	offset = s.Cmp(s.Slices[nextMid].Key, k)
	resolved = offset == 0
	return
}

func (s *SliceTree[K, V]) resolveNext(begin, end, mid int, k K) (nextBegin, nextEnd, nextMid, offset int, resolved bool) {

	cmp := s.Cmp(k, s.Slices[mid].Key)
	switch cmp {
	case 0:
		nextMid = mid
		resolved = true
		offset = cmp
		return
	case -1:
		nextBegin, nextEnd, nextMid, offset, resolved = s.lookBehind(begin, mid, k)
	case 1:
		nextBegin, nextEnd, nextMid, offset, resolved = s.lookAhead(end, mid, k)
	}

	return
}

// Returns an iterator for the current keys.
// The internals of this iterator  do not lock the tree or prevent updates.  You can safely call an iterator from with an iterator.
// and not run into deadlocks.
func (s *SliceTree[K, V]) Keys() iter.Seq2[int, K] {
	pos := 0
	return func(yield func(int, K) bool) {
		for pos < len(s.Slices) {
			if !yield(pos, s.Slices[pos].Key) {
				return
			}
			pos++
		}
	}
}

// Returns an iterator for the current values
// The internals of this iterator  do not lock the tree or prevent updates.  You can safely call an iterator from with an iterator.
// and not run into deadlocks.
func (s *SliceTree[K, V]) Values() iter.Seq2[int, V] {

	pos := 0
	return func(yield func(int, V) bool) {
		for pos < len(s.Slices) {
			if !yield(pos, s.Slices[pos].Value) {
				return
			}
			pos++
		}
	}
}

// Returns an iterator for key/value pars.
// The internals of this iterator  do not lock the tree or prevent updates.  You can safely call an iterator from with an iterator.
// and not run into deadlocks.
func (s *SliceTree[K, V]) All() iter.Seq2[K, V] {
	pos := 0
	return func(yield func(K, V) bool) {
		for pos < len(s.Slices) {
			kv := s.Slices[pos]
			if !yield(kv.Key, kv.Value) {
				return
			}
			pos++
		}
	}
}

// Attempts to remove the keys from the tree in bulk.  Returns the number of keys removed.
//
// This is almost always faster than just looping over a list of keys and calling Remove one key at a time. The internals of
// The MassRemove method deletes elements in squential contiguys blocks: reducing on the number of internal splice operations.
//
// Complexity:
//
// Worst case shown (Per key removed or k): o(log(n) + o(log k) + 2*k).
//
// In truth, the real world complexity is drastically reduced by the following:
//   - Deletetion of duplicate keys.
//   - Keys being deleted do not exist.
//
// The complexity is defined by the steps required:
//   - Index lookups: o(log n) +k
//   - child index is created to de-duplicate and order keys for deletion: o(log k)
//   - key deletion is done in contiguous blocks: k
func (s *SliceTree[K, V]) MassRemove(args ...K) (total int) {
	if len(s.Slices) == 0 {
		return 0
	}
	f := New[int, any](reverse)
	for _, k := range args {
		i, o := s.GetIndex(k)
		if o != 0 {
			continue
		}
		f.Put(i, nil)
	}

	total = f.Size()
	s.contig(total, f.Keys(), s.rangedel)

	return
}

// This method is by defenition, unsafe, but fast.
//
// Only use if the keys being removed meet all of the following requirements:
//   - no duplicate keys
//   - keys are in ascending ordered
//   - all keys currently exist in the internals of the tree
//
// Complexity: key=keys; o(log n +k)
func (s *SliceTree[K, V]) UnSafeMassRemove(keys ...K) {
	s.contig(len(keys), s.unsafeIter(keys), s.rangedel)
}

func (s *SliceTree[K, V]) unsafeIter(keys []K) iter.Seq2[int, int] {
	pos := len(keys) - 1
	id := 0
	return func(yield func(int, int) bool) {
		for pos > -1 {
			i, _ := s.GetIndex(keys[pos])
			if !yield(id, i) {
				return
			}
			id++
			pos--
		}
	}
}

func KvIter[K any, V any](set []*KvSet[K, V]) iter.Seq2[K, V] {

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

func (s *SliceTree[K, V]) contig(totalKeys int, r iter.Seq2[int, int], cb func(a, b int)) {
	var end = -1
	var last = -1
	size := totalKeys - 1
	for i, idx := range r {
		p := idx + 1
		if last == p {
			if end == -1 {
				end = last
			}
			if i == size {
				cb(idx, end)
			}

		} else if last != -1 && end == -1 {

			cb(last, last)
			if i == size {
				cb(idx, idx)
			}
		} else {
			if end != -1 {
				cb(last, end)
			} else if i == size {
				cb(idx, idx)
			}
			end = -1
		}

		last = idx
	}
}

// ThreadSafe implements [OrderedMap]
func (s *SliceTree[K, V]) ThreadSafe() bool { return false }

// RemoveBeforeI implements [OrderedMapExt]
func (s *SliceTree[K, V]) RemoveBeforeI(key K) iter.Seq2[K, V] {
	return KvIter(s.RemoveBeforeS(key))
}

// RemoveFromI implements [OrderedMapExt]
func (s *SliceTree[K, V]) RemoveFromI(key K) iter.Seq2[K, V] {
	return KvIter(s.RemoveFromS(key))
}

// RemoveAfterI implements [OrderedMapExt]
func (s *SliceTree[K, V]) RemoveAfterI(key K) iter.Seq2[K, V] {
	return KvIter(s.RemoveAfterS(key))
}

// RemoveAfterI implements [OrderedMapExt]
func (s *SliceTree[K, V]) RemoveToI(key K) iter.Seq2[K, V] {
	return KvIter(s.RemoveToS(key))
}

// Returns a thread safe instnace from the current instance.
func (s *SliceTree[K, V]) ToTs() OrderedMap[K, V] {
	return &ThreadSafeOrderedMap[K, V]{Tree: s}
}

// GetFirstKey implements [OrderedMap]
func (s *SliceTree[K, V]) FirstKey() (key K, ok bool) {
	if s.Size() == 0 {
		return
	}
	key = s.Slices[0].Key
	ok = true
	return
}

// GetLastKey implements [OrderedMap]
func (s *SliceTree[K, V]) LastKey() (key K, ok bool) {
	if s.Size() == 0 {
		return
	}
	key = s.Slices[len(s.Slices)-1].Key
	ok = true
	return
}

func (s *SliceTree[K, V]) betweenIter(a, b int) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if a < 0 {
			return
		}

		pos := a
		for pos < len(s.Slices) && pos <= b {
			kv := s.Slices[pos]
			if !yield(kv.Key, kv.Value) {
				return
			}
			pos++
		}
	}
}

func (s *SliceTree[K, V]) betweenChecks(a, b K) (begin, end, total int, ok bool) {
	if s.Size() == 0 {
		return
	}

	begin, c := s.GetIndex(a)
	//begin = i + o
	offset := c
	end, d := s.GetIndex(b)
	offset += d
	//end = i + o

	size := s.Size()
	final := size - 1
	if offset*offset == 4 && ((begin+end == final*2) || (begin+end == 0)) {
		// completly out of our ragne
		return
	}

	if d < 1 {
		end += d
	}
	if c > 0 {
		begin += c
	}

	if begin > end {
		return
	}

	total = 1 + end - begin
	ok = true

	return
}

// Between implements [OrderedMap]
func (s *SliceTree[K, V]) Between(a, b K) (total int, ok bool) {
	_, _, total, ok = s.betweenChecks(a, b)
	return
}

// BetweenKV implements [OrderedMap]
func (s *SliceTree[K, V]) BetweenKV(a, b K) (seq iter.Seq2[K, V]) {
	x, y, _, ok := s.betweenChecks(a, b)
	if ok {
		return s.betweenIter(x, y)
	} else {
		return s.betweenIter(-1, -1)
	}

}

func (s *SliceTree[K, V]) RemoveBetween(a, b K) (total int) {
	s.clearBetween(a, b, func(x, y, t int, ok bool) {
		if ok {
			total = 1 + y - x
		}
	})
	return
}

func (s *SliceTree[K, V]) clearBetween(a, b K, cb func(x, y, t int, ok bool)) {
	begin, end, total, ok := s.betweenChecks(a, b)
	cb(begin, end, total, ok)
	if ok {
		s.Slices = slices.Delete(s.Slices, begin, 1+end)
	}
}
