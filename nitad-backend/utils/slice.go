package utils

// this function remove the slice remove from the slice base
// ex:  base := []string{"test", "abc", "def", "ghi"}
//      remove := []string{"abc", "test"}
// return [def ghi]
// Used to remove deleteImages from ImageURLs
func RemoveSliceFromSlice(base []string, remove []string) []string {
	for i := 0; i < len(base); i++ {
		url := base[i]
		for _, rem := range remove {
			if url == rem {
				base = append(base[:i], base[i+1:]...)
				i--
				break
			}
		}
	}

	return base
}
