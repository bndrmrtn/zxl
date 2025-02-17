package lexer

// isLetter returns true if the given character is a letter
func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '$'
}

// isDigit returns true if the given character is a digit
func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
