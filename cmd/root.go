package cmd

import (
	"log"
	"os"

	"github.com/garthk/nag/naglib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "nag",
	Short: "Nag yourself to fix all the things",
	Long: `
Nag yourself to fix anything broken, by running Nagios plugins for you.

Read https://github.com/garthk/nag for more details.`,
	Example: `  nag list
  nag check
  nag check check_load
  nag run -- /usr/lib/nagios/plugins/check_load -w 15,10,5 -c 30,25,20
  nag run -S -- /bin/test -d /var/log/apache2`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		printStatusLine(naglib.UNKNOWN, err.Error())
		os.Exit(naglib.UNKNOWN)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nag.yaml)")
	RootCmd.PersistentFlags().BoolP("critical-warnings", "W", false, "Upgrade WARNING to CRITICAL")
	RootCmd.PersistentFlags().BoolP("critical-unknowns", "U", false, "Upgrade UNKNOWN to CRITICAL")
	RootCmd.PersistentFlags().BoolP("critical-excess", "X", false, "Upgrade excess (>3) exit status to CRITICAL")
	RootCmd.PersistentFlags().BoolP("tolerant", "x", false, "Pass excess (>3) exit status as-is, not as UNKNOWN")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".nag")
	viper.AddConfigPath("$HOME")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
