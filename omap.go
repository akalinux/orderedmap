// Yet another sorted map in go.. but not really.
//
// Technically the omap package implements very minital btree using a slice.
// The drivers of the design process, were the performance objectives.
// The btree implementation is ordered and does not allow for duplicates;
// The internals manage keys by splicing the internal slice, without the use of a temporary slice.
// The side effect of this design results in what operates exactly like an ordered map.
//
// Performance objectives:
//   - Lookups for both Put and Get operations are always a fixed complexity: o(log n).
//   - All iteration operations are fixed cost of o(n).
//   - Finding or removing elements between 2 elements is always a fixed cost of o(log(n) + log(n)).
//   - Finding elements before or after a given point is always a fixed cost of o(log n)
//   - Mass Removal of unordered elements that may or may not exist has a maximum complexity of o(log(n) + log(k) + k)
//   - Key/Value pairs are kept in a struct container, the underlying slice only holds a pointer to the struct, so all splice operations operate on pointer manipulation.
//   - Pre-emptive but predictable growth, this is done by setting the Growth size.
//
// The omap package provides a common interface [OrderedMap] implemented by the following:
//   - Thread safe [ThreadSafeOrderedMap]
//   - Not thread safe [SliceTree]
//
// Basic Example:
//
//	kv:=New[string,string](cmp.Compare)
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
//	kv.RemoveBetween("Sell","Universe")
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
