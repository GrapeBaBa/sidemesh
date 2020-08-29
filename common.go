package sidemesh

import (
	"github.com/zhigui-projects/sidemesh/pb"
)

const SideMeshPrefix string = "SIDE_MESH_"

type VerifyInfo struct {
	Contract string `json:"contract"`
	Function string `json:"function"`
}

type GlobalTransactionManager interface {
	StartGlobalTransaction(ttlHeight uint64, ttlTime string) error
	RegisterBranchTransaction(network string, chain string, contract string, function string, args []string) error
	PreparePrimaryTransaction() error
	StartBranchTransaction(primaryNetwork string, primaryChain string, primaryTxID string, primaryTxProof string) error
	PrepareBranchTransaction(primaryNetwork string, primaryChain string, primaryTxID string, globalTxQueryContract string, globalTxQueryFunc string) error
	ConfirmPrimaryTransaction(primaryPrepareTxID string, branchTxRes [][]string) error
	ConfirmBranchTransaction(branchPrepareTxID string, globalTxStatus int, primaryNetwork string, primaryChain string, primaryConfirmTxID string, primaryConfirmTxProof string) error
}

type LockManager interface {
	PutStateWithLock(key string, value []byte, primaryPrepareTxId *pb.TransactionID) error
	GetStateWithLock(key string) ([]byte, error)
}
