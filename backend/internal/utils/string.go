package utils

import "strings"

// Strip string of leading and trailing whitespaces
// NOTE: Return nil if s is nil
func TrimSpacePtr(s *string) *string {
	if s == nil {
		return nil
	}
	res := strings.TrimSpace(*s)
	return &res
}
