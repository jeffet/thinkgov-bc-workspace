// Note that using the "debug test" options will fail due to this issue:
// https://github.com/golang/go/issues/23733

package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

func TestRead(t *testing.T) {
	// Init
	marblesChaincode := new(Chaincode)
	stub := shim.NewMockStub("marbles", marblesChaincode)

	// Put a test marble in the world state
	testMarble := Marble{ObjectType: "Marble", Color: "green", Size: 19, Owner: "Matthew"}
	marbleAsBytes, err := json.Marshal(testMarble)
	if err != nil {
		fmt.Println("Marshalling the testMarble failed", err)
		t.FailNow()
	}
	stub.MockTransactionStart("1")
	stub.PutState("MARBLE1", marbleAsBytes)
	stub.MockTransactionEnd("1")

	// Read out the marble
	res := stub.MockInvoke("marbles", [][]byte{[]byte("read"), []byte("MARBLE1")})
	if res.Status != shim.OK {
		fmt.Println("read failed", string(res.Message))
		t.FailNow()
	}
	checkRes(t, res, `{"objectType":"Marble","color":"green","size":19,"owner":"Matthew"}`)
}

func TestAddMarble(t *testing.T) {
	// Init
	marblesChaincode := new(Chaincode)
	stub := shim.NewMockStub("marbles", marblesChaincode)

	// Add a marble
	res := stub.MockInvoke("marbles", [][]byte{[]byte("addMarble"), []byte("MARBLE1"), []byte("pink"), []byte("6"), []byte("Tom")})
	if res.Status != shim.OK {
		fmt.Println("addMarble failed", string(res.Message))
		t.FailNow()
	}

	// MockInvoke covers queries and invokes
	res = stub.MockInvoke("marbles", [][]byte{[]byte("read"), []byte("MARBLE1")})
	checkRes(t, res, `{"objectType":"Marble","color":"pink","size":6,"owner":"Tom"}`)
}

func TestChangeOwner(t *testing.T) {
	// Init
	marblesChaincode := new(Chaincode)
	stub := shim.NewMockStub("marbles", marblesChaincode)

	// Add a marble
	res := stub.MockInvoke("marbles", [][]byte{[]byte("addMarble"), []byte("MARBLE1"), []byte("pink"), []byte("6"), []byte("Tom")})
	if res.Status != shim.OK {
		fmt.Println("addMarble failed", string(res.Message))
		t.FailNow()
	}

	// Change it's owner
	res = stub.MockInvoke("marbles", [][]byte{[]byte("changeOwner"), []byte("MARBLE1"), []byte("Matthew")})
	if res.Status != shim.OK {
		fmt.Println("changeOwner failed", string(res.Message))
		t.FailNow()
	}

	// MockInvoke covers queries and invokes
	res = stub.MockInvoke("marbles", [][]byte{[]byte("read"), []byte("MARBLE1")})
	checkRes(t, res, `{"objectType":"Marble","color":"pink","size":6,"owner":"Matthew"}`)
}

func TestDelete(t *testing.T) {
	// Init
	marblesChaincode := new(Chaincode)
	stub := shim.NewMockStub("marbles", marblesChaincode)

	// Add a marble
	res := stub.MockInvoke("marbles", [][]byte{[]byte("addMarble"), []byte("MARBLE1"), []byte("pink"), []byte("6"), []byte("Tom")})
	if res.Status != shim.OK {
		fmt.Println("addMarble failed", string(res.Message))
		t.FailNow()
	}

	res = stub.MockInvoke("marbles", [][]byte{[]byte("delete"), []byte("MARBLE1")})
	if res.Status != shim.OK {
		fmt.Println("delete failed", string(res.Message))
		t.FailNow()
	}

	// MockInvoke covers queries and invokes
	res = stub.MockInvoke("marbles", [][]byte{[]byte("read"), []byte("MARBLE1")})
	if res.Status != shim.OK {
		fmt.Println("read() failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload != nil {
		fmt.Println("read() returned a value", res.Payload)
		t.FailNow()
	}
}

func checkRes(t *testing.T, res pb.Response, expected string) {
	if res.Status != shim.OK {
		fmt.Println("read() failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("read() failed to get value")
		t.FailNow()
	}
	if string(res.Payload) != expected {
		fmt.Println("read()'s value was not as expected")
		t.FailNow()
	}
}
