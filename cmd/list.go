package cmd

import (
	"log"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List NRPE commands",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("can't yet find commands\n")
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
