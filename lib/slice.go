package lib

func CompactMap[T any, U any](slice []T, fn func(T) (U, bool)) []U {
	var result []U

	for _, item := range slice {
		if newItem, keep := fn(item); keep {
			result = append(result, newItem)
		}
	}

	return result
}

func Map[T any, U any](slice []T, fn func(T) U) []U {
	var result []U

	for _, item := range slice {
		result = append(result, fn(item))
	}

	return result
}
