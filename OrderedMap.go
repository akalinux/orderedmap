package omap

import "iter"

type OrderedMap[K any, V any] interface {
	Map[K, V]

	// If ok is true, returns the first key.
	FirstKey() (key K, ok bool)

	// If ok is true, returns the last key.
	LastKey() (K, bool)

	// Returns total number of elements between a and b.
	Between(a, b K) (total int, ok bool)

	// Returns an iterator that contains the key value sets between a and b.
	// Neither a or b are required to exist.
	BetweenKV(a, b K) (seq iter.Seq2[K, V])

	// Trys to delete the elements between a and b, returns the total number of elements deleted.
	// Neither a or b are required to exist.
	RemoveBetween(a, b K) (total int, ok bool)
}
