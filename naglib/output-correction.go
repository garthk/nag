package naglib

import (
	"fmt"
	"strings"
)

const MAX_PLUGIN_OUTPUT_LENGTH = 4095

func DefaultOutputIfBlank(output string, status PluginStatus) string {
	if len(output) > 4 || len(strings.TrimSpace(output)) > 0 {
		return output
	}

	return fmt.Sprintf("%v\n", status)
}

func TruncateAndEnforceNewline(output string) string {
	return TruncateUTF8StringAtByte(output, MAX_PLUGIN_OUTPUT_LENGTH) + "\n";
}

func TruncateUTF8StringAtByte(str string, at int) string {
	if at >= len(str) {
		return str
	}
	if at <= 0 {
		return ""
	}
	maxidx := len(str) - 1
	for idx := at - 1; idx >= 0; idx-- {
		thisByte := str[idx]
		if !inCodePoint(thisByte) {
			return str[:idx+1] // up to thisByte inclusive
		}
		if idx < maxidx {
			nextByte := str[idx+1]
			if !inCodePoint(nextByte) || beginsCodePoint(nextByte) {
				// thisByte is the last in a code point
				return str[:idx+1]
			}
		}
	}
	return ""
}

func beginsCodePoint(c byte) bool {
	return c&128 == 128 && c&64 == 64
}

func inCodePoint(c byte) bool {
	return c&128 == 128
}
