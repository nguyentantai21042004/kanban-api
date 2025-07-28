package util

func ToPointer[T any](v T) *T {
	return &v
}

func DerefSlice[T any](ptrs []*T) []T {
	result := make([]T, 0, len(ptrs))
	for _, p := range ptrs {
		if p != nil {
			result = append(result, *p)
		}
	}
	return result
}
