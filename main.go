package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	"gopkg.in/ArcadiaMediaGroup/zcoin-rosetta-node.v0/client"
	"gopkg.in/ArcadiaMediaGroup/zcoin-rosetta-node.v0/configuration"
	"gopkg.in/ArcadiaMediaGroup/zcoin-rosetta-node.v0/services"
)

// NewBlockchainRouter creates a blockchain specific router
// that will handle common routes specified inside the rosetta API specification
func NewBlockchainRouter(client client.ZcoinClient) http.Handler {
	assert, err := asserter.NewServer([]*types.NetworkIdentifier{
		{
			Blockchain:           client.GetConfig().NetworkIdentifier.Blockchain,
			Network:              client.GetConfig().NetworkIdentifier.Network,
			SubNetworkIdentifier: nil,
		},
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create asserter: %v\n", err)
		os.Exit(1)
	}

	networkAPIController := server.NewNetworkAPIController(services.NewNetworkAPIService(client), assert)
	blockAPIController := server.NewBlockAPIController(services.NewBlockAPIService(client), assert)
	return server.NewRouter(networkAPIController, blockAPIController)
}

func main() {
	configPath := os.Getenv(configuration.ConfigPath)
	if configPath == "" {
		configPath = "config.yaml"
	}

	cfg, err := configuration.New(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to parse config: %v\n", err)
		os.Exit(1)
	}

	client := client.NewZcoinClient(cfg)
	router := NewBlockchainRouter(client)
	fmt.Println("Listening on ", "0.0.0.0:"+cfg.Server.Port)
	err = http.ListenAndServe("0.0.0.0:"+cfg.Server.Port, router)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Zcoin Rosetta Gateway server exited suddenly: %v\n", err)
		os.Exit(1)
	}
}
