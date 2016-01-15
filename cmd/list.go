package cmd

import (
	"fmt"

	"github.com/garthk/nag/naglib"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List NRPE commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		cfg, err := naglib.ParseConfig(nrpeCfgFile)
		if err != nil {
			return err
		}

		for k, v := range cfg.Commands {
			fmt.Printf("command[%s]=%s\n", k, v)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
