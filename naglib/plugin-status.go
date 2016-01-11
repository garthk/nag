package naglib

import (
	"github.com/daviddengcn/go-colortext"
)

type PluginStatus int

const (
	OK       PluginStatus = 0
	WARNING               = 1
	CRITICAL              = 2
	UNKNOWN               = 3
)

func (pluginStatus PluginStatus) String() string {
	switch {
	case pluginStatus == OK:
		return "OK"

	case pluginStatus == WARNING:
		return "WARNING"

	case pluginStatus == CRITICAL:
		return "CRITICAL"
	}

	return "UNKNOWN"
}

func (pluginStatus PluginStatus) Color() ct.Color {
	switch {
	case pluginStatus == OK:
		return ct.Green

	case pluginStatus == WARNING:
		return ct.Yellow

	case pluginStatus == CRITICAL:
		return ct.Red
	}

	return ct.Magenta
}
