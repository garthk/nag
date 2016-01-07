package cmd

import (
	"fmt"
	"github.com/garthk/nag/pkg/plugin-runner"
	"github.com/garthk/nag/pkg/safe-colortext"
	"github.com/spf13/cobra"
	"log"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a Nagios plugin or NRPE command",
	Long: `
Run a Nagios plugin or NRPE command, constraining their output and exit
status for compatibility.

If --sensitive, run any command as if it's an NRPE command, treating its
exit status as CRITICAL if non-zero.

Timeouts and problems running the command are treated as UNKNOWN, even if
run with --sensitive.

Use the -- option to end option parsing so you can use options in command
arguments, e.g. nag run -S -- /bin/test -d /var/log/apache2
`,
	Run: func(cmd *cobra.Command, args []string) {
		sensitive, err := cmd.Flags().GetBool("sensitive")
		if err != nil {
			panic(fmt.Sprintf("can't fetch sensitivity: %v", err))
		}

		options := runner.PluginRunOptions{
			Timeout:   runner.DEFAULT_TIMEOUT,
			Sensitive: sensitive,
		}
		if len(args) == 0 {
			// TODO: implement finding commands
			log.Printf("can't yet run all; echoing...\n")
			result, err := runner.RunPlugin(options, "echo", "1234567890")
			handle(result, err)
		} else {
			// TODO: implement finding commands
			// TODO: distinguish between commands and executables
			log.Printf("can't yet find commands; assuming executable...\n")
			result, err := runner.RunPlugin(options, args[0], args[1:]...)
			handle(result, err)
		}
	},
}

func handle(result runner.PluginResult, err error) {
	if err != nil {
		log.Println("command failed:", err.Error())
	}

	fmt.Printf("[")
	sc.Foreground(result.Status.Color(), false)
	fmt.Printf("%s", result.Status)
	sc.ResetColor()
	fmt.Printf("] %s", result.Output)
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("sensitive", "S", false, "Upgrade all >0 exit status to CRITICAL")
}

