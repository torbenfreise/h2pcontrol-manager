package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:               "h2pmanager",
	Short:             "h2pmanager is the service registry for the h2pcontrol ecosystem",
	Long:              "h2pmanager is the service registry for the h2pcontrol ecosystem.\n\nUse 'start' to run the manager server and 'list' to query registered services.",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func init() {
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
