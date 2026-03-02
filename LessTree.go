package omap

func LessIndex[K any, V any](k K, Less func(a, b K) bool, Slice []KvSet[K, V]) (idx, offset int) {

	end := len(Slice)

	Cmp := func(id int) int {
		if Less(k, Slice[id].Key) {
			return -1
		} else if Less(Slice[id].Key, k) {
			return 1
		} else {
			return 0
		}
	}

	if end < 3 {
		switch end {
		case 0:
			return 0, 0
		case 1:
			return 0, Cmp(0)
		case 2:
			// being lazy for just 2 elements
			offset = Cmp(0)
			if offset < 1 {
				return 0, offset
			}
			return 1, Cmp(1)
		}
		return
	}

	mid := getMid(end)
	offset = Cmp(mid)
	idx = mid
	end--
	begin := 0
	diff := 0
	for {
		switch offset {
		case 0:
			return mid, 0
		case -1:
			// moving left
			end = mid - 1
			diff = end - begin
		case 1:
			begin = mid + 1
			diff = end - begin
		}
		if diff == 0 {
			return begin, Cmp(begin)
		} else if diff < 0 {
			return begin, -1
		}
		mid = begin + getMid(diff+1)
		offset = Cmp(mid)
	}
}
