package string

import "strings"

// ContainsArray returns true if str contains all items in c
func ContainsArray(str string, c []string) bool {
	for _, s := range c {
		if !strings.Contains(str, s) {
			return false
		}
	}
	return true
}

// GetStringInBetween Returns ("", true)if no start string found. It returns
// non-empty string and false if it can find some useful string.
func GetStringInBetween(str string, start string, end string) (result string, empty bool) {
	s := strings.Index(str, start)
	if s == -1 {
		return result, true
	}
	newS := str[s+len(start):]
	e := strings.Index(newS, end)
	if e == -1 {
		return result, true
	}
	result = newS[:e]
	return result, false
}
