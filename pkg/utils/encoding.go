package utils

import (
	"os"
	"unicode/utf8"
)

// DetectFileEncoding returns the codec for the file, utf-8,latin-1
func DetectFileEncoding(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	// 1: Check BOM
	// 1. Check BOM
	if len(data) >= 3 && data[0] == 0xEF && data[1] == 0xBB && data[2] == 0xBF {
		return "utf-8", nil // UTF-8 BOM
	}
	if len(data) >= 2 && data[0] == 0xFF && data[1] == 0xFE {
		return "utf-16le", nil
	}
	if len(data) >= 2 && data[0] == 0xFE && data[1] == 0xFF {
		return "utf-16be", nil
	}

	// 2. Check UTF-8 validity
	if utf8.Valid(data) {
		return "utf-8", nil
	}

	// 3. Fallback to Latin-1 (ISO-8859-1) check
	latin1 := true
	for _, b := range data {
		if b > 0xFF { // Latin-1 is only single byte
			latin1 = false
			break
		}
	}
	if latin1 {
		return "latin-1", nil
	}

	// 4. Unknown
	return "unknown", nil
}
