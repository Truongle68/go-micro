package sliceutil

func DedupeString[T any](items []T, keyFn func(T) string) []string {
	seen := make(map[string]bool, len(items))
	results := make([]string, 0, len(items))

	for _, item := range items {
		val := keyFn(item)
		if val == "" || seen[val] {
			continue
		}
		seen[val] = true
		results = append(results, val)
	}
	return results
}
