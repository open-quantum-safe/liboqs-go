package oqstests

import "strings"

// stringMatchSlice returns true if str contains as a substring some element of
// the string slice s, and false otherwise. For example, the function returns
// true for the case str = "test" and s = ["testing", "element"], and false for
// the case str = "test" and s = ["happy", "dog"].
func stringMatchSlice(str string, s []string) bool {
	for _, pattern := range s {
		if strings.Contains(str, pattern) {
			return true
		}
	}
	return false
}
