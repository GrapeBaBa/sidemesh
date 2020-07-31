package sidemesh

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/zhigui-projects/sidemesh/pb"
)

type XABranchTransaction struct {
	PrepareTx  *BranchTransaction `json:"prepare_tx"`
	CommitTx   *BranchTransaction `json:"commit_tx"`
	RollbackTx *BranchTransaction `json:"rollback_tx"`
	ChainType  string
}

type BranchTransaction struct {
	Network  string   `json:"network"`
	Chain    string   `json:"chain"`
	Contract string   `json:"contract"`
	Function string   `json:"function"`
	Args     [][]byte `json:"args"`
}

type BranchTransactionResult struct {
	Network     string `json:"network"`
	Chain       string `json:"chain"`
	PrepareTxID string `json:"prepare_tx_id"`

	Proof []byte `json:"proof"`
}

type ConfirmGlobalTransactionRequest struct {
	Network                  string                     `json:"network"`
	Chain                    string                     `json:"chain"`
	PrimaryPrepareTxID       string                     `json:"primary_prepare_tx_id"`
	IsCommit                 bool                       `json:"is_commit"`
	BranchTransactionResults []*BranchTransactionResult `json:"branch_transaction_results"`
}

type VerifyInfo struct {
	Contract string `json:"contract"`
	Function string `json:"function"`
}

type GlobalTransactionManager interface {
	StartGlobalTransaction(ttlHeight uint64, ttlTime *timestamp.Timestamp) error
	RegisterBranchTransaction(xaBranchTx *XABranchTransaction) error
	PreparePrimaryTransaction(xaBranchTx *XABranchTransaction) error
	ConfirmPrimaryTransaction(request ConfirmGlobalTransactionRequest) error
}

type LockManager interface {
	PutStateWithLock(key string, value []byte, primaryPrepareTxId *pb.TransactionID) error
	GetStateOrReleaseExpiredLock(key string) ([]byte, error)
}
