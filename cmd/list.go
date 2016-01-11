package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List NRPE commands",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		return errors.New("can't yet list")
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
