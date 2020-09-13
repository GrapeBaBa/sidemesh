package resource

import (
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zhigui-projects/sidemesh"
	"github.com/zhigui-projects/sidemesh/pb"
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
	event := &pb.ResourceRegisteredOrUpdatedEvent{Uri: &pb.URI{Network: network, Chain: chain}, Connection: []byte(connection), Type: pb.ChainType_FABRIC}
	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return err
	}
	return ctx.GetStub().SetEvent(sidemesh.Prefix+"RESOURCE_REGISTERED_EVENT", eventBytes)
}

func (resourceRegistry *Registry) Resolve(ctx contractapi.TransactionContextInterface, network string, chain string) (string, error) {
	uri := sidemesh.Prefix + network + chain + ":connection"
	conn, err := ctx.GetStub().GetState(uri)
	return string(conn), err
}
