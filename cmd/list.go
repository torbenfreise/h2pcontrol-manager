package cmd

import (
	"log"

	managergrpc "buf.build/gen/go/beyer-labs/h2pcontrol/grpc/go/h2pcontrol/manager/v1/managerv1grpc"
	managerpb "buf.build/gen/go/beyer-labs/h2pcontrol/protocolbuffers/go/h2pcontrol/manager/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var list = &cobra.Command{
	Use:               "list",
	Short:             "list active services",
	Long:              "list all services registered with the h2pcontrol manager.",
	PersistentPreRun:  clientPreRun,
	PersistentPostRun: clientPostRun,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		client, ok := ctx.Value(clientKey).(managergrpc.ManagerServiceClient)
		if !ok {
			log.Fatal("manager client not found in context")
		}

		r, err := client.List(ctx, &managerpb.ListRequest{})
		if err != nil {
			log.Fatalf("could not list services: %v", err)
		}
		printServices(r.GetServices())
	},
}

func init() {
	viper.SetDefault("address", "127.0.0.1:50051")
	viper.SetEnvPrefix("H2PMANAGER")
	viper.AutomaticEnv()
	registerClientFlags(list)
	rootCmd.AddCommand(list)
}
