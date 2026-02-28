// The fastests sorted map possible for caching, time-series, and scheduling.
//
// The omap.OrderedMap instances offer a high-performance thread-safe sorted map for Go.
// Optimized for O(log n) lookups and O(1) boundary inserts using pre-allocated circular slices.
// 2x faster for time-series and sequential data.
// The drivers of the design process was the creation of a very good scheduler that could also double as a ttl cache deprecation engine.
// Technically the omap package implements very minital btree using a slice.
// The btree implementation is ordered and does not allow for duplicates;
// The internals manage keys by splicing the internal slice.
// The side effect of this design results in what operates exactly like sorted map.
// Under spesific conditions or very large data sets, omap.SliceTree is faster on "Get" operations than the built in go map.
// An omap.SliceTree instance uses signifigantly less the memory than the map feature in go.
//
// Unlike tree-based maps, omap.SliceTree and omap.CenterTree range searches use direct slice referencing, avoiding tree traversal entirely.
// This is in general the optimized solution for caching, and time-series maps.
//
// # Performance Matters
//
// Performance objectives while maintinaing a sorted map:
//   - Lookups for both Put and Get operations are always a fixed complexity: o(log n).
//   - All iteration operations are fixed cost of o(1).
//   - Finding or removing elements between 2 points is always a fixed cost of o(log(n) + log(n)).
//   - Finding elements: at, before, or after a given point is always a fixed cost of o(log n)
//   - Mass Removal of unordered elements that may or may not exist has a maximum complexity of o(log(n) + log(k) + k)
//   - Pre-emptive but predictable growth, this is done by setting the Growth size.
//   - omap.SliceTree and omap.CenterTree support tunable pre-allocation
//
// The omap package provides a common interface [OrderedMap] implemented by the following:
//   - Thread safe [ThreadSafeOrderedMap], this is a wrapper for [SliceTree]
//   - Not thread safe [SliceTree], but can return a thread safe wrapper.
//   - Not thread safe [CenterTree], but can be converted to a thread safe instance.
//
// Basic Example:
//
//	kv:=omap.NewCenterTree[string,string](2,cmp.Compare)
//	// Save a value
//	kv.Put("Hello"," ")
//	kv.Put("World","!\n")
//
//	// Itertor
//	for k,v :=range kv.All {
//		fmt.Printf("%s%s",k,v)
//	}
//
// The resulting output will be:
//
//	"Hello World!\n"
//
// We can now make things a bit smaller by removing things by a range.
//
//	kv.RemoveBetween("Sell","Zoo")
//
//	// Itertor
//	for k,v :=range kv.All {
//		fmt.Printf("%s%s\n",k,v)
//	}
//
// The resulting output will now be:
//
//	"Hello \n"
//
// The index lookup creates 2 values for each potential key:
//   - Array postion, example: 0
//   - Offset can be any of the following: -1,0,1
//
// Since lookups create both an index position and offsett, it becomes possible to look for the following:
//   - Elements before the array
//   - Positions between elements of the array
//   - Elements after the array
//   - Elements to overwrite
package omap
