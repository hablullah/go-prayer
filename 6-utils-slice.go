package prayer

func sliceRealIdx[T any](arr []T, idx int) int {
	arrLen := len(arr)

	if idx < 0 {
		for {
			idx += arrLen
			if idx >= 0 {
				return idx
			}
		}
	}

	if idx >= arrLen {
		for {
			idx -= arrLen
			if idx < arrLen {
				return idx
			}
		}
	}

	return idx
}

func sliceAt[T any](arr []T, idx int) T {
	realIdx := sliceRealIdx[T](arr, idx)
	return arr[realIdx]
}

func firstSliceItem[T any](arr []T) (T, bool) {
	var zero T
	if len(arr) == 0 {
		return zero, false
	}
	return arr[0], true
}

func lastSliceItem[T any](arr []T) (T, bool) {
	var zero T
	if len(arr) == 0 {
		return zero, false
	}
	return arr[len(arr)-1], true
}
