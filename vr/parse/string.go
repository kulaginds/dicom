package parse

import (
	"strings"
	"unicode"
)

func UI(data []byte) string {
	str := string(data)

	onlySpaces := true
	for _, char := range str {
		if !unicode.IsSpace(char) {
			onlySpaces = false
		}
	}
	if !onlySpaces {
		// String may have '\0' suffix if its length is odd.
		str = strings.Trim(str, " \000")
	}

	return str
}
