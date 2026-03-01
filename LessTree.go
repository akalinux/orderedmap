package omap

func LessIndex[K any, V any](Less func(K) bool, Slice []KvSet[K, V]) (index int, offset byte) {

	size := len(Slice)

	if size < 3 {
		if size == 0 {
			return 0, 0
		}
		return
	}

	return
}
