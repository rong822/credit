package main

import (
	"gitlab.bigtree.com/dashu-blockchain/credit/pkg/chaincode/points"
	"os"

	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
	shim.SetLoggingLevel(logLevel)

	err := shim.Start(new(points.Chaincode))
	if err != nil {
		log.Fatalf("error starting chaincode: %s\n", err)
	}
}
