# OMAP a Sort Ordered map
The fastests sorted map possible for searching by ranges.
[![Go Reference](https://pkg.go.dev/badge/github.com/akalinux/orderedmap.svg)](https://pkg.go.dev/github.com/akalinux/orderedmap)
[![Go Report Card](https://goreportcard.com/badge/github.com/akalinux/orderedmap)](https://goreportcard.com/report/github.com/akalinux/orderedmap)

The omap.OrderedMap instances offer a high-performance thread-safe sorted map for Go. Optimized for O(log n) lookups and O(1) boundary inserts using pre-allocated circular slices. The drivers of the design process was the creation of a very good scheduler that could also double as a ttl cache deprecation engine.  Technically the omap package implements very minital btree using a slice. The btree implementation is ordered and does not allow for duplicates; The internals manage keys by splicing the internal slice. The side effect of this design results in what operates exactly like sorted map.

Unlike tree-based maps, omap.SliceTree and omap.CenterTree range searches use direct slice referencing, avoiding tree traversal entirely.
This is in general the optimized solution for caching, and time-series maps.

What is this optimized for?
  - appending and prepending elements to the internal array
  - searching between 2 elements
  - mass element removal between 2 elements
  - merging objects via the FastMerge method

__omap.CenterTree__
  - Implements pre-allocated a sorted slice, the elements are sequential and start at the center of the slice.
  - This offers fast search speed, along with the fastest possible range based bulk element deletion.
  - Insertion to the beginning or end of the array are done using pre-allocated memory and are o(1)

__On Write for strings__:
  - Best Case, omap.CenterTree faster than all other known orderd maps maps
  - Worst case, omap.CenterTree is half as all other known maps

## When Should you use omap.CenterTree in place of a map?

Any one of these is a practical use case:
  - when all keys must be maintained in order
  - Fuzzy logic lookups required, IE the ability to find points in between keys
  - When elements are required to be primarily manipulated by ranges

## Basic usage

The omap package provides a common interface [OrderedMap](OrderedMap.go) implemented by the following:
   - Thread safe [ThreadSafeOrderedMap](./ThreadSafeOrderedMap.go), Wrapper for any [OrderedMap](./OrderedMap.go)
   - Not thread safe [SliceTree](./SliceTree.go), but can be converted to a thread safe instance.
   - Not thread safe [CenterTree](./CenterTree.go), but can be converted to a thread safe instance.

Creating ThreadSafe instance Example:

```
	kv:=omap.NewTs[string,string](cmp.Compare)
	// Save a value
	kv.Put("Hello"," ")
	kv.Put("World","!\n")

	// Itertor
	for k,v :=range kv.All {
		fmt.Printf("%s%s",k,v)
	}
```

The resulting output will be:
```
	"Hello World!\n"
```
We can now make things a bit smaller by removing things by a range.
 ```
// Note, both "Sell" and "Zoo", were never added to the instance,
// but the between operation works on these keys any ways.
kv.RemoveBetween("Sell","Zoo")

// Itertor
for k,v :=range kv.All() {
    fmt.Printf("%s%s\n",k,v)
}
```
The resulting output will now be:

	"Hello \n"

__Why this works?__
  - The string "Sell" comes before the string "World" 
  - The string "Zoo" comes after the string "World"

__How this works__

The index lookup creates 2 values for each potential key:
  - Array position, example: 0
  - Offset can be any of the following: -1,0,1

Since lookups create both an index position and offsett, it becomes possible to look for the following:
  - Elements before the array
  - Positions between elements of the array
  - Elements after the array
  - Elements to overwrite

## API

__Constructors__

The omap package supports any key type you can provide a compare function for, but a map in go only supports a comparable key type.  This means any map can be converted to an OrderedMap instance, but not every OrderedMap instance can be converted to a map.  
| Function | Types | Arguments | Returns | Thread Safe |
|-----|------------|--------|-|-|
| omap.New | K any, V any| func(a, b K) int | *SliceTree[K, V] | false|
| omap.NewFromMap |[K comparable, V any]| map[K]V, func(a, b K) int |*SliceTree[K, V] |false |
| omap.NewSliceTree | K any, V any | int, func(a, b K) int | *SliceTree[K, V] | false |
| omap.NewTs | K any, V any | func(a, b K) int | OrderedMap[K, V] | true |
| omap.ToMap | K comparable, V any | OrderedMap[K, V]| map[K]V | false |
| omap.NewCenterTree | K any, V any | int, func(a, b K) int | *CenterTree[K, V] | false |

As a note, any instance of SliceTree or CenterTree can create a thread safe instance of itself, by calling the s.ToTs() method.  If you create a thread safe instance, you should stop using the old instance.

Example conversion from map to a thread safe OrderedMap instance:
```go
s:=omap.NewFromMap(myMap,cb).ToTs()
```

To check if an instance is thread safe call the s.ThreadSafe() method.
```go
var om OrderedMap[string,int]=New(cb)

if !om.ThreadSafe() {
	om=om.ToTs()
}
```
Why not always provide a thread safe instance?  A thread safe instance requires mutex locking, this limits what can be done even when operations are atomic.  Example: You may have a perfectly valid reason to call an iterator from within an iterator on the same instance; This cannot be done when a mutex lock is applied to an existing instance.

__OrderedMap Methods__

The following table provides a general overview of the methods in OrderedMap.

| Method | Arguments | Teturn types | Description |
|-|-|-|-|
| All | | iter.Seq2[K, V]| iterator for all Keys and Values |
| Keys | | iter.Seq[K] | iterator for all keys |
| Values | | iter.Seq[V] | iterator for all Values |
| Exists | key K | bool | true if the key was found |
| Contains | key K | bool | true if the key is between both the FirstKey and LastKey |
| Put | key K, value V | int | Sets the key and value pair |
| Get | key K | value V, ok bool | Returned the value for the key if ok is true|
| GetKvSlice |  | []KvSet[K,V] | Returns the current internal slice used for read operations |
| Filter | func(k K,v V) bool |  | deletes the given element when the callback returns true |
| FilterBetween | cb func(k K, v V) bool,a,b K, opt ...int |  | like filter, but only runs between a and b, for options [See](#between-options) |
| Remove | key K | value V, ok bool | If ok is true, the returned value was removed based on the given key |
| RemoveAll | | int | Clears all elements and returns how many elements were removed |
| MassRemove | keys ...K | int |Tries to remove all keys provided, returns how many keys were removed |
| MassRemoveKV | keys ...K | iter.Seq2[K, V] |Tries to remove all keys provided. The iterator with a copy of all key value pairs that were removed |
| Size | | int | returns the number of key/value pairs in the instance |
| FirstKey | | key K, ok bool | When ok is true the first key in the instance is returned |
| LastKey | | key K, ok bool | When ok is true the last key in the instance is returned |
| Between | a,b K, opt ...int| total int | Returns the number of elements between a and b. For options  [See](#between-options) |
| BetweenKV | a,b K, opt ...int|  iter.Seq2[K, V] | Returns an iterator that contains the key/value pairs between a and b. For options  [See](#between-options) |
| GetBetweenKvSlice | a,b K, opt ...int| []KvSet[K,V] | Returns a slice containing elements between a and b. For options  [See](#between-options) |
| RemoveBetween | a,b K, opt ...int| int | Returns the number of elements removed between a and b. For options  [See](#between-options) |
| RemoveBetweenKV | a,b K, opt ...int|  iter.Seq2[K, V] | Returns an iterator that contains the key/value pairs that were moved from between a and b. For options  [See](#between-options) |
| ThreadSafe | | bool | Returns true if this instance is thread safe |
| Merge | set OrderedMap[K, V] |int | Merges set into this instance |
| FastMerge | set OrderedMap[K, V] |int | Merges set into this instance, this requires both set and the instance be sorted in the same order |
| SetOverwrite | cb func(key K, oldValue, newValue V) | | Sets the callback method that fires before a value is overwritten |
| SetGrowth | grow int| | Sets the internal growth value for the slice |
| ToTs | | OrderedMap[K, V] | If this instance is not contained in a thread safe wrapper, returns this instance in a thread safe wrapper, other wise returns this instance |
| Clone | | OrderedMap[K, V] | Returns a new instance of self |
### Between Options

The following table exlains the usage and possible values for functions that support between operations.

| opt[id] | Options | Description |
|-|-|-|
| 0 | omap.FIRST_KEY | When set, the a field is ignored and s.FirstKey is used in its place |
| 0 | omap.LAST_KEY | When set, the b field is ignored and s.LastKey is used in its place |
| 0 | omap.FIRST_KEY\|omap.LAST_KEY | This causes both a and b to be ignored |

Example using s.BetweenKV:
```go
  for k,v :=s.BetweenKV("","Tomorrow",omap.FIRSTKEY) {
    fmt.Printf("Key: [%s], Value: [%d]\n")
  }
```
Returns all values up to "Tomorrow".

## Benchmarks

All benchmarks have been moved [here](https://github.com/akalinux/benchmarksortedmaps).