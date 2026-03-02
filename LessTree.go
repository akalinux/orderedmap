package omap

/*
func LessIndex[K any, V any](k K, Less func(a, b K) bool, Slice []KvSet[K, V]) (idx, offset int) {

	end := len(Slice)

	if end < 3 {
		switch end {
		case 0:
			return 0, 0
		case 1:
			if Less(k, Slice[0].Key) {
				return 0, -1
			} else if Less(Slice[0].Key, k) {
				return 0, 1
			}
			return 0, 0
		case 2:
			// being lazy for just 2 elements
			for id := range 2 {
				if Less(k, Slice[id].Key) {
					return 0, -1
				} else if !Less(Slice[id].Key, k) {
					return 0, 0
				}
			}
			return 1, 1
		}
		return
	}

	mid := getMid(end)
	idx = mid
	end--
	begin := 0
	diff := 0
	for {
		lt := Less(k, Slice[mid].Key)
		if lt {
			end = mid - 1
			diff = end - begin
		} else if Less(Slice[mid].Key, k) {
			begin = mid + 1
			diff = end - begin
		} else {
			return mid, 0
		}
		if diff == 0 {
			if Less(k, Slice[begin].Key) {
				return begin, -1
			} else if Less(Slice[begin].Key, k) {
				return begin, 1
			}
			return begin, 0
		} else if diff < 0 {
			return begin, -1
		}
		mid = begin + getMid(diff+1)
	}
}
*/
