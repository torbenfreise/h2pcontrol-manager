package cmd

import (
	"github.com/spf13/cobra"
	"github.com/torbenfreise/h2pcontrol-manager/internal/server"
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
