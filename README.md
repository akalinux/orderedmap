# OMAP an Ordered Map
Yet another sorted map in go.. but not really.

Technically the omap package implements very minital btree using a slice.
The drivers of the design process, were the performance objectives.
The btree implementation is ordered and does not allow for duplicates;
The internals manage keys by splicing the internal slice.
The side effect of this design results in what operates exactly like ordered map.
Under spesific conditions or very large data sets, omap.SliceTree is faster on "Get" operations than the built in go map.
An omap.SliceTree instance uses signifigantly less the memory than the map feature ing go.

Performance objectives while maintinaing a sorted map:
  - Lookups for both Put and Get operations are always a fixed complexity: o(log n).
  - All iteration operations are fixed cost of o(1).
  - Finding or removing elements between 2 points is always a fixed cost of o(log(n) + log(n)).
  - Finding elements before or after a given point is always a fixed cost of o(log n)
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

__Disclamier:__ omap.SliceTree is built around a Compare function, this means the benchmark requires creating a proxy key that is equal to the key provided by the Cmp function.  In this benchmark the key used for the map has to be generated from the base string, and the original string pointer and value then need to be saved off in an additional data structure, this gives us a like for like compare between the native go map feature set and omap.SliceTree.

__How well does omap compare native map feature in go?:__
```
BenchmarkNew/Native_map,_size:_[4],_keys:_[1600]-10                26866             44408 ns/op           93056 B/op       1606 allocs/op
BenchmarkNew/SliceTree,_size_[4],_count_[1600]-10                   9758            115444 ns/op           41008 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get_size:_[4],_keys:_[1600]-10            10000            107863 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get_size:_[4],_keys:_[1600]-10              8061            152664 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_size:_[4],_keys:_[2500]-10                15519             77506 ns/op          169264 B/op       2510 allocs/op
BenchmarkNew/SliceTree,_size_[4],_count_[2500]-10                   6106            185458 ns/op           65584 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get_size:_[4],_keys:_[2500]-10             3975            280004 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get_size:_[4],_keys:_[2500]-10              5030            228491 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_size:_[4],_keys:_[3600]-10                10060            118377 ns/op          304880 B/op       3618 allocs/op
BenchmarkNew/SliceTree,_size_[4],_count_[3600]-10                   3798            296544 ns/op           90160 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get_size:_[4],_keys:_[3600]-10             2365            470665 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get_size:_[4],_keys:_[3600]-10              3288            347730 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map,_size:_[4],_keys:_[4900]-10                 7450            155436 ns/op          336080 B/op       4918 allocs/op
BenchmarkNew/SliceTree,_size_[4],_count_[4900]-10                   2485            467271 ns/op          122928 B/op          2 allocs/op
BenchmarkNew/Native_map,_Get_size:_[4],_keys:_[4900]-10              840           1379415 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree,_Get_size:_[4],_keys:_[4900]-10              2437            495088 ns/op               0 B/op          0 allocs/op
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

So which is better for performance?  Depends.. If you need very large sets of keys and you would have to create a proxy key to represent the raw key, then use omap.SliceTree, other wise use map.

__Why does the native map read slow down so much after 2500 elements?__

Simple, memory bandwidth. SliceTree for the same kind of work, uses between 45%-70% less memory than a native map in go, also the inner workings of SliceTree is a single slice in memory.  This means memory reads for SliceTree are usually contiguous. Although the native go map hasing operation is cheaper when it comes to cpu cycles. omap.SliceTree is built entierly around a single slice and can often times stay entierly inside the cpu cache on reads.  Even when going to main memory, a slice is usually a sequential set of reads, which greatly improve memory read performance.

__Why is SliceTree always slower on write?__

Simple, memory bandwidth.  SliceTree uses a slice in memory and the elemetns are copied back and forth in blocks.
Although the lookup and storage uses less memory, the write operation involes splicing an array, which will go to main memory 
more often.  The very thing that gives SliceTree its read speed at scale slows it down in write operations.

__So which is better SliceTree or a map?__

Very subjective.  Slicetree is built entirly around being able to find a range without scanning.  A map in go is
built around quick reads and writes with very simple keys.  After a certan point Slicetree will always read faster, but 
the native map feature will aways write faster.