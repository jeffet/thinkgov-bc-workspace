/*
	This is just a copy of example02.
*/

package main

import (
	"bytes"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// StringsChaincode - The Shim requires chaincode be structured as a struct with the
// 					 Invoke() and Init() functions attacked to it
// ===================================================================================
type StringsChaincode struct {
}

// main - Call the shim to fire up a new instance of this chaincode object
// ===================================================================================
func main() {
	err := shim.Start(new(StringsChaincode))
	if err != nil {
		fmt.Printf("Error starting Strings: %s", err)
	}
}

// Init - Setup the chaincode ready for exection
// ===================================================================================
func (t *StringsChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	//stub.PutState("Test", []byte("Hello, World!"))
	//return shim.Success(nil)

	// Declare variable
	fmt.Printf("Declaring")
	key := "Test"
	value := "Hello, World!"
	// Place into world state
	fmt.Printf("PutState")
	stub.PutState(key, []byte(value))
	// Exit with success
	return shim.Success(nil)
}

// read - Grab the value from the world state and return it
// ===================================================================================
func (t *StringsChaincode) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	valueAsBytes, _ := stub.GetState(args[0])
	return shim.Success(valueAsBytes)
}

// overwrite - Change the value in the world state to supplied argument
// ===================================================================================
func (t *StringsChaincode) overwrite(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	stub.PutState(args[0], []byte(args[1]))
	return shim.Success(nil)
}

// append - Add the submitted string to the
// ===================================================================================
func (t *StringsChaincode) append(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	oldValueAsBytes, _ := stub.GetState(args[0])
	var newValue bytes.Buffer
	newValue.Write(oldValueAsBytes)
	newValue.WriteString(args[1])
	stub.PutState(args[0], newValue.Bytes())
	return shim.Success(nil)
}

// delete - Remove a value from the world state and return it
// ===================================================================================
func (t *StringsChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	stub.DelState(args[0])
	return shim.Success(nil)
}

// Invoke - Router for the invoke/query requests
// ===================================================================================
func (t *StringsChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	// Grab the inputs to the transaction
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke has been called with", function, "and", args)

	// Call the apropriate function
	if function == "read" {
		return t.read(stub, args)
	} else if function == "add" {
		return t.overwrite(stub, args)
	} else if function == "overwrite" {
		return t.overwrite(stub, args)
	} else if function == "append" {
		return t.append(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	}

	// If 'function' doesn't match any of the above then return an error
	return shim.Error("Invalid invoke function name.")
}
