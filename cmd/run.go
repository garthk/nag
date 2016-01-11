package cmd

import (
	"errors"
	"fmt"
	"log"

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
	RunE: func(cmd *cobra.Command, args []string) error {

		treatment, err := checkExitStatusFlags(cmd)
		if err != nil {
			return err
		}

		options := naglib.PluginRunOptions{
			Timeout:   naglib.DEFAULT_TIMEOUT,
			Treatment: treatment,
		}

		if len(args) == 0 {
			return errors.New("nag run: command required")
		} else {
			log.Printf("can't yet find commands; assuming executable...\n")
			cmd.SilenceUsage = true
			result, err := naglib.RunPlugin(options, args[0], args[1:]...)
			return processPluginResult(result, err)
		}
	},
}

func processPluginResult(result naglib.PluginResult, err error) error {
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
