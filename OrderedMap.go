package omap

import "iter"

const (
	FIRST_KEY = 1
	LAST_KEY  = 2
)

type OrderedMap[K any, V any] interface {
	Map[K, V]

	// If ok is true, returns the first key.
	FirstKey() (key K, ok bool)

	// If ok is true, returns the last key.
	LastKey() (K, bool)

	// Returns total number of elements between a and b.
	// Neither a or b are required to exist.
	//
	// The the optional opt argument:
	//  - When: opt==omap.FIRST_KEY, a is ignored and the FirstKey is used.
	//  - When: opt==omap.LAST_KEY, b is ignored and the LastKey is used.
	//  - To Ignore both a and b, set: opt==omap.FIRST_KEY|omap.LAST_KEY.
	Between(a, b K, opt ...int) (total int, ok bool)

	// Returns an iterator that contains the key value sets between a and b.
	// Neither a or b are required to exist.
	//
	// The the optional opt argument:
	//  - When: opt==omap.FIRST_KEY, a is ignored and the FirstKey is used.
	//  - When: opt==omap.LAST_KEY, b is ignored and the LastKey is used.
	//  - To Ignore both a and b, set: opt==omap.FIRST_KEY|omap.LAST_KEY.
	BetweenKV(a, b K, opt ...int) (seq iter.Seq2[K, V])

	// Trys to delete the elements between a and b, returns the total number of elements deleted.
	// Neither a or b are required to exist.
	//
	// The the optional opt argument:
	//  - When: opt==omap.FIRST_KEY, a is ignored and the FirstKey is used.
	//  - When: opt==omap.LAST_KEY, b is ignored and the LastKey is used.
	//  - To Ignore both a and b, set: opt==omap.FIRST_KEY|omap.LAST_KEY.
	RemoveBetween(a, b K, opt ...int) (total int, ok bool)

	// Trys to delete the elements between a and b, returns an iterator that contains removed K,V pairs.
	// Neither a or b are required to exist.
	//
	// The the optional opt argument:
	//  - When: opt==omap.FIRST_KEY, a is ignored and the FirstKey is used.
	//  - When: opt==omap.LAST_KEY, b is ignored and the LastKey is used.
	//  - To Ignore both a and b, set: opt==omap.FIRST_KEY|omap.LAST_KEY.
	RemoveBetweenKV(a, b K, opt ...int) (seq iter.Seq2[K, V])
}
