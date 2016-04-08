package main

import (
    "errors"
    "fmt"

    "github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

const data_hash = "dataHash"
const user_key = "userKey"
const idp_key = "idpKey"

// SimpleChaincode example simple Chaincode implementation
type BundleChaincode struct {
}

// should call with args dataHash, userKey, idpKey
func (t *BundleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var err error
    
    if len(args) != 3 {
        return nil, errors.New("Incorrect number of arguments. Expecting 3")
    }
    
    dataHash := args[0]
    userKey := args[1]
    idpKey := args[2]
    
    fmt.Printf("dataHash = %s, userKey = %s, idpKey = %s\n", dataHash, userKey, idpKey)
    
    err = stub.PutState(data_hash, []byte(dataHash))
    if err != nil {
        return nil, err
    }
    
    err = stub.PutState(user_key, []byte(userKey))
    if err != nil {
        return nil, err
    }
    
    err = stub.PutState(idp_key, []byte(idpKey))
    if err != nil {
        return nil, err
    }
    
    return nil, nil
}

func (t *BundleChaincode) query(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var err error
    
    dataHash, err := stub.GetState(data_hash)
    if err != nil {
        return nil, err
    }
    userKey, err := stub.GetState(user_key)
    if err != nil {
        return nil, err
    }
    idpKey, err := stub.GetState(idp_key)
    if err != nil {
        return nil, err
    }
    
    json := fmt.Sprintf("{\"dataHash\":\"%s\",\"userKey\":\"%s\",\"idpKey\":\"%s\"}", dataHash, userKey, idpKey)
    return []byte(json), nil
}

func (t *BundleChaincode) getUserKey(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var err error
    
    userKey, err := stub.GetState(user_key)
    if err != nil {
        return nil, err
    }
    return userKey, nil
}

func (t *BundleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    if function == "query" {
        // Initialize the entities and their asset holdings
        return t.query(stub, args)
    } else if function == "getUserKey" {
        return t.getUserKey(stub, args)
    }
    
    return nil, errors.New("Received unknown function query")
}

func (t *BundleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    if function == "init" {
        // Initialize the entities and their asset holdings
        return t.init(stub, args)
    }
    
    return nil, errors.New("Received unknown function invocation")
}

func main() {
    err := shim.Start(new(BundleChaincode))
    if err != nil {
        fmt.Printf("Error starting Bundle chaincode: %s", err)
    }
}
