 # OMAP OrderMap
 Yet another sorted map in go.. but not really.

 Technically the omap package implements very minital btree using a slice.
 The btree implementation is ordered and does not allow for duplicates;
 The internals manage keys by splicing the internal slice, without the use of a temporary slice.
 The side effect of this design results in what operates like ordered map.

 Performance objectives:
   - Lookups for both Put and Get operations are always a fixed complexity: o(log n).
   - All iteration operations are fixed cost of o(n).
   - Finding or removing elements between 2 elements is always a fixed cost of o(log n + log n).
   - Finding elements before or after a given point is always a fixed cost of o(log n)
   - Mass Removal of unordered elements that may or may not exist has a maximum complexity of o(log(n) + log(k) + k)

 There are 2 Implementations provided by omap, both implement [OrderedMap](./OrderedMap.go) :
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
for k,v :=range kv.All {
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

## Getting Keys and Values:

There are may ways provided by the [OrderedMap](./OrderedMap.go) interface to fetch data,
for the full source code to these examples, please look [here](./examples/Iterators/iterators.go).

__Create our instance:__
```
  s:=omap.New[int,int](cmp.Compare)
  for i := range 5 {
    s.Put(i*5,fmt.Sprintf"")
  }
```

__Get a value by key:__
```
	if value, ok := s.Get(0); ok {
		fmt.Printf("Got: [%s]\n", value)
	}
```

__Get our first key:__
```
	key, ok := s.FirstKey()
	if ok {
		fmt.Printf("First Key: %d\n", key)
	}
```

__Get our last Key:__
```
	key, ok = s.LastKey()
	if ok {
		fmt.Printf("Last Key: %d\n", key)
	}
```

__Get all keys:__
```
	for i, key := range s.Keys() {
		fmt.Printf("  ID: %d, Key: %d\n", i, key)
	}
```

__Get all values:__
```
	for i, value := range s.Values() {
		fmt.Printf("  ID: %d, String: [%s]\n", i, value)
	}
```

__Get All Keys and values:__
```
	for key, value := range s.All() {
		fmt.Printf("  Key: %d, String: [%s]\n", key, value)
	}
```

__Get all Keys and Values beteen 3 and 11:__
```
	for key, value := range s.BetweenKV(3, 11) {
		fmt.Printf("  Key: %d, String: [%s]\n", key, value)
	}
```

__Get all keys and values from the first element to the virtual key 11:__
```
	// Note, when opt is set to FIRST_KEY, the value field a is ignored
	// and the FirstKey is used.
	for key, value := range s.BetweenKV(1000, 11, omap.FIRST_KEY) {
		fmt.Printf("  Key: %d, String: [%s]\n", key, value)
	}
```

__Get all keys and values from the 5 to out last element:__
```
	// Note, when opt is set to LAST_KEY, the value field b is ignored
	// and the LastKey is used.
	for key, value := range s.BetweenKV(10, 0, omap.LAST_KEY) {
		fmt.Printf("  Key: %d, String: [%s]\n", key, value)
	}
```

