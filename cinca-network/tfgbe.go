package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a forex
type SmartContract struct {
	contractapi.Contract
}

// forex describes basic details of what makes up a forex
type Forex struct {
	Forexid string `json:"forexid"`
	Fecha   string `json:"fecha"`
	Base string `json:"base"`
	Euros string `json:"euros"`
	Oro string `json:"oro"`
	Libra  string `json:"libra"`
	Bitcoin  string `json:"bitcoin"`
	Ethereum  string `json:"ethereum"`
}
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Forex
}


// InitLedger adds a base set of forexs to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	forexs := []Forex{
		Forex{Fecha:"2020-10-06 00:05:00+00", Base:"USD", Euros:"0.84823", Oro:"0.00052279", Libra:"0.769764", Bitcoin:"0.00009269872404841283", Ethereum:"0.002830175329361654"},

	}

	for i, forex := range forexs {
		forexAsBytes, _ := json.Marshal(forex)
		err := ctx.GetStub().PutState("Forex"+strconv.Itoa(i), forexAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateForex adds a new Forex to the world state with given details
func (s *SmartContract) CreateForex(ctx contractapi.TransactionContextInterface, forexid string, fecha string, base string, euros string, oro string, libra string, bitcoin string, ethereum string) error {
	forex := Forex{
		Fecha: fecha,
		Base: base,
		Euros: euros,
		Oro: oro,
		Libra: libra,
		Bitcoin: bitcoin,
		Ethereum: ethereum,
	}

	forexAsBytes, _ := json.Marshal(forex)

	return ctx.GetStub().PutState(forexid, forexAsBytes)
}



// QueryAllForexs returns all Forexs found in world state
func (s *SmartContract) QueryAllForexs(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		forex := new(Forex)
		_ = json.Unmarshal(queryResponse.Value, forex)

		queryResult := QueryResult{Key: queryResponse.Key, Record: forex}
		results = append(results, queryResult)
	}

	return results, nil
}


func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create forex chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting forex chaincode: %s", err.Error())
	}
}
