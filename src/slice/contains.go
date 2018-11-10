// Package slice provides ...
package slice

// Contains returns wether a slice contains a specific element
func Contains(slice []string, element string) bool {
	for _, sliceElement := range slice {
		if element == sliceElement {
			return true
		}
	}

	return false
}
