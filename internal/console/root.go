package console

import (
	"helpdesk-ticketing-system/internal/config"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "Todo service",
	Short: "Todo Service",
	Long:  `Todo Service`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	config.LoadWithViper()
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	config.SetupLogger()
}
