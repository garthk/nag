package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/garthk/nag/naglib"
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
		treatment, err := checkExitStatusFlags(cmd)
		if err != nil {
			return err
		}

		var command string
		switch len(args) {
		case 0:
			return errors.New(fmt.Sprintf("%s: command required", cmd.CommandPath()))
		case 1:
			command = args[0]
		default:
			// TODO: check AllowArguments; process if allowed
			return errors.New(fmt.Sprintf("%s: arguments not supported", cmd.CommandPath()))
		}

		cmd.SilenceUsage = true

		cfg, err := naglib.ParseConfig(nrpeCfgFile)
		if err != nil {
			return err
		}

		commandline := cfg.Commands[command]
		if commandline == "" {
			// TODO: exit CRITICAL
			return errors.New(fmt.Sprintf("%s: command %s not known; try %s",
				cmd.CommandPath(),
				command,
				listCmd.CommandPath()))
		}

		context := naglib.PluginContext{
			NagiosConfig:        cfg,
			ExitStatusTreatment: treatment,
		}

		// TODO: add command prefix
		// TODO: substitute arguments
		shell := os.ExpandEnv("$SHELL")
		if len(shell) == 0 {
			shell = "/bin/bash"
		}
		result, err := naglib.RunPlugin(&context, shell, "-c", commandline)
		return processPluginResult(context, result, err)
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
	addExitStatusFlags(runCmd)
}
