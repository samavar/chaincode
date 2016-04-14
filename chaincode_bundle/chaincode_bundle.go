package main

import (
    "errors"
    "fmt"
    "encoding/json"
    "github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

const data_hash = "dataHash"
const user_key = "userKey"
const idp_key = "idpKey"
var documentIndex = "_documentIndex"

// SimpleChaincode example simple Chaincode implementation
type DocumentChaincode struct {
}

type Document struct {
    DataHash, UserKey, IdpKey string
}

func (t *DocumentChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var err error
    var empty []string
    
    
    if len(args) != 0 && len(args) !=3 {
        return nil, errors.New("Incorrect number of arguments. Expecting 0 or 3")
    }
    
    if len(args) == 3 {
        fmt.Printf("deploying, number of arguments is 3")
        _, err = t.insert(stub, args)
        if err != nil {
            return nil, err
        }
        empty = append(empty, args[0])        
    }    
    fmt.Printf("inserting document hashesh", empty)
    
    jsonAsBytes, _ := json.Marshal(empty)
    
    err = stub.PutState(documentIndex, jsonAsBytes)
    if err != nil {
        return nil, err
    }
    
    return nil, nil
}

func (t *DocumentChaincode) query(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var err error
    
    if len(args) != 1 {
        return nil, errors.New("Incorrect number of arguments. Expecting 1")
    }
    
    document_hash := args[0]
    document, err := stub.GetState(document_hash)
    if err != nil {
        return nil, err
    }

    return document, nil
}

func (t *DocumentChaincode) getUserKey(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var err error
    
    userKey, err := stub.GetState(user_key)
    if err != nil {
        return nil, err
    }
    return userKey, nil
}

func (t* DocumentChaincode) getAllIndexes(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var err error
    
    allHashes, err := stub.GetState(documentIndex)
    if err != nil {
        return nil, err
    }
    
    return allHashes, nil
}

func (t *DocumentChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
    if function == "query" {
        // Initialize the entities and their asset holdings
        return t.query(stub, args)
    } else if function == "getUserKey" {
        return t.getUserKey(stub, args)
    } else if function == "getAllIndexes" {
        return t.getAllIndexes(stub, args)
    }
    
    return nil, errors.New("Received unknown function query")
}

func (t *DocumentChaincode) insert(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
    var err error
    
    if len(args) != 3 {
        return nil, errors.New("Incorrect number of arguments. Expecting 3")
    }
    document := Document{DataHash: args[0], UserKey: args[1], IdpKey : args[2]}
    fmt.Printf("inserting document with dataHash = %s, userKey = %s, idpKey = %s\n", document.DataHash, document.UserKey, document.IdpKey)
    
    jsonAsBytes, _ := json.Marshal(document)
    
    err = stub.PutState(document.DataHash, jsonAsBytes);
    if err != nil {
        return nil, err
    }
    
    
    jsonAsBytes, err = stub.GetState(documentIndex)
    if err != nil {
        return nil, err
    }
    var allHashes []string
    _ = json.Unmarshal(jsonAsBytes, &allHashes)
    
    allHashes = append(allHashes, document.DataHash)
    
    jsonAsBytes,_ = json.Marshal(allHashes)
    err = stub.PutState(documentIndex, jsonAsBytes)
    if(err != nil) {
        return nil, err
    }
    
    
    return nil, nil
}

func (t *DocumentChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {        
    if function == "init" {
        // Initialize the entities and their asset holdings
        return t.init(stub, args)
    } else if function == "insert" {
        return t.insert(stub, args);
    }
    
    return nil, errors.New("Received unknown function invocation")
}

func main() {
    err := shim.Start(new(DocumentChaincode))
    if err != nil {
        fmt.Printf("Error starting Bundle chaincode: %s", err)
    }
}