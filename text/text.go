package text

import "strings"

// Ellipsis a text
func Ellipsis(length int, text string) string {
	r := []rune(text)
	if len(r) > length {
		return string(r[0:length]) + "..."
	}
	return text
}

func Dos2Unix(s string) string {
	return strings.ReplaceAll(s, "\r", "")
}
