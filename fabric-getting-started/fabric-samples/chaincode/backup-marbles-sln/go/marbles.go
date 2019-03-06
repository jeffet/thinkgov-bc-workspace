package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//Chaincode Struct
type Chaincode struct {
}

//Marble Struct
type Marble struct {
	ObjectType string `json:"objectType"`
	Color      string `json:"color"`
	Size       int    `json:"size"`
	Owner      string `json:"owner"`
}

func main() {
	err := shim.Start(new(Chaincode))
	if err != nil {
		fmt.Printf("Error starting Chaincode: %s", err)
	}
}

//Init Chaincode
func (t *Chaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *Chaincode) initMarbles(stub shim.ChaincodeStubInterface) pb.Response {
	marbles := []Marble{
		Marble{ObjectType: "Marble", Color: "blue", Size: 20, Owner: "Tom"},
		Marble{ObjectType: "Marble", Color: "yellow", Size: 10, Owner: "Tom"},
		Marble{ObjectType: "Marble", Color: "red", Size: 22, Owner: "Tom"},
		Marble{ObjectType: "Marble", Color: "blue", Size: 21, Owner: "Tom"},
		Marble{ObjectType: "Marble", Color: "red", Size: 13, Owner: "Tom"},
		Marble{ObjectType: "Marble", Color: "yellow", Size: 34, Owner: "Tom"},
		Marble{ObjectType: "Marble", Color: "blue", Size: 25, Owner: "Matthew"},
		Marble{ObjectType: "Marble", Color: "red", Size: 9, Owner: "Matthew "},
		Marble{ObjectType: "Marble", Color: "blue", Size: 16, Owner: "Matthew "},
		Marble{ObjectType: "Marble", Color: "blue", Size: 13, Owner: "Matthew "},
	}

	i := 0
	for i < len(marbles) {
		fmt.Println("i is ", i)
		marbleAsBytes, _ := json.Marshal(marbles[i])
		stub.PutState("MARBLE"+strconv.Itoa(i), marbleAsBytes)
		fmt.Println("Added", marbles[i])
		i = i + 1
	}

	return shim.Success(nil)
}

//Invoke Chaincode
func (t *Chaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "read" {
		return t.read(stub, args)
	} else if function == "addMarble" {
		return t.addMarble(stub, args)
	} else if function == "changeOwner" {
		return t.changeOwner(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "initMarbles" {
		return t.initMarbles(stub)
	}

	return shim.Error("Invalid invoke function name.")
}

func (t *Chaincode) read(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	valueAsBytes, _ := stub.GetState(args[0])
	return shim.Success(valueAsBytes)
}

func (t *Chaincode) addMarble(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	size, _ := strconv.Atoi(args[2])
	var marble = Marble{ObjectType: "Marble", Color: args[1], Size: size, Owner: args[3]}

	marbleAsBytes, _ := json.Marshal(marble)
	stub.PutState(args[0], marbleAsBytes)

	return shim.Success(nil)
}

func (t *Chaincode) changeOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	marbleAsBytes, _ := stub.GetState(args[0])
	marble := Marble{}

	json.Unmarshal(marbleAsBytes, &marble)
	marble.Owner = args[1]

	marbleAsBytes, _ = json.Marshal(marble)
	stub.PutState(args[0], marbleAsBytes)

	return shim.Success(nil)
}

func (t *Chaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	stub.DelState(args[0])
	return shim.Success(nil)
}
