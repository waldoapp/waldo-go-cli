package lib

func Map[T any, U any](slice []T, fn func(T) U) []U {
	var result []U

	for _, item := range slice {
		result = append(result, fn(item))
	}

	return result
}
