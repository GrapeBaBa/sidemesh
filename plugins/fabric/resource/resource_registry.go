package resource

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zhigui-projects/sidemesh"
)

type Registry struct {
	contractapi.Contract
}

func (resourceRegistry *Registry) Register(ctx contractapi.TransactionContextInterface, network string, chain string, connection string) error {
	uri := sidemesh.Prefix + network + chain + ":connection"
	err := ctx.GetStub().PutState(uri, []byte(connection))
	if err != nil {
		return err
	}

	return ctx.GetStub().SetEvent("RESOURCE_REGISTERED_EVENT", []byte(connection))
}

func (resourceRegistry *Registry) Resolve(ctx contractapi.TransactionContextInterface, network string, chain string) (string, error) {
	uri := sidemesh.Prefix + network + chain + ":connection"
	conn, err := ctx.GetStub().GetState(uri)
	return string(conn), err
}
