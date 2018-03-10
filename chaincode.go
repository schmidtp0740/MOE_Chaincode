package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode ...
type Chaincode struct {
	// ...
}

type data struct {
	ObjectType string `json:"docType"`
	ID         string `json:"id"`
	HeartRate  string `json:"heartRate"`
	Unit       string `json:"unit"`
	TimeStamp  string `json:"timeStamp"`
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("error starting File Trace chaincode: %s", err)
	}
}

// Init ...
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke ...
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	if function == "insertData" {
		return t.insertData(stub, args)
	} else if function == "getHistory" {
		return t.getHistory(stub, args)
	}

	fmt.Println("invoke did not find func: " + function)
	return shim.Error("Received unknown function invocation")
}

// insertHeartRate ...
func (t *Chaincode) insertData(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//	 0	     1	         2	        3
	// "id", "heartRate",  "unit", "timeStamp"
	if len(args) != 4 {
		return shim.Error("Incorrect number of aguements. Expecting 4")
	}
	for key, value := range args {
		if len(value) <= 0 {
			return shim.Error("Arguement (" + strconv.Itoa(key) + ") must be a non empty string")
		}
	}

	id := args[0]
	heartRate := args[1]
	unit := args[2]
	timeStamp := args[2]

	dataAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get id: " + err.Error())
	} else if dataAsBytes != nil {

		fmt.Println("This id already exists: " + id)

		dataToModify := data{}

		err = json.Unmarshal(dataAsBytes, &dataToModify)
		if err != nil {
			return shim.Error(err.Error())
		}

		dataToModify.HeartRate = heartRate
		dataToModify.TimeStamp = timeStamp

		dataJSONasBytes, _ := json.Marshal(dataToModify)
		err = stub.PutState(id, dataJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	} else {

		objectType := "heartRate"
		data := &data{objectType, id, heartRate, unit, timeStamp}
		dataJSONasBytes, err := json.Marshal(data)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.PutState(id, dataJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
	}

	fmt.Println("- end insertHeartRate")
	return shim.Success(nil)
}

func (t *Chaincode) getHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguements. Expecting 1")
	}

	id := args[0]

	fmt.Printf("- Start getHistory: %s", id)
	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		buffer.WriteString("{\"TxId\": \"" + response.TxId + "\"")
		buffer.WriteString(", \"Value\": " + string(response.Value))

		buffer.WriteString(", \"Timestamp\": \"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistory returning: %s", buffer.String())

	return shim.Success(buffer.Bytes())
}
