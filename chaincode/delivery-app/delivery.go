package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// SmartContract structure was defined
type SmartContract struct {
}

// Trans struct with 3 properties was defined
type Trans struct {
	ID       string `json:"id"`
	Paid     string `json:"paid"`
	Received string `json:"received"`
}

// Init method
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

// Invoke method
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "queryTrans" {
		return s.queryTrans(APIstub, args)
	} else if function == "newTrans" {
		return s.newTrans(APIstub, args)
	} else if function == "queryAllTrans" {
		return s.queryAllTrans(APIstub)
	} else if function == "updateState" {
		return s.updateState(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// initLedger
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	trans := []Trans{}

	i := 0
	for i < len(trans) {
		fmt.Println("i is ", i)
		transAsBytes, _ := json.Marshal(trans[i])
		APIstub.PutState(strconv.Itoa(i+1), transAsBytes)
		fmt.Println("Added", trans[i])
		i = i + 1
	}

	return shim.Success(nil)
}

//queryClient
func (s *SmartContract) queryTrans(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	clientAsBytes, _ := APIstub.GetState(args[0])
	if clientAsBytes == nil {
		return shim.Error("Could not find transaction")
	}
	return shim.Success(clientAsBytes)
}

// newClient
func (s *SmartContract) newTrans(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var trans = Trans{ID: args[0], Paid: args[2], Received: args[3]}

	transAsBytes, _ := json.Marshal(trans)
	err := APIstub.PutState(args[0], transAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to creat new transaction: %s", args[0]))
	}

	return shim.Success(nil)
}

// query all trans
func (s *SmartContract) queryAllTrans(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add comma before array members,suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllTrans:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) updateState(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	transAsBytes, _ := APIstub.GetState(args[0])
	if transAsBytes == nil {
		return shim.Error("Could not find transaction")
	}

	var trans = Trans{ID: args[0], Paid: args[2], Received: args[3]}

	stateAsBytes, _ := json.Marshal(trans)
	err := APIstub.PutState(args[0], stateAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to update transaction: %s", args[0]))
	}

	return shim.Success(nil)
}

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
