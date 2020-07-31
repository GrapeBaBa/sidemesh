package resource

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Registry struct {
	contractapi.Contract
}

func (resourceRegistry *Registry) Register(ctx contractapi.TransactionContextInterface, network string, chain string, connection []byte) error {
	uri := network + chain + ":connection"
	err := ctx.GetStub().PutState(uri, connection)
	if err != nil {
		return err
	}

	return ctx.GetStub().SetEvent("RESOURCE_REGISTERED_EVENT", connection)
}

func (resourceRegistry *Registry) Resolve(ctx contractapi.TransactionContextInterface, network string, chain string) ([]byte, error) {
	uri := network + chain + ":connection"
	return ctx.GetStub().GetState(uri)
}
