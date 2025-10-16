package matchers

// IsBracket checks if a rune is a bracket character.
func IsBracket(r rune) bool {
	switch r {
	case '(', ')', '[', ']', '{', '}':
		return true
	}
	return false
}

// FindMatchingBracket finds the matching bracket for a given position.
// Returns -1 if not found.
func FindMatchingBracket(line string, pos int) int {
	if pos < 0 || pos >= len([]rune(line)) {
		return -1
	}

	runes := []rune(line)
	opening := map[rune]rune{'(': ')', '[': ']', '{': '}'}
	closing := map[rune]rune{')': '(', ']': '[', '}': '{'}
	target := runes[pos]

	// Forward scan
	if match, ok := opening[target]; ok {
		count := 1
		for i := pos + 1; i < len(runes); i++ {
			if runes[i] == target {
				count++
			} else if runes[i] == match {
				count--
				if count == 0 {
					return i
				}
			}
		}
	}

	// Backward scan
	if match, ok := closing[target]; ok {
		count := 1
		for i := pos - 1; i >= 0; i-- {
			if runes[i] == target {
				count++
			} else if runes[i] == match {
				count--
				if count == 0 {
					return i
				}
			}
		}
	}

	return -1
}
