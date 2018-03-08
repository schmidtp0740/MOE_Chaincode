package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {

}

// Chaincode ...
type Chaincode struct {
	// ....
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("error starting File Trace chaincode: %s", err)
	}
}
