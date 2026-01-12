package omap

import "iter"

type OrderedMapExt[K any, V any] interface {
	OrderedMap[K, V]
	// Removes all elements less than or equal to the key.
	// Returns the number of elements removed.
	// The key is not required to exist.
	RemoveTo(key K) (total int)

	// Removes all elements less than or equal to the key
	// Returns a slice that contains the removed KvSet elements.
	// The key is not required to exist.
	RemoveToS(key K) (result []*KvSet[K, V])

	// Removes all elements less than or equal to the key
	// Returns a iter.Seq2 instance that contains the removed KvSet elements.
	// The key is not required to exist.
	RemoveToI(key K) iter.Seq2[K, V]

	// Removes all elements less than the key.
	// Returns the number of elements removed.
	// The key is not required to exist.
	RemoveBefore(key K) (total int)

	// Removes all elements less than the key.
	// Returns a slice that contains the removed KvSet elements.
	// The key is not required to exist.
	RemoveBeforeS(key K) (result []*KvSet[K, V])

	// Removes all elements less than the key.
	// Returns a iter.Seq2 instance that contains the removed KvSet elements.
	// The key is not required to exist.
	RemoveBeforeI(key K) iter.Seq2[K, V]

	// Removes all elements greater than or equal to the key.
	// Returns the number of elements removed.
	// The key is not required to exist.
	RemoveFrom(key K) (total int)

	// Removes all elements greater than or equal to the key.
	// Returns a slice that contains the removed KvSet elements.
	// The key is not required to exist.
	RemoveFromS(key K) (result []*KvSet[K, V])

	// Removes all elements greater than or equal to the key.
	// Returns a iter.Seq2 instance that contains the removed KvSet elements.
	// The key is not required to exist.
	RemoveFromI(key K) iter.Seq2[K, V]

	// Removes all elements less than the key.
	// Returns the number of elements removed.
	// The key does not need to exist
	RemoveAfter(key K) (total int)

	// Removes all elements less than the key.
	// Returns a slice that contains the removed KvSet elements.
	// The key is not required to exist.
	RemoveAfterS(key K) (result []*KvSet[K, V])

	// Removes all elements less than the key.
	// Returns a iter.Seq2 instance that contains the removed KvSet elements.
	// The key is not required to exist.
	RemoveAfterI(key K) iter.Seq2[K, V]
}
