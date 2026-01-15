# OMAP OrderMap
Yet another sorted map in go.. but not really.

Technically the omap package implements very minital btree using a slice.
The drivers of the design process, were the performance objectives.
The btree implementation is ordered and does not allow for duplicates;
The internals manage keys by splicing the internal slice, without the use of a temporary slice.
The side effect of this design results in what operates exactly like ordered map.

Performance objectives while maintinaing a sorted map:
  - Lookups for both Put and Get operations are always a fixed complexity: o(log n).
  - All iteration operations are fixed cost of o(1).
  - Finding or removing elements between 2 points is always a fixed cost of o(log(n) + log(n)).
  - Finding elements before or after a given point is always a fixed cost of o(log n)
  - Mass Removal of unordered elements that may or may not exist has a maximum complexity of o(log(n) + log(k) + k)
  - Pre-emptive but predictable growth, this is done by setting the Growth size.

## When Should you use omap?

Any one of these is a Practical use case:
  - An ordered map is required
  - Fuzzy logic is required, IE the ability to find points in between keys
  - When a combination of freequent updates and searching by ranges is requried
  - Memory is more important than cpu time for your maps

## Basic usage

The omap package provides a common interface [OrderedMap](OrderedMap.go) implemented by the following:
   - Thread safe [ThreadSafeOrderedMap](./ThreadSafeOrderedMap.go)
   - Not thread safe [SliceTree](./SliceTree.go), but can be converted to a thread safe instance.


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

__Why this works?__
  - The string "Sell" comes before the string "World" 
  - The string "Universe" comes after the string "World"

__How this works__

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
| All | | iter.Seq2[K, V]| iterator for all Keys and Values |
| Keys | | iter.Seq[K] | iterator for all keys |
| Values | | iter.Seq[V] | iterator for all Values |
| Exists | key K | bool | true if the key was found |
| Contains | key K | bool | true if the key is between both the FirstKey and LastKey |
| Put | key K, value V | int | Sets the key and value pair |
| Get | key K | value V, ok bool | Returned the value for the key if ok is true|
| Remove | key K | value V, ok bool | If ok is true, the returned value was removed based on the given key |
| RemoveAll | | int | Clears all elements and returns how many elements were removed |
| MassRemove | keys ...K | int |Tries to remove all keys provided, returns how many keys were removed |
| MassRemoveKV | keys ...K | iter.Seq2[int, K] |Tries to remove all keys provided. The iterator with a copy of all key value pairs that were removed |
| Size | | int | returns the number of key/value pairs in the instance |
| FirstKey | | key K, ok bool | When ok is true the first key in the instance is returned |
| LastKey | | key K, ok bool | When ok is true the last key in the instance is returned |
| Between | a,b K, opt ...int| total int | Returns the number of elements between a and b. For options  [See](#between-options) |
| BetweenKV | a,b K, opt ...int|  iter.Seq2[int, K] | Returns an iterator that contains the key/value pairs between a and b. For options  [See](#between-options) |
| RemoveBetween | a,b K, opt ...int| int | Returns the number of elements removed between a and b. For options  [See](#between-options) |
| RemoveBetweenKV | a,b K, opt ...int|  iter.Seq2[int, K] | Returns an iterator that contains the key/value pairs that were moved from between a and b. For options  [See](#between-options) |
| ThreadSafe | | bool | Returns true if this instance is thread safe |
| Merge | set OrderedMap[K, V] |int | Merges set into this instance |
| SetOverwrite | cb func(key K, oldValue, newValue V) | | Sets the callback method that fires before a value is overwritten |
| SetGrowth | grow int| | Sets the internal growth value for the slice |
| ToTs() | | OrderedMap[K, V] | If this instance is not contained in a thread safe wrapper, returns this instance in a thread safe wrapper, other wise returns this instance |
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

So benchmarks are always very subjective, but the real question is: what do we compare omap too?  The only real answer is the native map in go.
Now this is in no way a fair comparison.. The omap package can use any data set, so long as a compare function can be provided, while the map in go only needs to be optimized internally for hashing bytes, so we would expected the native map feature to be an order of magnitude faster.

__How well does omap compare native map feature in go?:__
```
BenchmarkNew/Native_map,_size:_[4],_keys:_[100]-10                623304              1788 ns/op            5464 B/op          4 allocs/op
BenchmarkNew/Slicetree,_size_[4],_count_[100]-10                  298898              3862 ns/op            2736 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get_size:_[4],_keys:_[100]-10            755082              1581 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get_size:_[4],_keys:_[100]-10             331154              3565 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_size:_[4],_keys:_[400]-10                151936              7871 ns/op           21848 B/op          4 allocs/op
BenchmarkNew/Slicetree,_size_[4],_count_[400]-10                   57602             20789 ns/op            9776 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get_size:_[4],_keys:_[400]-10             89725             13153 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get_size:_[4],_keys:_[400]-10              37513             31975 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_size:_[4],_keys:_[900]-10                 53892             22983 ns/op           98432 B/op          6 allocs/op
BenchmarkNew/Slicetree,_size_[4],_count_[900]-10                   21109             56476 ns/op           21808 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get_size:_[4],_keys:_[900]-10             31167             38754 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get_size:_[4],_keys:_[900]-10              15006             79956 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_size:_[4],_keys:_[1600]-10                43485             27773 ns/op           98432 B/op          6 allocs/op
BenchmarkNew/Slicetree,_size_[4],_count_[1600]-10                   9874            117807 ns/op           41008 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get_size:_[4],_keys:_[1600]-10            10000            108882 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get_size:_[4],_keys:_[1600]-10              7345            147573 ns/op               0 B/op          0 allocs/op
```

So what do these numbers really tell us?  Well nothing we didn't all ready know prior to the benchmark. The map feature of
go trades memory for read and write speed,  in particular on wirte.  Usually platforms are more cpu constrained than memory constrained, but that isn't always the case.

On write:
  - Go map map is on average 40%-70% faster than omap.SliceTree.

On read:
  - Go map on average is about 27%-55% faster Than omap.SliceTree

Where the native map in go always perfoms worse is in memory usage:
  - Go map uses about 2-4 times more memory than omap.SliceTree.
