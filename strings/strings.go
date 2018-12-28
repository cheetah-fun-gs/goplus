package strings

import "strings"

// StringFillLeft 左边填充 string
func StringFillLeft(in string, char string, n int) string {
	inSplit := strings.Split(in, "")
	out := []string{}
	for i := 0; i < n; i++ {
		out = append(out, char)
	}
	for _, char := range inSplit {
		out = append(out, char)
	}
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
