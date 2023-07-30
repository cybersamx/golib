package stringsutils

import (
	"fmt"
)

func RuneToString(r rune) string {
	if r == rune(0) {
		return ""
	}

	return string(r)
}

// EllipsisString returns a string truncated and appended with an ellipsis ... if the text goes over
// the limit.
func EllipsisString(text string, limit int) string {
	const ellipsisLen = 3
	if limit <= ellipsisLen {
		return "..."
	}

	if len(text) <= limit {
		return text
	}

	return fmt.Sprintf("%s...", text[:limit-3])
}

// TruncateString returns a string truncated at the specified limit.
func TruncateString(text string, limit int) string {
	if len(text) < limit {
		return text
	}

	return text[:limit]
}
