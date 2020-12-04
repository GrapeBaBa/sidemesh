package verify

import (
	"encoding/json"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zhigui-projects/sidemesh"
)

type Registry struct {
	contractapi.Contract
}

func (verifyRegistry *Registry) Register(ctx contractapi.TransactionContextInterface, network string, chain string, contract string, function string) error {
	uri := sidemesh.Prefix + network + chain + ":verify"
	verifyInfo := &sidemesh.VerifyInfo{Contract: contract, Function: function}
	verifyInfoBytes, err := json.Marshal(verifyInfo)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(uri, verifyInfoBytes)
}

func (verifyRegistry *Registry) Resolve(ctx contractapi.TransactionContextInterface, network string, chain string) (string, error) {
	uri := sidemesh.Prefix + network + chain + ":verify"
	verifyInfoBytes, err := ctx.GetStub().GetState(uri)
	if err != nil {
		return "", err
	}

	return string(verifyInfoBytes), nil
}
