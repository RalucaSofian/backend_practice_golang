package utils

import "strings"

func CamelToPascalCase(input string) string {

	if len(input) > 0 {
		return strings.ToUpper(input[:1]) + input[1:]
	}

	return ""
}
