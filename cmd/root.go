package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	managergrpc "buf.build/gen/go/beyer-labs/h2pcontrol/grpc/go/h2pcontrol/manager/v1/managerv1grpc"
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

var rootCmd = &cobra.Command{
	Use:   "h2pcontrol",
	Short: "h2pcontrol is a tool for managing grpc communication between different services",
	Long:  "h2pcontrol is a tool for managing grpc communication between different services. This is the h2pcontrol client which allows you to register your service and consume other services. ",
}

func init() {}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}

func initClientFlags(cmd *cobra.Command) {
	viper.SetDefault("manager", "127.0.0.1:50051")
	viper.SetEnvPrefix("H2PCONTROL")
	viper.AutomaticEnv()
	cmd.Flags().String("manager", viper.GetString("manager"), "Address of the manager gRPC server")
	err := viper.BindPFlag("manager", cmd.Flags().Lookup("manager"))
	if err != nil {
		log.Fatalf("Failed to bind viper flag: %v", err)
	}
}

func clientPreRun(cmd *cobra.Command, _ []string) {
	managerAddr, _ := cmd.Flags().GetString("manager")
	conn, err := grpc.NewClient(managerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to initialise grpc client: %v", err)
	}

	// Store connection and context in command
	ctx := context.Background()
	ctx = context.WithValue(ctx, connKey, conn)
	client := managergrpc.NewManagerServiceClient(conn)
	ctx = context.WithValue(ctx, clientKey, client)
	cmd.SetContext(ctx)
}

func clientPostRun(cmd *cobra.Command, _ []string) {
	if conn, ok := cmd.Context().Value(connKey).(*grpc.ClientConn); ok && conn != nil {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close grpc connection: %v", err)
		}
	}
}
