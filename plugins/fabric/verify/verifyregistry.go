package verify

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zhigui-projects/sidemesh"
)

type Registry struct {
	contractapi.Contract
}

func (verifyRegistry *Registry) Register(ctx contractapi.TransactionContextInterface, network string, chain string, verifyInfo string) error {
	uri := sidemesh.SideMeshPrefix + network + chain + ":verify"
	return ctx.GetStub().PutState(uri, []byte(verifyInfo))
}

func (verifyRegistry *Registry) Resolve(ctx contractapi.TransactionContextInterface, network string, chain string) (string, error) {
	uri := sidemesh.SideMeshPrefix + network + chain + ":verify"
	verifyInfoBytes, err := ctx.GetStub().GetState(uri)
	if err != nil {
		return "", err
	}

	return string(verifyInfoBytes), nil
}
