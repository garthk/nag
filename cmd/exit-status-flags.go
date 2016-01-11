package cmd

import (
	"errors"
	"fmt"

	"github.com/garthk/nag/naglib"
	"github.com/spf13/cobra"
)

func checkExitStatusFlags(cmd *cobra.Command) (*naglib.ExitStatusTreatment, error) {
	flags := new(naglib.ExitStatusTreatment)

	flags.CriticalWarnings = getBoolOrPanic(cmd, "critical-warnings")
	flags.CriticalUnknowns = getBoolOrPanic(cmd, "critical-unknowns")
	flags.CriticalExcess = getBoolOrPanic(cmd, "critical-excess")
	flags.TolerantExcess = getBoolOrPanic(cmd, "tolerant")

	if flags.CriticalExcess && flags.TolerantExcess {
		return flags, errors.New("--tolerant clashes with --critical-excess")
	}

	return flags, nil
}

func getBoolOrPanic(cmd *cobra.Command, flagName string) bool {
	result, err := cmd.Flags().GetBool(flagName)

	if err != nil {
		panic(fmt.Sprintf("can't fetch %s: %v", flagName, err))
	}

	return result
}

func addExitStatusFlags(cmd *cobra.Command) {
}
