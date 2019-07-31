package pkg

import (
	"strings"
)

func ExtractRoute(phpDocBlock string) (bool, string) {
	if strings.Contains(phpDocBlock, "@Route") {

		// split
		phpDocBlockStrings := strings.Split(phpDocBlock, "\n")

		for _, line := range phpDocBlockStrings {
			// Trim surrounding whitespace
			line = strings.TrimSpace(line)
			// Skip over the first DocBlock line
			if line == "/**" {
				continue
			}
			// Skip over the last DocBlock line
			if line == "*/" {
				continue
			}
			// TODO: Add support for multiline config. Current implementation expects definition on a single line.
			if strings.Contains(line, "@Route") {
				// Remove "* @Route(" from beginning and ")" from end of line
				line = line[9 : len(line)-1]
				rtn := ""
				lineParts := splitDocBlockLineParts(line)
				for i, part := range lineParts {
					part = strings.TrimSpace(part)
					key, value := extractKeyValuePair(part)
					if i > 0 {
						rtn += ", "
					}
					rtn = rtn + key + " => " + value
				}
				return true, "[" + rtn + "]"
			}
		}
		return true, phpDocBlock
	}
	return false, ""
}

func extractKeyValuePair(part string) (string, string) {
	openQuotePos := -1
	closeQuotePos := -1
	key := ""
	value := ""

	for i := 0; i < len(part)-1; i++ {
		c := string([]rune(part)[i])
		if c == "\"" || c == "'" || c == "{" || c == "}" {
			if openQuotePos == -1 {
				openQuotePos = i
				continue
			}
			closeQuotePos = i
		}
	}

	if openQuotePos < 1 {
		key = "_"
		openQuotePos = 0
	} else {
		key = part[0 : openQuotePos-1]
	}

	if closeQuotePos-1 > len(part)-1 {
		value = part[openQuotePos+1:]
	} else {
		if string([]rune(part)[len(part)-1]) == "}" {
			value = part[openQuotePos:]
		} else {
			value = part[openQuotePos+1 : len(part)-1]
		}
	}

	return key, value
}

func splitDocBlockLineParts(line string) []string {

	var parts []string
	var scope []string

	currentPartBeginPos := 0

	for i := 0; i < len(line)-1; i++ {
		switch string([]rune(line)[i]) {
		case "\"":
			if len(scope) > 0 && scope[len(scope)-1] == "\"" {
				scope = scope[:len(scope)-1]
				continue
			}
			scope = append(scope, "\"")

		case "'":
			if len(scope) > 0 && scope[len(scope)-1] == "'" {
				scope = scope[:len(scope)-1]
				continue
			}
			scope = append(scope, "'")

		case "{":
			scope = append(scope, "{")

		case "}":
			if len(scope) > 0 && scope[len(scope)-1] == "{" {
				scope = scope[:len(scope)-1]
				continue
			}

		case ",":
			if len(scope) == 0 {
				parts = append(parts, line[currentPartBeginPos:i])
				currentPartBeginPos = i + 1
			}
		}
	}

	// Get the final part
	parts = append(parts, line[currentPartBeginPos:])

	return parts
}
