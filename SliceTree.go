package omap

import (
	"iter"
	"slices"
)

const SUB_ONE int8 = -1

type SliceTree[K any, V any] struct {

	// Internally managed keys slice
	Slices []KvSet[K, V]

	// Compare function.
	Cmp func(a, b K) int

	// Required non 0 value, determins by what capacity we grow the internal
	// Slice.  Default is 1.
	Growth int

	// Required non nil value, called when ever a value is overwritten.
	// Seting this funtion saves on having to write a check when data is overwritten.
	OnOverWrite func(key K, oldValue V, newValue V)
}

// Creatss a new SliceTree with the internal Slice set to "size".
func NewSliceTree[K any, V any](size int, cb func(a, b K) int) *SliceTree[K, V] {
	return &SliceTree[K, V]{
		Slices: make([]KvSet[K, V], 0, size),
		Cmp:    cb,
		Growth: 100,
	}
}

// Creates a new SliceTee with the default Slices size of 100.  If you require more control over the starting size of the slice
// use the NewSliceTree function in stead.
func New[K any, V any](cb func(a, b K) int) *SliceTree[K, V] {
	return NewSliceTree[K, V](100, cb)
}

func NewFromMap[K comparable, V any](m map[K]V, cb func(a, b K) int) *SliceTree[K, V] {
	s := NewSliceTree[K, V](100, cb)
	for k, v := range m {
		s.Put(k, v)
	}
	return s
}

// Tries to remove the element of k, returns false if it fails.
// Complexity: o(log n)
func (s *SliceTree[K, V]) Remove(k K) (value V, ok bool) {

	idx, offset := GetIndex(k, s.Cmp, s.Slices)
	i := idx + int(offset)
	if i >= 0 && i < len(s.Slices) {
		value = s.Slices[i].Value
	}

	ok = s.clearIdx(idx, offset)
	return
}

// Sets the key/vale pair and returns the index id.
// Comlexity: o(log n)
func (s *SliceTree[K, V]) Put(k K, v V) {
	total := len(s.Slices)

	if total == 0 {
		// 0 size.. just append
		s.Slices = append(s.Slices, KvSet[K, V]{k, v})
	}
	idx, offset := GetIndex(k, s.Cmp, s.Slices)
	s.SetIndex(idx, offset, k, v)
}

func (s *SliceTree[K, V]) Contains(key K) bool {
	if s.Size() == 0 {
		return false
	} else if s.Cmp(s.Slices[0].Key, key) == 1 {
		return false
	} else if s.Cmp(s.Slices[len(s.Slices)-1].Key, key) == -1 {
		return false
	}

	return true
}

// Sets the value in the index to v.  The last index value returned from Put to update the last index point.
// This lets you bypass the o(log n) update complexity for writing to the same element over and over again.
// The internals still call s.OnOverWrite for you.
func (s *SliceTree[K, V]) Set(index int, v V) (status bool) {
	if index < 0 || len(s.Slices) <= index {
		return

	}
	el := &s.Slices[index]
	if s.OnOverWrite != nil {
		s.OnOverWrite(el.Key, el.Value, v)
	}
	el.Value = v

	status = true
	return
}

// Tries to fetch value based on key of k, if k does not exist, found is false.
func (s *SliceTree[K, V]) Get(k K) (value V, found bool) {
	if len(s.Slices) == 0 {
		return
	}
	i, o := GetIndex(k, s.Cmp, s.Slices)
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
	_, o := GetIndex(k, s.Cmp, s.Slices)
	return o == 0
}

func (s *SliceTree[K, V]) clearIdx(idx int, offset int8) (result bool) {

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

// Returns the total number key/value pairs in the slice.
func (s *SliceTree[K, V]) Size() int {
	return len(s.Slices)
}

// Sets the given k,v pair based on the index and offset provided by a call to GetIndex.
// Returns the resulting array index id.
//
// Using a combinaiton of GetIndex and SetIndex lets you bypass the o(log n) comlexity when wiring to the same node over and over again.
// The value reutrned from Put can be used to update the internals using SetIndex with the offset being 0.
func (s *SliceTree[K, V]) SetIndex(idx int, offset int8, k K, v V) (index int) {
	size := len(s.Slices)
	if offset != 0 {
		ns := size + 1
		s.grow(ns)
		s.Slices = s.Slices[:ns]
		os := (offset + 1) >> 1
		switch idx {
		case 0:
			copy(s.Slices[1+os:], s.Slices[os:size])
			s.Slices[os] = KvSet[K, V]{k, v}
			return int(os)
		default:
			index = idx + int(os)
			copy(s.Slices[index+1:], s.Slices[index:size])
			s.Slices[index] = KvSet[K, V]{k, v}
			return
		}
	} else {
		ns := idx + 1
		if size < ns {
			// empty slice!
			s.grow(ns)
			s.Slices = s.Slices[:ns]
			s.Slices[idx] = KvSet[K, V]{k, v}
		} else {
			// overwrite
			if s.OnOverWrite != nil {
				s.OnOverWrite(k, s.Slices[idx].Value, v)
			}
			s.Slices[idx].Value = v
		}

		return idx
	}
}

func (s *SliceTree[K, V]) grow(size int) {
	if cap(s.Slices) < size {
		grow := s.Growth
		// be nice to people who make an empty struct
		if grow <= 0 {
			grow = 1
		}
		s.Slices = slices.Grow(s.Slices, grow)
	}
}

// Returns an iterator for the current keys.
// The internals of this iterator  do not lock the tree or prevent updates.  You can safely call an iterator from with an iterator.
// and not run into deadlocks.
func (s *SliceTree[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		size := len(s.Slices)
		for pos := 0; pos < size; pos++ {
			if !yield(s.Slices[pos].Key) {
				return
			}
		}
	}
}

// Returns an iterator for the current values
// The internals of this iterator  do not lock the tree or prevent updates.  You can safely call an iterator from with an iterator.
// and not run into deadlocks.
func (s *SliceTree[K, V]) Values() iter.Seq[V] {

	return func(yield func(V) bool) {
		size := len(s.Slices)
		for pos := 0; pos < size; pos++ {
			if !yield(s.Slices[pos].Value) {
				return
			}
		}
	}
}

// Returns an iterator for key/value pars.
// The internals of this iterator  do not lock the tree or prevent updates.  You can safely call an iterator from with an iterator.
// and not run into deadlocks.
func (s *SliceTree[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		size := len(s.Slices)
		for pos := 0; pos < size; pos++ {
			kv := s.Slices[pos]
			if !yield(kv.Key, kv.Value) {
				return
			}
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
	return s.massremove(args, s.rangedel)
}

func (s *SliceTree[K, V]) massremove(args []K, cb func(a, b int)) int {
	if len(s.Slices) == 0 {
		return 0
	}
	f := New[int, any](reverse)
	for _, k := range args {
		i, o := GetIndex(k, s.Cmp, s.Slices)
		if o != 0 {
			continue
		}
		f.Put(i, nil)
	}

	s.contig(f.Size(), f.keys(), cb)
	return f.Size()
}

func (s *SliceTree[K, V]) keys() iter.Seq2[int, K] {
	return func(yield func(int, K) bool) {
		size := len(s.Slices)
		for pos := 0; pos < size; pos++ {
			if !yield(pos, s.Slices[pos].Key) {
				return
			}
		}
	}
}

func (s *SliceTree[K, V]) MassRemoveKV(args ...K) iter.Seq2[K, V] {
	// worst case is the len of args.. so we always pre-allocate for worst case
	list := make([][]KvSet[K, V], 0, len(args))
	s.massremove(args, func(a, b int) {
		list = append(list, slices.Clone(s.Slices[a:b+1]))
		s.rangedel(a, b)
	})
	return func(yield func(K, V) bool) {
		/// outer loop is in reverse order
		for i := len(list) - 1; i > -1; i-- {
			// inner loop is in forward order
			for _, kv := range list[i] {
				if !yield(kv.Key, kv.Value) {
					return
				}
			}
		}
	}

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

func (s *SliceTree[K, V]) FilterBetween(cb func(K, V) bool, a, b K, opt ...int) {
	begin, end, _, ok := s.betweenChecks(a, b, opt...)
	if !ok {
		return
	}
	count := 0
	Slices := s.Slices
	size := len(Slices)
	filtered := make([]KvSet[K, V], end-begin+1)
	end += 1
	offset := size - end
	for _, kv := range Slices[begin:end] {
		if !cb(kv.Key, kv.Value) {
			filtered[count] = kv
			count++
		}
	}
	copy(Slices[begin:begin+count], filtered[0:count])
	if offset != 0 {
		copy(Slices[begin+count:], Slices[end:size])
	}
	s.Slices = Slices[0 : begin+count+offset]

}

func (s *SliceTree[K, V]) unsafeIter(keys []K) iter.Seq2[int, int] {
	pos := len(keys) - 1
	id := 0
	return func(yield func(int, int) bool) {
		for pos > -1 {
			i, _ := GetIndex(keys[pos], s.Cmp, s.Slices)
			if !yield(id, i) {
				return
			}
			id++
			pos--
		}
	}
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

// Merge implements [OrderedMap]
func (s *SliceTree[K, V]) Merge(set OrderedMap[K, V]) int {
	// do not add to ourself!
	if s == set {
		return 0
	}

	size := s.Size()
	// is this worth trying to optimize?
	for k, v := range set.All() {
		s.Put(k, v)
	}

	return s.Size() - size
}

func (s *SliceTree[K, V]) FastMerge(set OrderedMap[K, V]) int {
	src := set.GetKvSlice()
	size := s.Size()
	if set.Size() == 0 {
		return 0
	}
	if size == 0 {
		copy(s.Slices, src)
		return set.Size()
	}

	res, end := MergeKvSet(s.Slices, src, make([]KvSet[K, V], s.Size()+len(src)), 0, s.Size()-1, s.Cmp, s.OnOverWrite)
	s.Slices = res[0 : end+1]
	return s.Size() - size
}

func (s *SliceTree[K, V]) GetKvSlice() []KvSet[K, V] {
	return s.Slices
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

func (s *SliceTree[K, V]) betweenChecks(a, b K, opt ...int) (begin, end, total int, ok bool) {
	if s.Size() == 0 {
		return
	}

	var c int8
	var d int8
	if len(opt) == 0 {
		begin, c = GetIndex(a, s.Cmp, s.Slices)
		end, d = GetIndex(b, s.Cmp, s.Slices)
	} else {
		if FIRST_KEY != opt[0]&FIRST_KEY {
			begin, c = GetIndex(a, s.Cmp, s.Slices)
		} else {
			c = -1
		}
		if LAST_KEY == opt[0]&LAST_KEY {
			end = len(s.Slices) - 1
		} else {
			end, d = GetIndex(b, s.Cmp, s.Slices)
		}
	}
	offset := c
	offset += d

	size := s.Size()
	final := size - 1
	if offset*offset == 4 && ((begin+end == final*2) || (begin+end == 0)) {
		return
	}

	if d < 1 {
		end += int(d)
	}
	if c > 0 {
		begin += int(c)
	}

	if begin > end {
		return
	}

	total = 1 + end - begin
	ok = true

	return
}

// Between implements [OrderedMap]
func (s *SliceTree[K, V]) Between(a, b K, opt ...int) (total int) {
	_, _, total, _ = s.betweenChecks(a, b, opt...)
	return
}

// BetweenKV implements [OrderedMap]
func (s *SliceTree[K, V]) BetweenKV(a, b K, opt ...int) (seq iter.Seq2[K, V]) {
	x, y, _, ok := s.betweenChecks(a, b, opt...)
	if ok {
		return s.betweenIter(x, y)
	} else {
		return s.betweenIter(-1, -1)
	}

}

// Returns a slice containing the elements between a and b
func (s *SliceTree[K, V]) GetBetweenKvSlice(a, b K, opt ...int) []KvSet[K, V] {
	x, y, _, ok := s.betweenChecks(a, b, opt...)
	if ok {
		return s.Slices[x : y+1]
	} else {
		return []KvSet[K, V]{}
	}
}

// RemoveBetween implements [OrderedMap]
func (s *SliceTree[K, V]) RemoveBetween(a, b K, opt ...int) (total int) {
	s.clearBetween(a, b, func(x, y, t int, res bool) {
		if res {
			total = 1 + y - x
		}
	}, opt...)
	return
}

// RemoveBetween implements [OrderedMap]
func (s *SliceTree[K, V]) RemoveBetweenKV(a, b K, opt ...int) (removed iter.Seq2[K, V]) {
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

func (s *SliceTree[K, V]) clearBetween(a, b K, cb func(x, y, t int, ok bool), opt ...int) {
	begin, end, total, ok := s.betweenChecks(a, b, opt...)
	cb(begin, end, total, ok)
	if ok {
		s.Slices = slices.Delete(s.Slices, begin, 1+end)
	}
}

// Returns false.
func (s *SliceTree[K, V]) ThreadSafe() bool {
	return false
}

// Sets the internal OnOverWrite function.
func (s *SliceTree[K, V]) SetOverwrite(cb func(key K, oldValue, newValue V)) {
	s.OnOverWrite = cb
}

// SetGrowth implements [OrderedMap]
func (s *SliceTree[K, V]) SetGrowth(grow int) {
	if grow <= 0 {
		s.Growth = 1
		return
	}
	s.Growth = grow
}

// Utility method, to convert from an OrderedMap instance to a regular map.
// Due to constraints placed on maps in go, this feature is implemented as a function, not a method.
func ToMap[K comparable, V any](src OrderedMap[K, V]) map[K]V {

	m := make(map[K]V, src.Size())
	for k, v := range src.All() {
		m[k] = v
	}
	return m
}

// Deletes the given element when the callback returns true
func (s *SliceTree[K, V]) Filter(cb func(K, V) bool) {
	ns := make([]KvSet[K, V], 0, s.Size())
	for k, v := range s.All() {
		if !cb(k, v) {
			ns = append(ns, KvSet[K, V]{k, v})
		}
	}
	s.Slices = ns
}
