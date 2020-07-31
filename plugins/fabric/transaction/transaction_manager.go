package transaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/zhigui-projects/sidemesh"
	"github.com/zhigui-projects/sidemesh/plugins/fabric/lock"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/zhigui-projects/sidemesh/pb"
)

type GlobalTransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetGlobalTransactionManager() sidemesh.GlobalTransactionManager
	GetLockManager() sidemesh.LockManager
}

type GlobalTransactionContext struct {
	contractapi.TransactionContext
	once        sync.Once
	globalTxMgr sidemesh.GlobalTransactionManager
	lockMgr     sidemesh.LockManager
}

func (gtxc *GlobalTransactionContext) GetGlobalTransactionManager() sidemesh.GlobalTransactionManager {
	if gtxc.globalTxMgr == nil {
		gtxc.once.Do(func() {
			gtxc.globalTxMgr = &GlobalTransactionManagerImpl{stub: gtxc.GetStub(), clientIdentity: gtxc.GetClientIdentity()}
			gtxc.lockMgr = &lock.ManagerImpl{Stub: gtxc.GetStub(), ClientIdentity: gtxc.GetClientIdentity()}
		})
	}

	return gtxc.globalTxMgr
}

func (gtxc *GlobalTransactionContext) GetLockManager() sidemesh.LockManager {
	if gtxc.lockMgr == nil {
		gtxc.once.Do(func() {
			gtxc.globalTxMgr = &GlobalTransactionManagerImpl{stub: gtxc.GetStub(), clientIdentity: gtxc.GetClientIdentity()}
			gtxc.lockMgr = &lock.ManagerImpl{Stub: gtxc.GetStub(), ClientIdentity: gtxc.GetClientIdentity()}
		})
	}

	return gtxc.lockMgr
}

type GlobalTransactionManagerImpl struct {
	stub               shim.ChaincodeStubInterface
	clientIdentity     cid.ClientIdentity
	globalTransactions map[string]*pb.GlobalTransaction
}

func (gtxm *GlobalTransactionManagerImpl) StartGlobalTransaction(ttlHeight uint64, ttlTime *timestamp.Timestamp) error {
	network, err := gtxm.stub.GetState("NetworkID")
	if err != nil {
		return err
	}
	xidKey := string(network) + gtxm.stub.GetChannelID() + gtxm.stub.GetTxID()
	xid := &pb.TransactionID{Uri: &pb.URI{Network: string(network), Chain: gtxm.stub.GetChannelID()}, Id: gtxm.stub.GetTxID()}
	var gtxTtlHeight uint64
	if ttlHeight <= 0 {
		chainInfoResponse := gtxm.stub.InvokeChaincode("qscc", [][]byte{[]byte("GetChainInfo"), []byte(xid.Uri.Chain)}, xid.Uri.Chain)
		chainInfo := &common.BlockchainInfo{}
		err = proto.Unmarshal(chainInfoResponse.Payload, chainInfo)
		if err != nil {
			return err
		}
		gtxTtlHeight = chainInfo.Height + 5
	} else {
		gtxTtlHeight = ttlHeight
	}

	var gtxTtlTime *timestamp.Timestamp
	if ttlTime == nil {
		gtxTtlTime, err = ptypes.TimestampProto(time.Now().UTC().Add(30 * time.Second))
		if err != nil {
			return err
		}
	} else {
		gtxTtlTime = ttlTime
	}

	crossChainTx := &pb.GlobalTransaction{PrimaryPrepareTxId: xid, TtlHeight: gtxTtlHeight, TtlTime: gtxTtlTime}
	gtxm.globalTransactions[xidKey] = crossChainTx

	return nil
}

func (gtxm *GlobalTransactionManagerImpl) RegisterBranchTransaction(xaBranchTx *sidemesh.XABranchTransaction) error {
	network, err := gtxm.stub.GetState("NetworkID")
	if err != nil {
		return err
	}
	xidKey := string(network) + gtxm.stub.GetChannelID() + gtxm.stub.GetTxID()

	crossChainTx := gtxm.globalTransactions[xidKey]

	branchPrepareTx := &pb.BranchTransaction{TxId: &pb.TransactionID{Uri: &pb.URI{Network: xaBranchTx.PrepareTx.Network, Chain: xaBranchTx.PrepareTx.Chain}}, Invocation: &pb.Invocation{Contract: xaBranchTx.PrepareTx.Contract, Func: xaBranchTx.PrepareTx.Function, Args: xaBranchTx.PrepareTx.Args}}
	var branchCommitTx *pb.BranchTransaction
	var branchRollbackTx *pb.BranchTransaction
	switch xaBranchTx.ChainType {
	case "fabric":
		//TODO: AUTO GENERATE COMMIT TX AND ROLLBACK TX
		branchCommitTx = &pb.BranchTransaction{}
		branchRollbackTx = &pb.BranchTransaction{}
	case "xuperchain":
		//TODO: AUTO GENERATE COMMIT TX AND ROLLBACK TX
		branchCommitTx = &pb.BranchTransaction{}
		branchRollbackTx = &pb.BranchTransaction{}
	default:
		branchCommitTx = &pb.BranchTransaction{TxId: &pb.TransactionID{Uri: &pb.URI{Network: xaBranchTx.CommitTx.Network, Chain: xaBranchTx.CommitTx.Chain}}, Invocation: &pb.Invocation{Contract: xaBranchTx.CommitTx.Contract, Func: xaBranchTx.CommitTx.Function, Args: xaBranchTx.CommitTx.Args}}
		branchRollbackTx = &pb.BranchTransaction{TxId: &pb.TransactionID{Uri: &pb.URI{Network: xaBranchTx.RollbackTx.Network, Chain: xaBranchTx.RollbackTx.Chain}}, Invocation: &pb.Invocation{Contract: xaBranchTx.RollbackTx.Contract, Func: xaBranchTx.RollbackTx.Function, Args: xaBranchTx.RollbackTx.Args}}
	}

	crossChainTx.BranchPrepareTxs = append(crossChainTx.BranchPrepareTxs, branchPrepareTx)
	crossChainTx.BranchCommitTxs = append(crossChainTx.BranchCommitTxs, branchCommitTx)
	crossChainTx.BranchRollbackTxs = append(crossChainTx.BranchRollbackTxs, branchRollbackTx)

	return nil
}

func (gtxm *GlobalTransactionManagerImpl) PreparePrimaryTransaction(xaBranchTx *sidemesh.XABranchTransaction) error {
	network, err := gtxm.stub.GetState("NetworkID")
	if err != nil {
		return err
	}
	xidKey := string(network) + gtxm.stub.GetChannelID() + gtxm.stub.GetTxID()

	crossChainTx := gtxm.globalTransactions[xidKey]

	var branchCommitTx *pb.BranchTransaction
	var branchRollbackTx *pb.BranchTransaction
	switch xaBranchTx.ChainType {
	case "fabric":
		//TODO: AUTO GENERATE COMMIT TX AND ROLLBACK TX
		branchCommitTx = &pb.BranchTransaction{}
		branchRollbackTx = &pb.BranchTransaction{}
	case "xuperchain":
		//TODO: AUTO GENERATE COMMIT TX AND ROLLBACK TX
		branchCommitTx = &pb.BranchTransaction{}
		branchRollbackTx = &pb.BranchTransaction{}
	default:
		branchCommitTx = &pb.BranchTransaction{TxId: &pb.TransactionID{Uri: &pb.URI{Network: xaBranchTx.CommitTx.Network, Chain: xaBranchTx.CommitTx.Chain}}, Invocation: &pb.Invocation{Contract: xaBranchTx.CommitTx.Contract, Func: xaBranchTx.CommitTx.Function, Args: xaBranchTx.CommitTx.Args}}
		branchRollbackTx = &pb.BranchTransaction{TxId: &pb.TransactionID{Uri: &pb.URI{Network: xaBranchTx.RollbackTx.Network, Chain: xaBranchTx.RollbackTx.Chain}}, Invocation: &pb.Invocation{Contract: xaBranchTx.RollbackTx.Contract, Func: xaBranchTx.RollbackTx.Function, Args: xaBranchTx.RollbackTx.Args}}
	}

	crossChainTx.PrimaryCommitTx = branchCommitTx
	crossChainTx.PrimaryRollbackTx = branchRollbackTx

	crossChainTxStatus := &pb.GlobalTransactionStatus{PrimaryPrepareTxId: crossChainTx.PrimaryPrepareTxId, Status: pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_PREPARED}
	crossChainTxStatusBytes, err := proto.Marshal(crossChainTxStatus)
	if err != nil {
		return err
	}
	err = gtxm.stub.PutState(xidKey+":status", crossChainTxStatusBytes)
	if err != nil {
		return err
	}

	crossChainTxBytes, err := proto.Marshal(crossChainTx)
	err = gtxm.stub.PutState(xidKey, crossChainTxBytes)
	if err != nil {
		return err
	}

	primaryTxPreparedEvent := &pb.PrimaryTransactionPreparedEvent{PrimaryPrepareTxId: crossChainTx.PrimaryPrepareTxId, PrimaryCommitTx: crossChainTx.PrimaryCommitTx, PrimaryRollbackTx: crossChainTx.PrimaryRollbackTx, BranchPrepareTxs: crossChainTx.BranchPrepareTxs}
	primaryTxPreparedEventBytes, err := proto.Marshal(primaryTxPreparedEvent)
	if err != nil {
		return err
	}
	err = gtxm.stub.SetEvent("PRIMARY_TRANSACTION_PREPARED_EVENT", primaryTxPreparedEventBytes)
	if err != nil {
		return err
	}
	return nil
}

func (gtxm *GlobalTransactionManagerImpl) ConfirmPrimaryTransaction(request sidemesh.ConfirmGlobalTransactionRequest) error {
	xidKey := request.Network + request.Chain + request.PrimaryPrepareTxID
	crossChainTxStatusBytes, err := gtxm.stub.GetState(xidKey + ":status")
	if err != nil {
		return err
	}

	if crossChainTxStatusBytes == nil {
		return errors.New("not found primary prepare tx")
	}
	crossChainTxStatus := &pb.GlobalTransactionStatus{}
	err = proto.Unmarshal(crossChainTxStatusBytes, crossChainTxStatus)
	if err != nil {
		return err
	}

	if crossChainTxStatus.Status != pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_PREPARED {
		return errors.New("duplicate confirm global transaction request")
	}

	crossChainTx := &pb.GlobalTransaction{}
	crossChainTxBytes, err := gtxm.stub.GetState(xidKey)
	err = proto.Unmarshal(crossChainTxBytes, crossChainTxStatus)
	if err != nil {
		return err
	}

	if len(crossChainTx.BranchPrepareTxs) != len(request.BranchTransactionResults) {
		return errors.New("branch transaction result count not match branch prepare transaction count")
	}

	var wg sync.WaitGroup
	var numVerified uint32 = 0
	for i, branchTxRes := range request.BranchTransactionResults {
		wg.Add(1)
		go func(branchTxRes *sidemesh.BranchTransactionResult) {
			verifyInfoKey := branchTxRes.Network + branchTxRes.Chain + ":verify"
			verifyInfoBytes, err := gtxm.stub.GetState(verifyInfoKey)
			if err != nil {
				fmt.Println(err)
				wg.Done()
				return
			}

			verifyInfo := &sidemesh.VerifyInfo{}
			err = json.Unmarshal(verifyInfoBytes, verifyInfo)
			if err != nil {
				fmt.Println(err)
				wg.Done()
				return
			}

			args := [][]byte{[]byte(verifyInfo.Function), branchTxRes.Proof}
			res := gtxm.stub.InvokeChaincode(verifyInfo.Contract, args, gtxm.stub.GetChannelID())
			if res.Status != shim.OK {
				fmt.Println(res.Message)
				wg.Done()
				return
			}

			if ok, err := strconv.ParseBool(string(res.Payload)); err == nil && ok {
				atomic.AddUint32(&numVerified, 1)
				crossChainTx.BranchPrepareTxs[i].TxId = &pb.TransactionID{Uri: &pb.URI{Network: branchTxRes.Network, Chain: branchTxRes.Chain}, Id: branchTxRes.PrepareTxID}
			} else {
				fmt.Println(err)
			}
			wg.Done()
		}(branchTxRes)
	}
	wg.Wait()

	var branchConfirmTxs []*pb.BranchTransaction
	if int(numVerified) == len(crossChainTx.BranchPrepareTxs) {
		branchConfirmTxs = crossChainTx.BranchCommitTxs
		crossChainTxStatus.Status = pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED
	} else {
		branchConfirmTxs = crossChainTx.BranchRollbackTxs
		crossChainTxStatus.Status = pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_CANCELED
	}
	crossChainTxStatusBytes, err = proto.Marshal(crossChainTxStatus)
	if err != nil {
		return err
	}
	err = gtxm.stub.PutState(xidKey+":status", crossChainTxStatusBytes)
	if err != nil {
		return err
	}

	crossChainTxBytes, err = proto.Marshal(crossChainTx)
	if err != nil {
		return err
	}
	err = gtxm.stub.PutState(xidKey, crossChainTxBytes)
	if err != nil {
		return err
	}

	network, err := gtxm.stub.GetState("NetworkID")
	if err != nil {
		return err
	}
	primaryTxConfirmedEvent := &pb.PrimaryTransactionConfirmedEvent{PrimaryConfirmTxId: &pb.TransactionID{Uri: &pb.URI{Network: string(network), Chain: gtxm.stub.GetChannelID()}, Id: gtxm.stub.GetTxID()}, BranchCommitOrRollbackTxs: branchConfirmTxs}
	primaryTxConfirmedEventBytes, err := proto.Marshal(primaryTxConfirmedEvent)
	if err != nil {
		return err
	}
	err = gtxm.stub.SetEvent("PRIMARY_TRANSACTION_CONFIRMED_EVENT", primaryTxConfirmedEventBytes)
	if err != nil {
		return err
	}
	return nil
}
