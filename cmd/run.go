package cmd

import (
	"errors"
	"fmt"

	"github.com/garthk/nag/naglib"
	"github.com/garthk/nag/pkg/safe-colortext"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a Nagios plugin or other executable",
	Long: `
Run a Nagios plugin or other executable, constraining its output and exit
status for compatibility.

Use the -- option to end option parsing so you can give executables options.
`,
	Example: `  nag run /usr/lib/nagios/plugins/check_http http://localhost:8001
  nag run -WUX -- test -d /var/log/apache2`,
	DisableAutoGenTag: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		treatment, err := checkExitStatusFlags(cmd)
		if err != nil {
			return err
		}

		cfg, err := naglib.ParseConfig(nrpeCfgFile)
		if err != nil {
			return err
		}

		context := naglib.PluginContext{
			NagiosConfig:        cfg,
			ExitStatusTreatment: treatment,
		}

		if len(args) == 0 {
			return errors.New("nag run: executable required")
		} else {
			cmd.SilenceUsage = true
			result, err := naglib.RunPlugin(&context, args[0], args[1:]...)
			return processPluginResult(context, result, err)
		}
	},
}

func processPluginResult(context naglib.PluginContext, result naglib.PluginResult, err error) error {
	for _, msg := range context.Messages {
		printStatusLine(msg.Severity, msg.Message)
	}
	if err == nil {
		printStatusLine(result.Status, result.Output)
	}
	return err
}

func printStatusLine(status naglib.PluginStatus, msg string) {
	fmt.Printf("[")
	sc.Foreground(status.Color(), false)
	fmt.Printf("%s", status)
	sc.ResetColor()
	fmt.Printf("] %s\n", msg)
}

func init() {
	RootCmd.AddCommand(runCmd)
	addExitStatusFlags(runCmd)
}
