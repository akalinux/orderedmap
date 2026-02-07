# OMAP an Ordered Map
Yet another sorted map in go.. but not really.

Technically the omap package implements very minital btree using a slice.
The drivers of the design process, were the performance objectives.
The btree implementation is ordered and does not allow for duplicates;
The internals manage keys by splicing the internal slice.
The side effect of this design results in what operates exactly like ordered map.
Under spesific conditions or very large data sets, omap.SliceTree is faster on "Get" operations than the built in go map.
An omap.SliceTree instance uses signifigantly less the memory than the map feature in go.

Performance objectives while maintinaing a sorted map:
  - Lookups for both Put and Get operations are always a fixed complexity: o(log n).
  - All iteration operations are fixed cost of o(1).
  - Finding or removing elements between 2 points is always a fixed cost of o(log(n) + log(n)).
  - Finding elements: at, before, or after a given point is always a fixed cost of o(log n)
  - Mass Removal of unordered elements that may or may not exist has a maximum complexity of o(log(n) + log(k) + k)
  - Pre-emptive but predictable growth, this is done by setting the Growth size.

## When Should you use omap.SliceTree in place of a map?

Any one of these is a practical use case:
  - An ordered map is required
  - Memory constrained systems
  - Fuzzy logic is required, IE the ability to find points in between keys
  - When a combination of freequent updates and searching by ranges is requried
  - Very large data sets where read speed is more important than write speed
  - Keys that can not be represented as a comparable value
  - When managing elements between ranges is required

## Basic usage

The omap package provides a common interface [OrderedMap](OrderedMap.go) implemented by the following:
   - Thread safe [ThreadSafeOrderedMap](./ThreadSafeOrderedMap.go), Wrapper for [SliceTree](./SliceTree.go)
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
| MassRemoveKV | keys ...K | iter.Seq2[K, V] |Tries to remove all keys provided. The iterator with a copy of all key value pairs that were removed |
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

__Disclamier:__ omap.SliceTree is built around a Compare function, this means the benchmark requires creating a proxy key that is equal to the key provided by the Cmp function.  In this benchmark the key used for the map has to be generated from the base string, and the original string pointer and value then need to be saved off in an additional data structure, this gives us a like for like compare between the native go map feature and omap.SliceTree.

__How well does omap compare native map feature in go?:__
```
BenchmarkNew/Native_map_Put,_keys:_[1600]-10                               26467             44379 ns/op           93056 B/op       1606 allocs/op
BenchmarkNew/Slicetree,_Put,_keys:_[1600]-10                                9781            114559 ns/op           41008 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get,_keys:_[1600]-10                              10000            109108 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get,_keys:_[1600]-10                                7849            150049 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_Count_between,_keys:_[1600]-10                       67          17179143 ns/op           22552 B/op       2944 allocs/op
BenchmarkNew/SliceTree,_Count_Between_nodes_keys:_[1600]-10                 4364            267578 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Put,_keys:_[2500]-10                               15544             77330 ns/op          169264 B/op       2510 allocs/op
BenchmarkNew/Slicetree,_Put,_keys:_[2500]-10                                5958            186216 ns/op           65584 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get,_keys:_[2500]-10                               5620            200192 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get,_keys:_[2500]-10                                5156            231305 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_Count_between,_keys:_[2500]-10                       27          41700145 ns/op           36988 B/op       4744 allocs/op
BenchmarkNew/SliceTree,_Count_Between_nodes_keys:_[2500]-10                 2847            425970 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Put,_keys:_[3600]-10                               10038            119693 ns/op          304880 B/op       3618 allocs/op
BenchmarkNew/Slicetree,_Put,_keys:_[3600]-10                                3668            302038 ns/op           90160 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get,_keys:_[3600]-10                               2881            373605 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get,_keys:_[3600]-10                                3358            351491 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_Count_between,_keys:_[3600]-10                        9         114708859 ns/op           54710 B/op       6944 allocs/op
BenchmarkNew/SliceTree,_Count_Between_nodes_keys:_[3600]-10                 1914            627752 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Put,_keys:_[4900]-10                                7454            159133 ns/op          336080 B/op       4918 allocs/op
BenchmarkNew/Slicetree,_Put,_keys:_[4900]-10                                2288            480656 ns/op          122928 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get,_keys:_[4900]-10                                764           1435220 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get,_keys:_[4900]-10                                2323            499635 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_Count_between,_keys:_[4900]-10                        4         270712341 ns/op           75738 B/op       9545 allocs/op
BenchmarkNew/SliceTree,_Count_Between_nodes_keys:_[4900]-10                 1357            881794 ns/op               0 B/op          0 allocs/op
```

__How to read the benchmark__

So what do these numbers really tell us?  Well nothing we didn't all ready know prior to the benchmark. The map feature of
go trades memory for read and write speed,  in particular on wirte.  Usually platforms are more cpu constrained than memory constrained, but that isn't always the case.

On write:
  - Go map is faster 

On read:
  - Go map is faster with smaller sets of keys
  - omap.SliceTree is faster on larger sets of keys

Where the native map in go always perfoms worse is in memory usage:
  - Go map uses about 45%-70% more memory than omap.SliceTree.

Scanning for any key between 2 strings:
  - Go map must be iterated over and each key must be compared o(n)
  - omap.SliceTree is just o(log(n)+log(n))

So which is better for performance?  Depends.. If you need very large sets of keys and you would have to create a proxy key to represent the raw key, then use omap.SliceTree, other wise use map.  

Scanning all keys for every key is a known worst case.. so why include it?  People do it, and most ordered map packages on [pkg.go.dev](https://pkg.go.dev) will force you do do that at least until you reach the end of your range.

__Why does the native go map read slow down so much after 2500 elements?__

Its complicated, but its a combination of memory bandwidth and the internals the go native map doing full scans due to a large number of collisions.

__Why is SliceTree always slower on write?__

Simple: memory bandwidth.  omap.SliceTree uses a slice in memory and the elemetns are copied back and forth in blocks.
Although the lookup and storage uses less memory, the write operation involes splicing an array, which will go to main memory 
more often.  The very thing that gives omap.SliceTree its read speed at scale slows it down in write operations.

__So which is better SliceTree or a map?__

Very subjective, omap.Slicetree is built entirly around being able to find a range of keys without scanning.  A map in go is
built around hashing bytes.  After a certan point omap.Slicetree will always read faster, but 
the go native map feature will aways write faster.  The go map implementation will always use more memory than omap.SliceTree, but the native go map perfoms sligtly less than half the memory modification of omap.SliceTree on insertion/update.

__Why include memory in benchmarks?__

This is a complex topic, but here is a short answer: Try turning memory benchmarks on for other ordered map pacakges on [pkg.go.dev](https://pkg.go.dev), they use orders of magnitued more memory than the native map in go. Most ordered map implementations arn't performacne competative with the native map in go.  The omap.SliceTree is at least competative with the native go map implementation.  In spesific use cases omap.SliceTree is signifigantly faster than the native map feature of go.  An omap.SliceTree instance does all this while being an ordered map, that is no small feat.

__Comapring go map o(1) and omap.SliceTee o(log(n))__

__Go map: o(1) How so?__  In truth go map uses the first 2 bytes as keys in a 2 tier tree, the remaing bytes then hit the o(1) or full scan.  This is the sweet spot on most use cases.  The side effect is, keys are never going to be ordered.

__omap.Slicetree is o(log n)?__ The omap.SliceTree is a btree with os type int(32|64) as its limit, indexed by sequence order. An omap.SliceTree instance is never a full scan, but an Order First Search is always more expensive in smaller sets and always cheaper in larger sets.  Effectivly omap.SliceTree is pure an Order First search with the root at the median of the array.  This hits the sweet spot for range based lookups and massive data sets.  The side effect is an ordered index of keys.

__Odd Quirks of indexing__

So which is faster a 1 byte key using a map or omap.SliceTree?
  - at a one byte key omap.SliceTree is faster

So which is faster a 2 byte key using a map in go with 65535 elements or omap.SliceTree?
  - at 2 bytes a normal map in go is faster

So when does omap.Slicetree actually become faster? 
  - when a map in go would end up with colisions on the 3rd tier, this causes a full scan of that tier.
  - when does that happen?  Usually after a few thousand unique strings, but it depends on which buckets they land.

Is omap.SliceTree ever faster with ints or floats?
  - Yes on range scans
  - On Writes, Never
  - On Read, never on a16 bit number
  - On Read after: the first 16 bits cause a collision and ternary logic is faster operating on the entire integer or float value
