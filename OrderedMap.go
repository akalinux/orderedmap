package omap

import "iter"

const (
	FIRST_KEY = 1
	LAST_KEY  = 2
)

type OrderedMap[K any, V any] interface {
	// Should return an iterator for all key/value pairs
	All() iter.Seq2[K, V]

	// Returns all keys, the int just expected to be a sequential number
	Keys() iter.Seq[K]

	// Returns all the values, the int is expected to be a sequential number
	Values() iter.Seq[V]

	// Returns true if the key exists, false if it does not.
	Exists(key K) bool

	// Returns true if the key is between the first and last key
	Contains(key K) bool

	// Sets the key to the value, returns the index id.
	Put(key K, value V)

	// Tries to get the value using the given key.
	// If the key exists, found is set to true, if the key does not exists then found is set to false.
	Get(key K) (value V, found bool)

	// Tries to remove the given key, returns value is set if ok is true.
	Remove(key K) (value V, ok bool)

	// Removes all elements.
	// Returns how many element were removed.
	RemoveAll() int

	// Attempts to remove all keys, returns the number of keys removed.
	MassRemove(keys ...K) (total int)

	// Attempts to remove all keys, returns an iterator with the key/value pairs
	MassRemoveKV(keys ...K) iter.Seq2[K, V]

	// Returns the current number of key/value pairs.
	Size() int

	// If ok is true, returns the first key.
	FirstKey() (key K, ok bool)

	// If ok is true, returns the last key.
	LastKey() (key K, ok bool)

	// Returns total number of elements between a and b.
	// Neither a or b are required to exist.
	//
	// The the optional opt argument:
	//  - When: opt[0]==omap.FIRST_KEY, a is ignored and the FirstKey is used.
	//  - When: opt[0]==omap.LAST_KEY, b is ignored and the LastKey is used.
	//  - To Ignore both a and b, set: opt[0]==omap.FIRST_KEY|omap.LAST_KEY.
	Between(a, b K, opt ...int) (total int)

	// Returns an iterator that contains the key value sets between a and b.
	// Neither a or b are required to exist.
	//
	// The the optional opt argument:
	//  - When: opt[0]==omap.FIRST_KEY, a is ignored and the FirstKey is used.
	//  - When: opt[0]==omap.LAST_KEY, b is ignored and the LastKey is used.
	//  - To Ignore both a and b, set: opt[0]==omap.FIRST_KEY|omap.LAST_KEY.
	BetweenKV(a, b K, opt ...int) (seq iter.Seq2[K, V])

	// Trys to delete the elements between a and b, returns the total number of elements deleted.
	// Neither a or b are required to exist.
	//
	// The the optional opt argument:
	//  - When: opt[0]==omap.FIRST_KEY, a is ignored and the FirstKey is used.
	//  - When: opt[0]==omap.LAST_KEY, b is ignored and the LastKey is used.
	//  - To Ignore both a and b, set: opt[0]==omap.FIRST_KEY|omap.LAST_KEY.
	RemoveBetween(a, b K, opt ...int) (total int)

	// Trys to delete the elements between a and b, returns an iterator that contains removed K,V pairs.
	// Neither a or b are required to exist.
	//
	// The the optional opt argument:
	RemoveBetweenKV(a, b K, opt ...int) (seq iter.Seq2[K, V])

	// Should return true if the instance is thread safe, fakse if not.
	ThreadSafe() bool

	// Merges a given [OrderedMap] into this one.
	// Returns the number of keys added
	Merge(set OrderedMap[K, V]) int

	// Sets the overwrite notice.  Its good to know when things change!
	SetOverwrite(cb func(key K, oldValue, newValue V))

	// Sets the growth value.
	SetGrowth(grow int)

	// Returns a thread safe instance.
	// If the instance is all ready thread safe, then the current instance is returned.
	ToTs() OrderedMap[K, V]

	// Deletes the given element when the callback returns true
	Filter(func(k K, v V) bool)
}
