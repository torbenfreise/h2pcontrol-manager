package cmd

import (
	managergrpc "buf.build/gen/go/beyer-labs/h2pcontrol/grpc/go/h2pcontrol/manager/v1/managerv1grpc"
	managerpb "buf.build/gen/go/beyer-labs/h2pcontrol/protocolbuffers/go/h2pcontrol/manager/v1"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
)

var watch = &cobra.Command{
	Use:               "watch",
	Short:             "Watch active services",
	Long:              "Stream registry updates from the h2pcontrol manager, reprinting the service list on every change.",
	PersistentPreRun:  clientPreRun,
	PersistentPostRun: clientPostRun,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client, ok := ctx.Value(clientKey).(managergrpc.ManagerServiceClient)
		if !ok {
			log.Fatal("manager client not found in context")
		}

		stream, err := client.Watch(ctx, &managerpb.WatchRequest{})
		if err != nil {
			log.Fatalf("could not start watch stream: %v", err)
		}

		for {
			resp, err := stream.Recv()
			if err != nil {
				if errors.Is(err, io.EOF) || ctx.Err() != nil {
					return
				}
				log.Fatalf("watch stream error: %v", err)
			}
			fmt.Print("\033[H\033[2J")
			PrettyPrintServices(&managerpb.ListResponse{Services: resp.GetServices()})
		}
	},
}

func init() {
	registerClientFlags(watch)
	rootCmd.AddCommand(watch)
}
