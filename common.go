package sidemesh

const Prefix string = "SIDE_MESH_"

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
	PutLockedStateWithPrimaryLock(key string, value []byte) error
	PutLockedStateWithBranchLock(key string, value []byte, primaryNetwork string, primaryChain string, primaryTxID string) error
	GetStateMaybeLocked(key string) ([]byte, error)
	PutStateMaybeLocked(key string, value []byte) error
}
