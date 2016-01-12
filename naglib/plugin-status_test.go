package naglib

import (
	"testing"

	"github.com/daviddengcn/go-colortext"
	"github.com/stretchr/testify/assert"
)

func Test_PluginStatusStringConversion_OK(t *testing.T) {
	assert.Equal(t, "OK", OK.String())
}

func Test_PluginStatusStringConversion_WARNING(t *testing.T) {
	assert.Equal(t, "WARNING", WARNING.String())
}

func Test_PluginStatusStringConversion_CRITICAL(t *testing.T) {
	assert.Equal(t, "CRITICAL", CRITICAL.String())
}

func Test_PluginStatusStringConversion_UNKNOWN(t *testing.T) {
	assert.Equal(t, "UNKNOWN", UNKNOWN.String())
}

func Test_PluginStatusColorConversion_OK(t *testing.T) {
	assert.Equal(t, ct.Green, OK.Color())
}

func Test_PluginStatusColorConversion_WARNING(t *testing.T) {
	assert.Equal(t, ct.Yellow, WARNING.Color())
}

func Test_PluginStatusColorConversion_CRITICAL(t *testing.T) {
	assert.Equal(t, ct.Red, CRITICAL.Color())
}

func Test_PluginStatusColorConversion_UNKNOWN(t *testing.T) {
	assert.Equal(t, ct.Magenta, UNKNOWN.Color())
}
