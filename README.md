# OMAP OrderMap
Yet another sorted map in go.. but not really.

Technically the omap package implements very minital btree using a slice.
The drivers of the design process, were the performance objectives.
The btree implementation is ordered and does not allow for duplicates;
The internals manage keys by splicing the internal slice, without the use of a temporary slice.
The side effect of this design results in what operates exactly like ordered map.

Performance objectives:
  - Lookups for both Put and Get operations are always a fixed complexity: o(log n).
  - All iteration operations are fixed cost of o(n).
  - Finding or removing elements between 2 points is always a fixed cost of o(log(n) + log(n)).
  - Finding elements before or after a given point is always a fixed cost of o(log n)
  - Mass Removal of unordered elements that may or may not exist has a maximum complexity of o(log(n) + log(k) + k)
  - Pre-emptive but predictable growth, this is done by setting the Growth size.

The omap package provides a common interface [OrderedMap](OrderedMap.go) implemented by the following:
   - Thread safe [ThreadSafeOrderedMap](./ThreadSafeOrderedMap.go)
   - Not thread safe [SliceTree](./SliceTree.go)

## Basic usage

Creating ThreadSafe instance Example:

```
	kv:=NewTs[string,string](cmp.Compare)
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
// Note, both "Sell" and "Universe", were never added to the instance,
// but the between operation works on these keys any ways.
kv.RemoveBetween("Sell","Universe")

// Itertor
for k,v :=range kv.All() {
    fmt.Printf("%s%s\n",k,v)
}
```
The resulting output will now be:

	"Hello \n"

### Why this works?
  - The string "Sell" comes before the string "World" 
  - The string "Universe" comes after the string "World"

### How this works

The index lookup creates 2 values for each potential key:
  - Array postion, example: 0
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

As a note, any instance of SliceTree can create a thread safe instance of itself, by calling the s.ToTs() method.  If you create a thread safe instance, you should stop using the old instance.

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
| All | | iter.Seq2[K, V]| iterator for all Veys and Values |
| Keys | | iter.Seq2[int, K] | iterator for all keys |
| Values | | iter.Seq2[int, K] | iterator for all Values |
| Exists | key K | bool | true if the key was found |
| Put | key K, value V | int | Sets the key and value pair, and returns the index id |
| Get | key K | value V, ok bool | Returned the value for the key if ok is true|
| Remove | key K | value V, ok bool | If ok is true, the returned value was removed based on the given key |
| RemoveAll | | int | Clears all elements and returns how many elements were removed |
| MassRemove | keys ...K | int |Tries to remove all keys provided, returns how many keys were removed |
| MassRemoveKV | keys ...K | iter.Seq2[int, K] |Tries to remove all keys provided, returns an iterator with a copy of all key value pairs that were removed |
| Size | | int | returns the number of key/value pairs in the instance |
| FirstKey | | key K, ok bool | When ok is true the first key in the instance is returned |
| LastKey | | key K, ok bool | When ok is true the last key in the instance is returned |
| Between | a,b K, opt ...int| total int | Returns the number of elements between a and b. For options  [See](#between-options) |

### Between Options