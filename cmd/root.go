package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
)

var cfgFile string

var RootCmd = &cobra.Command{
	Use:   "nag",
	Short: "Nag yourself to fix all the things",
	Long: `
Nag yourself to fix anything broken, by running Nagios plugins for you.

Read https://github.com/garthk/nag for more details.

Examples:
  nag list
  nag run
  nag run check_load
  nag run -- /usr/lib/nagios/plugins/check_load -w 15,10,5 -c 30,25,20
  nag run -S -- /bin/test -d /var/log/apache2
`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nag.yaml)")
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
