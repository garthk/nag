package naglib

import (
	"testing"
)

func check(t *testing.T, status PluginStatus, expected string) {
	actual := status.String()
	if actual != expected {
		t.Errorf("%s != %s", actual, expected)
	}
}

// TODO: testify

func Test_PluginStatusStringConversion_OK(t *testing.T) {
	check(t, OK, "OK")
}

func Test_PluginStatusStringConversion_WARNING(t *testing.T) {
	check(t, WARNING, "WARNING")
}

func Test_PluginStatusStringConversion_CRITICAL(t *testing.T) {
	check(t, CRITICAL, "CRITICAL")
}

func Test_PluginStatusStringConversion_UNKNOWN(t *testing.T) {
	check(t, UNKNOWN, "UNKNOWN")
}
