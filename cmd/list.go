package cmd

import (
	"fmt"
	"log"

	managergrpc "buf.build/gen/go/beyer-labs/h2pcontrol/grpc/go/h2pcontrol/manager/v1/managerv1grpc"
	managerpb "buf.build/gen/go/beyer-labs/h2pcontrol/protocolbuffers/go/h2pcontrol/manager/v1"
	"github.com/spf13/cobra"
)

var list = &cobra.Command{
	Use:               "list",
	Short:             "list active services",
	Long:              "list all active services from the h2p registry.",
	PersistentPreRun:  clientPreRun,
	PersistentPostRun: clientPostRun,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client := ctx.Value(clientKey).(managergrpc.ManagerServiceClient)

		r, err := client.List(ctx, &managerpb.ListRequest{})
		if err != nil {
			log.Fatalf("could not list services: %v", err)
		}
		PrettyPrintServices(r)
	},
}

func init() {
	initClientFlags(list)
	rootCmd.AddCommand(list)
}

func PrettyPrintServices(resp *managerpb.ListResponse) {
	if len(resp.Services) == 0 {
		fmt.Println("No servers found.")
		return
	}
	fmt.Println("Registered Services:")
	for _, server := range resp.Services {
		definition := server.GetDefinition()
		lastSeen := server.GetLastSeen().AsTime().Format("2006-01-02 15:04:05")
		status := map[bool]string{true: "healthy", false: "unhealthy"}[server.GetHealthy()]

		fmt.Printf("  %s\n", definition.GetName())
		fmt.Printf("    Description : %s\n", definition.GetDescription())
		fmt.Printf("    Address     : %s\n", definition.GetAddress())
		fmt.Printf("    Last seen   : %s\n", lastSeen)
		fmt.Printf("    Status      : %s\n", status)
		fmt.Println()
	}
}
