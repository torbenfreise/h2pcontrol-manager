package cmd

import (
	"context"
	"fmt"
	"log"

	managergrpc "buf.build/gen/go/beyer-labs/h2pcontrol/grpc/go/h2pcontrol/manager/v1/managerv1grpc"
	managerpb "buf.build/gen/go/beyer-labs/h2pcontrol/protocolbuffers/go/h2pcontrol/manager/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type contextKey string

const (
	connKey   contextKey = "conn"
	clientKey contextKey = "client"
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
		PrettyPrintServices(r)
	},
}

func init() {
	viper.SetDefault("address", "127.0.0.1:50051")
	viper.SetEnvPrefix("H2PMANAGER")
	viper.AutomaticEnv()
	registerClientFlags(list)
	rootCmd.AddCommand(list)
}

func registerClientFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("address", "a", viper.GetString("address"), "Address of the h2pcontrol manager to query")
	if err := viper.BindPFlag("address", cmd.Flags().Lookup("address")); err != nil {
		log.Fatalf("Failed to bind viper flag: %v", err)
	}
}

func clientPreRun(cmd *cobra.Command, _ []string) {
	managerAddr, _ := cmd.Flags().GetString("address")
	conn, err := grpc.NewClient(managerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to initialise grpc client: %v", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, connKey, conn)
	client := managergrpc.NewManagerServiceClient(conn)
	ctx = context.WithValue(ctx, clientKey, client)
	cmd.SetContext(ctx)
}

func clientPostRun(cmd *cobra.Command, _ []string) {
	if conn, ok := cmd.Context().Value(connKey).(*grpc.ClientConn); ok && conn != nil {
		if err := conn.Close(); err != nil {
			log.Fatalf("Failed to close grpc connection: %v", err)
		}
	}
}

func PrettyPrintServices(resp *managerpb.ListResponse) {
	if len(resp.GetServices()) == 0 {
		fmt.Println("No servers found.")
		return
	}
	fmt.Println("Registered Services:")
	for _, server := range resp.GetServices() {
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
