package omap

type resolveNext[K any] struct {
	resolved bool
	begin    int
	end      int
	mid      int
	key      K
	offset   int
}
