# OMAP a Sorted sorted map
Yet another sorted map in go.. but not really.

Technically the omap package implements very minital btree using a slice.
The drivers of the design process, were the performance objectives.
The btree implementation is ordered and does not allow for duplicates;
The internals manage keys by splicing the internal slice.
The side effect of this design results in what operates exactly like sorted map.
Under spesific conditions or very large data sets, omap.SliceTree is faster on "Get" operations than the built in go map.
An omap.SliceTree instance uses signifigantly less the memory than the map feature in go.

# Performance Matters

Performance objectives while maintinaing a sorted map:
  - Lookups for both Put and Get operations are always a fixed complexity: o(log n).
  - All iteration operations are fixed cost of o(1).
  - Finding or removing elements between 2 points is always a fixed cost of o(log(n) + log(n)).
  - Finding elements: at, before, or after a given point is always a fixed cost of o(log n)
  - Mass Removal of unordered elements that may or may not exist has a maximum complexity of o(log(n) + log(k) + k)
  - Pre-emptive but predictable growth, this is done by setting the Growth size.

__On Read For strings__:
  - Small number of keys: omap.CenterTree and omap.SliceTree are slightly faster than go's interal map
  - Large number of keys: omap.CenterTree and omap.SliceTree is o(log n) faster than go's internal map

__On Write for strings__:
  - Best Case, omap.CenterTree twice as fast on write over go's internal map
  - Worst case, omap.CenterTree is half as fast on write as go's internal map 

## When Should you use omap.CenterTree in place of a map?

Any one of these is a practical use case:
  - You need something that can find values based on strings strings faster than go's internal map
  - A sorted map is required
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
| Remove | key K | value V, ok bool | If ok is true, the returned value was removed based on the given key |
| RemoveAll | | int | Clears all elements and returns how many elements were removed |
| MassRemove | keys ...K | int |Tries to remove all keys provided, returns how many keys were removed |
| MassRemoveKV | keys ...K | iter.Seq2[K, V] |Tries to remove all keys provided. The iterator with a copy of all key value pairs that were removed |
| Size | | int | returns the number of key/value pairs in the instance |
| FirstKey | | key K, ok bool | When ok is true the first key in the instance is returned |
| LastKey | | key K, ok bool | When ok is true the last key in the instance is returned |
| Between | a,b K, opt ...int| total int | Returns the number of elements between a and b. For options  [See](#between-options) |
| BetweenKV | a,b K, opt ...int|  iter.Seq2[K, V] | Returns an iterator that contains the key/value pairs between a and b. For options  [See](#between-options) |
| RemoveBetween | a,b K, opt ...int| int | Returns the number of elements removed between a and b. For options  [See](#between-options) |
| RemoveBetweenKV | a,b K, opt ...int|  iter.Seq2[K, V] | Returns an iterator that contains the key/value pairs that were moved from between a and b. For options  [See](#between-options) |
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
Now this is in no way a fair comparison.. The omap package can use any data set, so long as a compare function can be provided, while the map in go only needs to be optimized internally for hashing bytes, so we would expected the native map feature to faster.

__Disclamier:__ omap.SliceTree is built around a Compare function and omap.CenterTree is optimized for: appending and prepending, this means the benchmark requires creating a proxy key that is equal to the key provided by the Cmp function.  In this benchmark the key used for the map has to be generated from the base string, and the original string pointer and value then need to be saved off in an additional data structure, this gives us a like for like compare between the native go map feature and omap.SliceTree.

The following holds true for these benchmarks
  - All read/get operations are best case for go's internal map and worst case for omap.SliceTree and omap.CenterTree
  - All write operations are worst case for omap.SliceTree and best case for go's internal map and omap.CenterTree

__How well does omap compare native map feature in go?:__
```
BenchmarkNew/Native_map_Put_keys:_[1600]-10                    27043                   44257 ns/op           93058 B/op       1606 allocs/op
BenchmarkNew/Slicetree_Put_keys:_[1600]-10                     16694                   71888 ns/op           90160 B/op          4 allocs/op
BenchmarkNew/CenterTree_Put_keys:_[1600]-10                    45520                   26238 ns/op           82016 B/op          3 allocs/op
BenchmarkNew/Native_map_Get_keys:_[1600]-10                    10000                  110410 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Get_keys:_[1600]-10                     10000                  105361 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Get_keys:_[1600]-10                    10000                  106131 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Count_between:_[1600]-10                  67                17463820 ns/op           22549 B/op       2944 allocs/op
BenchmarkNew/SliceTree_Count_Between:_[1600]-10                 6213                  191290 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Count_Between:_[1600]-10                6247                  192631 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Put_keys:_[2500]-10                    15507                   77001 ns/op          169264 B/op       2510 allocs/op
BenchmarkNew/Slicetree_Put_keys:_[2500]-10                      8641                  121615 ns/op          155696 B/op          4 allocs/op
BenchmarkNew/CenterTree_Put_keys:_[2500]-10                    27486                   43528 ns/op          122976 B/op          3 allocs/op
BenchmarkNew/Native_map_Get_keys:_[2500]-10                     5215                  202513 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Get_keys:_[2500]-10                      7119                  166795 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Get_keys:_[2500]-10                     6871                  166984 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Count_between:_[2500]-10                  25                43347567 ns/op           37002 B/op       4744 allocs/op
BenchmarkNew/SliceTree_Count_Between:_[2500]-10                 4071                  293999 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Count_Between:_[2500]-10                4074                  288385 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Put_keys:_[3600]-10                     9924                  116389 ns/op          304880 B/op       3618 allocs/op
BenchmarkNew/Slicetree_Put_keys:_[3600]-10                      6532                  183078 ns/op          204848 B/op          4 allocs/op
BenchmarkNew/CenterTree_Put_keys:_[3600]-10                    20108                   59646 ns/op          180320 B/op          3 allocs/op
BenchmarkNew/Native_map_Get_keys:_[3600]-10                     2839                  372251 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Get_keys:_[3600]-10                      4876                  247516 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Get_keys:_[3600]-10                     4776                  249299 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Count_between:_[3600]-10                   9               116627112 ns/op           54710 B/op       6944 allocs/op
BenchmarkNew/SliceTree_Count_Between:_[3600]-10                 2787                  430323 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Count_Between:_[3600]-10                2760                  432141 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Put_keys:_[4900]-10                     7615                  155677 ns/op          336080 B/op       4918 allocs/op
BenchmarkNew/Slicetree_Put_keys:_[4900]-10                      3898                  305689 ns/op          417840 B/op          5 allocs/op
BenchmarkNew/CenterTree_Put_keys:_[4900]-10                    13425                   89542 ns/op          237664 B/op          3 allocs/op
BenchmarkNew/Native_map_Get_keys:_[4900]-10                      770                 1505944 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Get_keys:_[4900]-10                      3392                  353147 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Get_keys:_[4900]-10                     3356                  356358 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_map_Count_between:_[4900]-10                   4               274719299 ns/op           75738 B/op       9545 allocs/op
BenchmarkNew/SliceTree_Count_Between:_[4900]-10                 1929                  620151 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Count_Between:_[4900]-10                1921                  625068 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint64_Put_keys:_100-10           1000000000               0.0000011 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint64_Put_keys:_100-10       1000000000               0.0000046 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint64_Put_keys:_100-10        1000000000               0.0000037 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint64_Get_keys:_100-10           1000000000               0.0000011 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint64_int_Get_keys:_100-10   1000000000               0.0000030 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint64_int_Get_keys:_100-10    1000000000               0.0000028 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint64_Put_keys:_1000-10          1000000000               0.0000068 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint64_Put_keys:_1000-10      1000000000               0.0000323 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint64_Put_keys:_1000-10       1000000000               0.0000333 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint64_Get_keys:_1000-10          1000000000               0.0000053 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint64_int_Get_keys:_1000-10  1000000000               0.0000340 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint64_int_Get_keys:_1000-10   1000000000               0.0000293 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint64_Put_keys:_10000-10         1000000000               0.0000723 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint64_Put_keys:_10000-10     1000000000               0.0004179 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint64_Put_keys:_10000-10      1000000000               0.0003858 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint64_Get_keys:_10000-10         1000000000               0.0000539 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint64_int_Get_keys:_10000-10 1000000000               0.0003762 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint64_int_Get_keys:_10000-10  1000000000               0.0003694 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint8_Put_keys:_10-10             1000000000               0.0000009 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint8_int_Put_keys:_10-10     1000000000               0.0000008 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint8_int_Put_keys:_10-10      1000000000               0.0000009 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint8_Get_keys:_10-10             1000000000               0.0000007 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint8_int_Get_keys:_10-10     1000000000               0.0000008 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint8_int_Get_keys:_10-10      1000000000               0.0000007 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint8_Put_keys:_100-10            1000000000               0.0000021 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint8_int_Put_keys:_100-10    1000000000               0.0000039 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint8_int_Put_keys:_100-10     1000000000               0.0000033 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint8_Get_keys:_100-10            1000000000               0.0000022 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint8_int_Get_keys:_100-10    1000000000               0.0000034 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint8_int_Get_keys:_100-10     1000000000               0.0000022 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint8_Put_keys:_255-10            1000000000               0.0000047 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint8_int_Put_keys:_255-10    1000000000               0.0000079 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint8_int_Put_keys:_255-10     1000000000               0.0000071 ns/op               0 B/op          0 allocs/op
BenchmarkNew/Native_Map_Uint8_Get_keys:_255-10            1000000000               0.0000049 ns/op               0 B/op          0 allocs/op
BenchmarkNew/CenterTree_Map_Uint8_int_Get_keys:_255-10    1000000000               0.0000080 ns/op               0 B/op          0 allocs/op
BenchmarkNew/SliceTree_Map_Uint8_int_Get_keys:_255-10     1000000000               0.0000055 ns/op               0 B/op          0 allocs/op
```

__How to read the benchmark__

So what do these numbers really tell us?  Well nothing we didn't all ready know prior to the benchmark. The map feature of
go trades memory for read and write speed,  in particular on wirte.  Usually platforms are more cpu constrained than memory constrained, but that isn't always the case.  So we are looking at worst case reads for all omap based benhmarks and best case write for go's
internal map and omap.CenterTree.  In that aspect omap.CenterTree is a little more twice as fast best case over go's internal map for 
writes, but why?  Read to the end of this file for details.. Benchmarks are always a zero sum game!

What version of go did you run this on?
  - 1.26

What setup did you use?
  - cpu: AMD Ryzen 9 9950X 16-Core Processor
  - mem: dd5 6k mt
  - set -cpu 10
  - VM: Container under windows in docker desktop: Linux d2f8535c6a45 6.6.87.2-microsoft-standard-WSL2 #1 SMP PREEMPT_DYNAMIC Thu Jun  5 18:30:46 UTC 2025 x86_64 x86_64 x86_64 GNU/Linux

On write:
  - omap.CenterTree is faster ( when appending or prepending )
  - Most cases in smaller to medium data sets go's map is faster.. But not always.. 

On read:
  - Go map map is faster when it comes to integers on both read and write
  - omap.SliceTree and omap.CenterTree are faster when it comes to stings

Where the native map in go always perfoms worse is in memory usage:
  - Go map uses about 45%-70% more memory than omap.SliceTree or omap.CenterTree.

Scanning for any key between 2 strings:
  - Go map must be iterated over and each key must be compared o(n)
  - omap.SliceTree and omap.CenterTree are just o(log(n)+log(n))

So which is better for performance? For strings omap.SliceTree or omap.CenterTree, for int/float there is a use case for go's internal map
but its benefits don' outweigh its losses.

Scanning all keys for every key is a known worst case.. so why include it?  People do it, and most sorted map packages on [pkg.go.dev](https://pkg.go.dev) will force you do do that at least until you reach the end of your range.

__Why does the native go map read slow down so much after just 1600 strings?__

Its complicated, but its a combination of memory bandwidth and the internals the go native map doing full scans due to a large number of collisions.

__Why is SliceTree always slower on write?__

Simple: O(log n) on lookup, has to be done prior to a write.  Go's internal map skips the compare operation all until the 3rd tier.
The very thing that gives omap.SliceTree its read speed at scale slows it down in write operations.

__So which is better SliceTree CenterTree or a map?__

When it comes to read performance of strings, omap.SliceTree and omap.CenterTree are always faster.
The omap.Slicetree object is built entirly around being able to find a range of keys without scanning. 
The omap.CenterTree is optimized for pre-pending and appending data to the slice.  The map provided by go is built around hashing.
WHen it comes to a map of uint64, go's internal map has both a read and write performance advantage.

__Why include memory in benchmarks?__

This is a complex topic, but here is a short answer: Try turning memory benchmarks on for other sorted map pacakges on [pkg.go.dev](https://pkg.go.dev), they use orders of magnitued more memory than the native map in go. Most sorted map implementations arn't performacne competative with the native map in go.  The omap.SliceTree/omap.CenterTree is at least competative with the native go map implementation.  In spesific use cases omap.SliceTree and omap.CenterTree are signifigantly faster than the native map feature of go.  An instance of omap.SliceTree or omap.CenterTree do all this while being an sorted map, that is no small feat.

__Comapring go map o(1) and omap.SliceTee o(log(n))__

__Go map: o(1) How so?__  In truth go map uses the first 2 bytes of a hash as keys in a 2 tier tree, the remaing bytes then hit the o(1) or full scan.  This is the sweet spot on most use cases.  The side effect is, keys are never going to be ordered.

__omap.Slicetree is o(log n)?__ The omap.SliceTree is a btree with os type int(32|64) as its limit, indexed by sequence order. An omap.SliceTree instance is never a full scan, but an Order First Search is always more expensive in smaller sets and always cheaper in larger sets.  Effectivly omap.SliceTree is pure an Order First search with the root at the median of the array.  This hits the sweet spot for range based lookups and massive data sets.  The side effect is an ordered index of keys.

__omap.CenterTree is way faster than go's map on these benchmarks.. how so?__  The omap.CenterTree is optimized for pushing data to the
begin and end of the array.  The array is pre-allocated and the data is stored in the middle of the array.  Thus putting things before
the fist value or after the last value are super cheap. When you write to the middle go's internal map should usually be faster, but not always.

__Odd Quirks of indexing__

So which is faster a 1 byte key using a map or omap.SliceTree?
  - On Write: omap.CenterTree is faster for append and prepend
  - On Read: go's map is slightly faster than omap.SliceTree

So which is faster a 2 byte key using a map in go with 65535 elements or omap.SliceTree?
  - at 2 bytes a normal map in go is faster

So when does omap.Slicetree actually become faster? 
  - when a map in go would end up with colisions on the 3rd tier, this causes a full scan of that tier.
  - when does that happen?  Usually after a few thousand unique strings, but it depends on which buckets they land.

Is omap.SliceTree or omap.CenterTree ever faster with ints or floats?
  - Yes on range scans
  - On Writes for strings, On append or prepend omap.CenterTree is always faster
  - On Read, never on a 16 bit number
  - On Read after: the first 16 bits cause a collision and ternary logic is faster operating on the entire integer or float value

Why does omap.CenterTree have such good write perfomance?? is the 2x performance over go's internal map real?
  - Its a trap!  In reality it omap.CenterTree is optimized for append and prepend.
  - if not appending or pre-pending there is roughly a 30% write performance penalty due to the O(log n) lookup tax