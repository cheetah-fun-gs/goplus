package strings

import "strings"

// StringFillLeft 左边填充 string
func StringFillLeft(in string, char string, n int) string {
	inSplit := strings.Split(in, "")
	out := []string{}
	for i := 0; i < n; i++ {
		out = append(out, char)
	}
	out = append(out, inSplit...)
	return strings.Join(out, "")
}

// StringFillRight 右边填充 string
func StringFillRight(in string, char string, n int) string {
	inSplit := strings.Split(in, "")
	for i := 0; i < n; i++ {
		inSplit = append(inSplit, char)
	}
	return strings.Join(inSplit, "")
}

// IsInStringSlice 是否在列表中
func IsInStringSlice(slice []string, find string) bool {
	for _, s := range slice {
		if s == find {
			return true
		}
	}
	return false
}
