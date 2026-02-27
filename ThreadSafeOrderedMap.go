package omap

import (
	"iter"
	"sync"
)

type ThreadSafeOrderedMap[K any, V any] struct {
	// Instance to wrap for locking
	Tree OrderedMap[K, V]
	lock sync.RWMutex
}

// Creates a new thread safe [OrderedMap] instance.
func NewTs[K any, V any](Cmp func(a, b K) int) (Map OrderedMap[K, V]) {

	Map = &ThreadSafeOrderedMap[K, V]{
		Tree: New[K, V](Cmp),
	}
	return
}

// Always returns this instance.
func (s *ThreadSafeOrderedMap[K, V]) ToTs() OrderedMap[K, V] {
	return s
}

// All implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		s.lock.RLock()
		defer s.lock.RUnlock()
		next, stop := iter.Pull2(s.Tree.All())
		defer stop()
		for k, v, ok := next(); ok; k, v, ok = next() {
			if !yield(k, v) {
				return
			}
		}

	}
}

// RemoveAll implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) RemoveAll() int {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.RemoveAll()
}

// Exists implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Exists(key K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Tree.Exists(key)
}
func (s *ThreadSafeOrderedMap[K, V]) Contains(key K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Tree.Contains(key)
}

// Get implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Get(key K) (value V, found bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Tree.Get(key)
}

type rlocker interface {
	rlock()
	runlock()
}

func (s *ThreadSafeOrderedMap[K, V]) rlock() {
	s.lock.RLock()
}
func (s *ThreadSafeOrderedMap[K, V]) runlock() {
	s.lock.RUnlock()
}

func invertSeq[K any](s rlocker, cb func() iter.Seq[K]) iter.Seq[K] {
	return func(yield func(K) bool) {
		s.rlock()
		defer s.runlock()
		next, stop := iter.Pull(cb())
		defer stop()
		for k, ok := next(); ok; k, ok = next() {
			if !yield(k) {
				return
			}
		}
	}
}

// Keys implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Keys() iter.Seq[K] {
	return invertSeq(s, s.Tree.Keys)
}

// MassRemove implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) MassRemove(keys ...K) (total int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.MassRemove(keys...)
}

func (s *ThreadSafeOrderedMap[K, V]) MassRemoveKV(keys ...K) iter.Seq2[K, V] {
	s.lock.Lock()
	defer s.lock.Unlock()
	return TsKvIterWrapper(s.Tree.MassRemoveKV(keys...))
}

func (s *ThreadSafeOrderedMap[K, V]) Filter(cb func(K, V) bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Tree.Filter(cb)
}

// Merge implements [OrderedMap]
func (s *ThreadSafeOrderedMap[K, V]) Merge(set OrderedMap[K, V]) int {
	if s == set {
		// do not merge onto ourself!
		return 0
	}
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.Tree.Merge(set)
}

// Put implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Put(key K, value V) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Tree.Put(key, value)
}

// Remove implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) Remove(key K) (V, bool) {
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
func (s *ThreadSafeOrderedMap[K, V]) Values() iter.Seq[V] {
	return invertSeq(s, s.Tree.Values)
}

// Always returns true.
func (s *ThreadSafeOrderedMap[K, V]) ThreadSafe() bool {
	return true
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
func (s *ThreadSafeOrderedMap[K, V]) Between(a, b K, opt ...int) (total int) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	total = s.Tree.Between(a, b, opt...)
	return
}

// BetweenKV implements [OrderedMap]
func (s *ThreadSafeOrderedMap[K, V]) BetweenKV(a, b K, opt ...int) (seq iter.Seq2[K, V]) {
	return func(yield func(K, V) bool) {
		s.lock.RLock()
		defer s.lock.RUnlock()
		next, stop := iter.Pull2(s.Tree.BetweenKV(a, b, opt...))
		defer stop()
		for k, v, ok := next(); ok; k, v, ok = next() {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Takes an existing K,V iterator and returns a thread safe version.
func TsKvIterWrapper[K any, V any](seq iter.Seq2[K, V]) iter.Seq2[K, V] {
	var l sync.RWMutex
	return func(yield func(K, V) bool) {
		l.RLock()
		defer l.RUnlock()
		next, stop := iter.Pull2(seq)
		defer stop()
		for k, v, ok := next(); ok; k, v, ok = next() {
			if !yield(k, v) {
				return
			}
		}
	}
}

// RemoveBetweenKV implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) RemoveBetweenKV(a, b K, opt ...int) (seq iter.Seq2[K, V]) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return TsKvIterWrapper(s.Tree.RemoveBetweenKV(a, b, opt...))
}

// RemoveBetween implements [OrderedMap].
func (s *ThreadSafeOrderedMap[K, V]) RemoveBetween(a, b K, opt ...int) (total int) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.Tree.RemoveBetween(a, b, opt...)
}

// SetOverwrite implements [OrderedMap]
func (s *ThreadSafeOrderedMap[K, V]) SetOverwrite(cb func(key K, oldValue, newValue V)) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Tree.SetOverwrite(cb)
}

func (s *ThreadSafeOrderedMap[K, V]) SetGrowth(grow int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Tree.SetGrowth(grow)

}
