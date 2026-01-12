package main

import (
	"cmp"
	"fmt"

	omap "github.com/akalinux/orderedmap"
)

func main() {
	kv := omap.NewTs[string, string](cmp.Compare)

	kv.Put("Hello", " ")
	kv.Put("World", "!\n")
	for k, v := range kv.All() {
		fmt.Printf("%s%s", k, v)
	}
	kv.RemoveBetween("Sell", "Zoo")

	for k, v := range kv.All() {
		fmt.Printf("%s%s\n", k, v)
	}
}
