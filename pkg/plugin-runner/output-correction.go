package runner

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
	length := len(output)
	switch {
	case length == 0:
		return "\n"

	case length > (MAX_PLUGIN_OUTPUT_LENGTH - 1):
		return output[:MAX_PLUGIN_OUTPUT_LENGTH-1] + "\n"

	case output[length-1] == '\n':
		return output
	}
	return output + "\n"
}
