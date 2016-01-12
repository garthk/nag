package naglib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExitStatusToPluginStatus_DefaultTreatment_0(t *testing.T) {
	assert.Equal(t, OK, ExitStatusToPluginStatus(0, new(ExitStatusTreatment)))
}

func Test_ExitStatusToPluginStatus_DefaultTreatment_1(t *testing.T) {
	assert.Equal(t, WARNING, ExitStatusToPluginStatus(1, new(ExitStatusTreatment)))
}

func Test_ExitStatusToPluginStatus_DefaultTreatment_2(t *testing.T) {
	assert.Equal(t, CRITICAL, ExitStatusToPluginStatus(2, new(ExitStatusTreatment)))
}

func Test_ExitStatusToPluginStatus_DefaultTreatment_3(t *testing.T) {
	assert.Equal(t, UNKNOWN, ExitStatusToPluginStatus(3, new(ExitStatusTreatment)))
}

func Test_ExitStatusToPluginStatus_DefaultTreatment_23(t *testing.T) {
	assert.Equal(t, UNKNOWN, ExitStatusToPluginStatus(23, new(ExitStatusTreatment)))
}

func Test_ExitStatusToPluginStatus_CriticalWarnings_1(t *testing.T) {
	treatment := ExitStatusTreatment{CriticalWarnings: true}
	assert.Equal(t, CRITICAL, ExitStatusToPluginStatus(1, &treatment))
}

func Test_ExitStatusToPluginStatus_CriticalUnknowns_3(t *testing.T) {
	treatment := ExitStatusTreatment{CriticalUnknowns: true}
	assert.Equal(t, CRITICAL, ExitStatusToPluginStatus(3, &treatment))
}

func Test_ExitStatusToPluginStatus_CriticalExcess_23(t *testing.T) {
	treatment := ExitStatusTreatment{CriticalExcess: true}
	assert.Equal(t, CRITICAL, ExitStatusToPluginStatus(23, &treatment))
}
