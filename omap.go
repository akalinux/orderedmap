// The fastests sorted map possible for searching by ranges.
//
// Example:
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
//   - Array position, example: 0
//   - Offset can be any of the following: -1,0,1
//
// Since lookups create both an index position and offsett, it becomes possible to look for the following:
//   - Elements before the array
//   - Positions between elements of the array
//   - Elements after the array
//   - Elements to overwrite
package omap
