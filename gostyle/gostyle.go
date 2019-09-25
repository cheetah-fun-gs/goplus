package gostyle

import (
	"fmt"
	"strings"
)

// golint 的变量名检查
var commonInitialisms = map[string]bool{
	// golint
	"ACL":   true,
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XMPP":  true,
	"XSRF":  true,
	"XSS":   true,
}

func firstToUpper(s string) string {
	split := strings.Split(s, "")
	dst := []string{strings.ToUpper(split[0])}
	dst = append(dst, split[1:]...)
	return strings.Join(dst, "")
}

// FormatToCamelCase 键名转换为驼峰命名法
func FormatToCamelCase(s string) string {
	var split []string
	for _, word := range strings.Split(s, "_") {
		if _, ok := commonInitialisms[strings.ToUpper(word)]; ok {
			split = append(split, strings.ToUpper(word))
			continue
		}
		split = append(split, firstToUpper(word))
	}
	return strings.Join(split, "")
}

// FormatToGoPackage 去除下划线 获取包名或文件名
func FormatToGoPackage(srcStr string) (dstStr string) {
	dstStr = strings.Replace(srcStr, "_", "", -1)
	return strings.ToLower(dstStr)
}

// Tag tag
type Tag struct {
	Key   string
	Value string
}

// TagJoin 拼接tag
func TagJoin(tags []*Tag) string {
	split := []string{}
	for _, tag := range tags {
		split = append(split, fmt.Sprintf(`%s:"%s"`, tag.Key, tag.Value))
	}
	return strings.Join(split, " ")
}
