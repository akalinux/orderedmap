package omap

import (
	"iter"
	"sync"
)

type ThreadSafeOrderedMap[K any, V any] struct {
	// Instance to wrap for locking
	Tree *SliceTree[K, V]
	lock sync.RWMutex
}

// Creates a new thread safe OrderedMap.
func NewTs[K any, V any](Cmp func(a, b K) int) (Map OrderedMap[K, V]) {

	Map = &ThreadSafeOrderedMap[K, V]{
		Tree: New[K, V](Cmp),
	}
	return
}

// All implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		pos := 0
		s.lock.RLock()
		defer s.lock.RUnlock()
		for {
			k, v, ok := s.nextKv(pos)
			if !ok {
				return
			}
			if !yield(k, v) {
				return
			}
			pos++
		}
	}
}

// RemoveAfter implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) RemoveAfter(key K) (total int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveAfter(key)
}

// RemoveAfterS implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) RemoveAfterS(key K) (result []*KvSet[K, V]) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveAfterS(key)
}

// RemoveAll implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) RemoveAll() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveAll()
}

// RemoveBefore implements [OrderedMapExt].
func (s *ThreadSafeOrderedMap[K, V]) RemoveBefore(key K) (total int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveBefore(key)
}

// RemoveBeforeS implements [OrderedMapExt].
func (s *ThreadSafeOrderedMap[K, V]) RemoveBeforeS(key K) (result []*KvSet[K, V]) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveBeforeS(key)
}

// Creates a thread safe iterator for a slice of *KvSet.
func TsKvIter[K any, V any](set []*KvSet[K, V]) iter.Seq2[K, V] {
	var lock sync.RWMutex
	return func(yield func(K, V) bool) {
		lock.RLock()
		defer lock.RUnlock()
		for _, row := range set {
			if !yield(row.Key, row.Value) {
				return
			}
		}
	}
}

// RemoveBeforeI implements [OrderedMapExt].
// Returns a thread safe iterator for the deleted values.
func (s *ThreadSafeOrderedMap[K, V]) RemoveBeforeI(key K) iter.Seq2[K, V] {
	return TsKvIter(s.RemoveBeforeS(key))
}

// RemoveFromI implements [OrderedMapExt].
// Returns a thread safe iterator for the deleted values.
func (s *ThreadSafeOrderedMap[K, V]) RemoveFromI(key K) iter.Seq2[K, V] {
	return TsKvIter(s.RemoveFromS(key))
}

// RemoveAfterI implements [OrderedMapExt].
// Returns a thread safe iterator for the deleted values.
func (s *ThreadSafeOrderedMap[K, V]) RemoveAfterI(key K) iter.Seq2[K, V] {
	return TsKvIter(s.RemoveAfterS(key))
}

// RemoveRemoveToI implements [OrderedMapExt].
// Returns a thread safe iterator for the deleted values.
func (s *ThreadSafeOrderedMap[K, V]) RemoveToI(key K) iter.Seq2[K, V] {
	return TsKvIter(s.RemoveToS(key))
}

// RemoveFrom implements [OrderedMapExt].
func (s *ThreadSafeOrderedMap[K, V]) RemoveFrom(key K) (total int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveFrom(key)
}

// RemoveFromS implements [OrderedMapExt].
func (s *ThreadSafeOrderedMap[K, V]) RemoveFromS(key K) (result []*KvSet[K, V]) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveFromS(key)
}

// RemoveTo implements [OrderedMapExt].
func (s *ThreadSafeOrderedMap[K, V]) RemoveTo(key K) (total int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveTo(key)
}

// RemoveToS implements [OrderedMapExt].
func (s *ThreadSafeOrderedMap[K, V]) RemoveToS(key K) (result []*KvSet[K, V]) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveToS(key)
}

// Exists implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Exists(key K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Tree.Exists(key)
}

// Get implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Get(key K) (value V, found bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Tree.Get(key)
}

// Keys implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Keys() iter.Seq2[int, K] {
	return func(yield func(int, K) bool) {
		pos := 0
		s.lock.RLock()
		defer s.lock.RUnlock()
		for {
			k, _, ok := s.nextKv(pos)
			if !ok {
				return
			}
			if !yield(pos, k) {
				return
			}
			pos++
		}
	}
}

// MassRemove implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) MassRemove(keys ...K) (total int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.MassRemove(keys...)
}

// Put implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Put(key K, value V) (index int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.Put(key, value)
}

// Remove implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Remove(key K) bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.Remove(key)
}

// Size implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Tree.Size()
}

// Values implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Values() iter.Seq2[int, V] {
	return func(yield func(int, V) bool) {
		pos := 0
		s.lock.RLock()
		defer s.lock.RUnlock()
		for {
			_, v, ok := s.nextKv(pos)
			if !ok {
				return
			}
			if !yield(pos, v) {
				return
			}
			pos++
		}
	}
}

// Always returns true.
func (s *ThreadSafeOrderedMap[K, V]) ThreadSafe() bool {
	return true
}

func (s *ThreadSafeOrderedMap[K, V]) nextKv(pos int) (key K, value V, ok bool) {
	if s.Tree.Size() > pos && pos > -1 {
		ok = true
		key = s.Tree.Slices[pos].Key
		value = s.Tree.Slices[pos].Value
	}
	return
}

// GetFirstKey implements [OrderedMap]
func (s *ThreadSafeOrderedMap[K, V]) FirstKey() (key K, ok bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	key, ok = s.Tree.FirstKey()
	return
}

// GetLastKey implements [OrderedMap]
func (s *ThreadSafeOrderedMap[K, V]) LastKey() (key K, ok bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	key, ok = s.Tree.LastKey()
	return
}

// Between implements [OrderedMap]
func (s *ThreadSafeOrderedMap[K, V]) Between(a, b K) (total int, ok bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	total, ok = s.Tree.Between(a, b)
	return
}

// BetweenKV implements [OrderedMap]
func (s *ThreadSafeOrderedMap[K, V]) BetweenKV(a, b K) (seq iter.Seq2[K, V]) {
	return func(yield func(K, V) bool) {
		s.lock.RLock()
		defer s.lock.RUnlock()
		next, stop := iter.Pull2(s.Tree.BetweenKV(a, b))
		defer stop()
		for k, v, ok := next(); ok; k, v, ok = next() {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (s *ThreadSafeOrderedMap[K, V]) RemoveBetween(a, b K) (total int) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Tree.RemoveBetween(a, b)
}
