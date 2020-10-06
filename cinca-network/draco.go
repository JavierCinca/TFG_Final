package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a iot
type SmartContract struct {
	contractapi.Contract
}

// iot describes basic details of what makes up a iot
type Iot struct {
	Id   string `json:"id"`
	Time string `json:"time"`
	Tipo string `json:"tipo"`
	State string `json:"state"`
	Id_iot  string `json:"id_iot"`
}

// QueryResult structure used for handling result of query
type QueryResult struct {
	Key    string `json:"Key"`
	Record *Iot
}

// InitLedger adds a base set of iots to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	iots := []Iot{
		Iot{Id: "Lamp:001", Time: "2020-09-16T14:07:32.00Z", Tipo: "Lamp", State: "ON", Id_iot: "urn:ngsi-ld:Store:001" },
		Iot{Id: "Lamp:002", Time: "2020-09-16T14:07:32.00Z", Tipo: "Lamp", State: "ON", Id_iot: "urn:ngsi-ld:Store:002" },
	}

	for i, iot := range iots {
		iotAsBytes, _ := json.Marshal(iot)
		err := ctx.GetStub().PutState("Iot"+strconv.Itoa(i), iotAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

// CreateIot adds a new iot to the world state with given details
func (s *SmartContract) CreateIot(ctx contractapi.TransactionContextInterface, iotNumber string, id string, time string, tipo string, state string, id_iot string) error {
	iot := Iot{
		Id: id,
		Time: time,
		Tipo: tipo,
		State: state,
		Id_iot: id_iot,
	}

	iotAsBytes, _ := json.Marshal(iot)

	return ctx.GetStub().PutState(iotNumber, iotAsBytes)
}

// QueryIot returns the iot stored in the world state with given id
func (s *SmartContract) QueryIot(ctx contractapi.TransactionContextInterface, iotNumber string) (*Iot, error) {
	iotAsBytes, err := ctx.GetStub().GetState(iotNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if iotAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", iotNumber)
	}

	iot := new(Iot)
	_ = json.Unmarshal(iotAsBytes, iot)

	return iot, nil
}

// QueryAllIots returns all iots found in world state
func (s *SmartContract) QueryAllIots(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
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

		iot := new(Iot)
		_ = json.Unmarshal(queryResponse.Value, iot)

		queryResult := QueryResult{Key: queryResponse.Key, Record: iot}
		results = append(results, queryResult)
	}

	return results, nil
}


func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabiot chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabiot chaincode: %s", err.Error())
	}
}
