package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check the result of a local NRPE command",
	Long: `
Check an NRPE command, constraining its output and exit status for
compatibility.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return errors.New("can't yet check")
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
	addExitStatusFlags(runCmd)
}
