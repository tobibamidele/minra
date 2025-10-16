// package utils
//
// import "regexp"
//
// var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)
//
// // StripANSI removes ANSI escape codes from a string
// func StripANSI(s string) string {
// 	return ansi.ReplaceAllString(s, "")
// }
//
// // VisibleWidth returns the printable width of a string (ignores ANSI).
// func VisibleWidth(s string) int {
// 	return len([]rune(StripANSI(s)))
// }
//
// // SafeSliceANSI slices a string by visible character positions, skipping escape codes.
// func SafeSliceANSI(s string, start, end int) string {
// 	if start < 0 {
// 		start = 0
// 	}
// 	if end < start {
// 		end = start
// 	}
//
// 	out := make([]rune, 0, len(s))
// 	visibleCount := 0
// 	inEscape := false
//
// 	for _, r := range s {
// 		if r == '\x1b' {
// 			inEscape = true
// 			out = append(out, r)
// 			continue
// 		}
// 		if inEscape {
// 			out = append(out, r)
// 			if (r >= 'A' && r <= 'z') {
// 				inEscape = false
// 			}
// 			continue
// 		}
//
// 		if visibleCount >= start && visibleCount < end {
// 			out = append(out, r)
// 		}
// 		visibleCount++
// 		if visibleCount >= end {
// 			break
// 		}
// 	}
//
// 	return string(out)
// }

package utils

import (
	"regexp"
	"unicode/utf8"
)

var ansi = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// StripANSI removes ANSI escape codes from a string.
func StripANSI(s string) string {
	return ansi.ReplaceAllString(s, "")
}

// VisibleWidth returns printable width of a string (ignores ANSI).
func VisibleWidth(s string) int {
	return utf8.RuneCountInString(StripANSI(s))
}

// SafeSliceANSI slices by visible rune positions, ignoring ANSI escapes.
func SafeSliceANSI(s string, start, end int) string {
	if start < 0 {
		start = 0
	}
	if end < start {
		end = start
	}

	out := make([]rune, 0, len(s))
	visibleCount := 0
	inEscape := false

	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			out = append(out, r)
			continue
		}
		if inEscape {
			out = append(out, r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				inEscape = false
			}
			continue
		}

		if visibleCount >= start && visibleCount < end {
			out = append(out, r)
		}
		visibleCount++
		if visibleCount >= end {
			break
		}
	}

	return string(out)
}

