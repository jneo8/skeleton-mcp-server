package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AddPersistentFlags adds persistent flags to the given cobra command
func AddPersistentFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("config", "", "config file path")
	cmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")

	if err := viper.BindPFlag("config", cmd.PersistentFlags().Lookup("config")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding config flag: %v\n", err)
		os.Exit(1)
	}
	if err := viper.BindPFlag("logging.level", cmd.PersistentFlags().Lookup("log-level")); err != nil {
		fmt.Fprintf(os.Stderr, "Error binding log-level flag: %v\n", err)
		os.Exit(1)
	}
}
