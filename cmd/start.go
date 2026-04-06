package cmd

import (
	"h2pcontrol.manager/internal/server"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the h2pcontrol manager server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		server.RunServer(port)
	},
}

func init() {
	startCmd.Flags().Int("port", 50051, "Port to listen on")
	rootCmd.AddCommand(startCmd)
}
