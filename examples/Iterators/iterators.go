package main

import (
	"cmp"
	"fmt"

	omap "github.com/akalinux/orderedmap"
)

func buildMap() omap.OrderedMap[int, string] {
	s := omap.New[int, string](cmp.Compare)

	// Create our data set
	// Keys: 0,5,10,15,20
	for i := range 5 {
		s.Put(i*5, fmt.Sprintf("Value: %d", i))
	}

	return s
}
func main() {
	// Create our map
	s := buildMap()

	// Get our first element
	key, ok := s.FirstKey()
	if ok {
		fmt.Printf("First Key: %d\n", key)
	}

	// Get our last element
	key, ok = s.LastKey()
	if ok {
		fmt.Printf("Last Key: %d\n", key)
	}

	fmt.Printf("\n#All Keys\n")
	// Get all keys
	for i, key := range s.Keys() {
		fmt.Printf("  ID: %d, Key: %d\n", i, key)
	}

	// Get all values
	fmt.Printf("\n#All Values\n")
	for i, value := range s.Values() {
		fmt.Printf("  ID: %d, String: [%s]\n", i, value)
	}

	// All Keys and values
	fmt.Printf("\n#All Keys and Values\n")
	for key, value := range s.All() {
		fmt.Printf("  Key: %d, String: [%s]\n", key, value)
	}

	// Get all Keys and Values beteen 3 and 11
	fmt.Printf("\n#Keys and Values, between a range\n")
	for key, value := range s.BetweenKV(3, 11) {
		fmt.Printf("  Key: %d, String: [%s]\n", key, value)
	}

	// Get all keys and values from the first element to 11
	fmt.Printf("\n#Keys and Values, from start to point\n")
	// Note, when the opt is set to FIRST_KEY, the value field a is ignored
	// and the FirstKey is used.
	for key, value := range s.BetweenKV(1000, 11, omap.FIRST_KEY) {
		fmt.Printf("  Key: %d, String: [%s]\n", key, value)
	}

	// Get all keys and values from the 5 to out last element
	fmt.Printf("\n#Keys and Values, from start to point\n")
	// Note, when the opt is set to LAST_KEY, the value field b is ignored
	// and the LastKey is used.
	for key, value := range s.BetweenKV(10, 0, omap.LAST_KEY) {
		fmt.Printf("  Key: %d, String: [%s]\n", key, value)
	}
}
