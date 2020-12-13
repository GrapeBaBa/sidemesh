package transaction

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
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
			lockImpl := &lock.ManagerImpl{Stub: gtxc.GetStub(), ClientIdentity: gtxc.GetClientIdentity(), WriteKeySet: map[string][]string{}}
			gtxc.lockMgr = lockImpl
			gtxc.globalTxMgr = &GlobalTransactionManagerImpl{stub: gtxc.GetStub(), clientIdentity: gtxc.GetClientIdentity(), globalTransactions: make(map[string]*pb.GlobalTransaction), lockMgr: lockImpl}
		})
	}

	return gtxc.globalTxMgr
}

func (gtxc *GlobalTransactionContext) GetLockManager() sidemesh.LockManager {
	if gtxc.lockMgr == nil {
		gtxc.once.Do(func() {
			lockImpl := &lock.ManagerImpl{Stub: gtxc.GetStub(), ClientIdentity: gtxc.GetClientIdentity(), WriteKeySet: map[string][]string{}}
			gtxc.lockMgr = lockImpl
			gtxc.globalTxMgr = &GlobalTransactionManagerImpl{stub: gtxc.GetStub(), clientIdentity: gtxc.GetClientIdentity(), globalTransactions: make(map[string]*pb.GlobalTransaction), lockMgr: lockImpl}
		})
	}

	return gtxc.lockMgr
}

type GlobalTransactionManagerImpl struct {
	stub               shim.ChaincodeStubInterface
	clientIdentity     cid.ClientIdentity
	lockMgr            *lock.ManagerImpl
	globalTransactions map[string]*pb.GlobalTransaction
}

func (gtxm *GlobalTransactionManagerImpl) StartGlobalTransaction(ttlHeight uint64, ttlTime string) error {
	network, err := gtxm.stub.GetState(sidemesh.Prefix + "NetworkID")
	if err != nil {
		return err
	}
	xidKey := sidemesh.Prefix + string(network) + gtxm.stub.GetChannelID() + gtxm.stub.GetTxID()
	xid := &pb.TransactionID{Uri: &pb.URI{Network: string(network), Chain: gtxm.stub.GetChannelID()}, Id: gtxm.stub.GetTxID()}
	if ttlTime == "" && ttlHeight <= 0 {
		return errors.New("ttlHeight and ttlTime at least one should be set")
	}
	var gtxTtlHeight uint64
	if ttlHeight > 0 {
		// Exist undetermined issue if peer height inconsistent
		//chainInfoResponse := gtxm.stub.InvokeChaincode("qscc", [][]byte{[]byte("GetChainInfo"), []byte(xid.Uri.Chain)}, xid.Uri.Chain)
		//chainInfo := &common.BlockchainInfo{}
		//err = proto.Unmarshal(chainInfoResponse.Payload, chainInfo)
		//if err != nil {
		//	return err
		//}
		//gtxTtlHeight = chainInfo.Height + 5
		gtxTtlHeight = ttlHeight
	}

	var gtxTtlTime *timestamp.Timestamp
	if ttlTime != "" {
		// Exist undetermined issue if peer time inconsistent
		//gtxTtlTime, err = ptypes.TimestampProto(time.Now().UTC().Add(30 * time.Second))
		//if err != nil {
		//	return err
		//}
		ttlTimestamp, err := time.Parse(time.RFC3339, ttlTime)
		if err != nil {
			return err
		}
		gtxTtlTime, err = ptypes.TimestampProto(ttlTimestamp)
		if err != nil {
			return err
		}
	}

	globalTx := &pb.GlobalTransaction{PrimaryPrepareTxId: xid, BranchPrepareTxs: make([]*pb.BranchTransaction, 0), BranchConfirmTxs: make([]*pb.BranchTransaction, 0), TtlHeight: gtxTtlHeight, TtlTime: gtxTtlTime}
	gtxm.globalTransactions[xidKey] = globalTx

	return nil
}

func (gtxm *GlobalTransactionManagerImpl) RegisterBranchTransaction(network string, chain string, contract string, function string, args []string) error {
	primaryNetwork, err := gtxm.stub.GetState(sidemesh.Prefix + "NetworkID")
	if err != nil {
		return err
	}
	xidKey := sidemesh.Prefix + string(primaryNetwork) + gtxm.stub.GetChannelID() + gtxm.stub.GetTxID()

	globalTx := gtxm.globalTransactions[xidKey]

	branchPrepareTx := &pb.BranchTransaction{TxId: &pb.TransactionID{Uri: &pb.URI{Network: network, Chain: chain}}, Invocation: &pb.Invocation{Contract: contract, Func: function, Args: args}}
	globalTx.BranchPrepareTxs = append(globalTx.BranchPrepareTxs, branchPrepareTx)

	return nil
}

func (gtxm *GlobalTransactionManagerImpl) PreparePrimaryTransaction() error {
	network, err := gtxm.stub.GetState(sidemesh.Prefix + "NetworkID")
	if err != nil {
		return err
	}
	xidKey := sidemesh.Prefix + string(network) + gtxm.stub.GetChannelID() + gtxm.stub.GetTxID()
	globalTx := gtxm.globalTransactions[xidKey]

	// Get Caller function name, this function should be invoked in cross chain prepare function, so that cross chain prepare function will be got.
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	funcName := frame.Function[strings.LastIndex(frame.Function, ".")+1:]

	primaryConfirmTx := &pb.BranchTransaction{TxId: &pb.TransactionID{Uri: &pb.URI{Network: string(network), Chain: gtxm.stub.GetChannelID()}}, Invocation: &pb.Invocation{Func: "Confirm" + funcName, Args: []string{gtxm.stub.GetTxID()}}}
	globalTx.PrimaryConfirmTx = primaryConfirmTx

	globalTxStatus := &pb.GlobalTransactionStatus{PrimaryPrepareTxId: globalTx.PrimaryPrepareTxId, Status: pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_PREPARED}
	globalTxStatusBytes, err := proto.Marshal(globalTxStatus)
	if err != nil {
		return err
	}
	err = gtxm.stub.PutState(xidKey+":status", globalTxStatusBytes)
	if err != nil {
		return err
	}

	globalTxBytes, err := proto.Marshal(globalTx)
	err = gtxm.stub.PutState(xidKey, globalTxBytes)
	if err != nil {
		return err
	}

	if writeKeySet, ok := gtxm.lockMgr.WriteKeySet[xidKey]; ok {
		var buf bytes.Buffer
		err = gob.NewEncoder(&buf).Encode(&writeKeySet)
		if err != nil {
			return err
		}
		fmt.Println(len(writeKeySet))
		fmt.Println("wset")
		fmt.Println(buf.Bytes())
		err = gtxm.stub.PutState(xidKey+":wset", buf.Bytes())
		if err != nil {
			return err
		}
	}

	queryGlobalTxInvocation := &pb.Invocation{Func: "QueryGlobalTxStatus"}
	primaryTxPreparedEvent := &pb.PrimaryTransactionPreparedEvent{PrimaryPrepareTxId: globalTx.PrimaryPrepareTxId, PrimaryConfirmTx: globalTx.PrimaryConfirmTx, BranchPrepareTxs: globalTx.BranchPrepareTxs, GlobalTxStatusQuery: queryGlobalTxInvocation}
	primaryTxPreparedEventBytes, err := proto.Marshal(primaryTxPreparedEvent)
	if err != nil {
		return err
	}
	err = gtxm.stub.SetEvent(sidemesh.Prefix+"PRIMARY_TRANSACTION_PREPARED_EVENT", primaryTxPreparedEventBytes)
	if err != nil {
		return err
	}
	return nil
}

func (gtxm *GlobalTransactionManagerImpl) StartBranchTransaction(primaryNetwork string, primaryChain string, primaryTxID string, primaryTxProof string) error {
	if primaryTxProof == "" {
		return errors.New("need primary tx proof")
	}
	args := [][]byte{[]byte("Resolve"), []byte(primaryNetwork), []byte(primaryChain)}
	res := gtxm.stub.InvokeChaincode("verifyregistry", args, gtxm.stub.GetChannelID())
	if res.Status != shim.OK {
		return errors.New(res.Message)
	}
	fmt.Println(res)
	fmt.Println(res.Payload)
	verifyInfo := &sidemesh.VerifyInfo{}
	err := json.Unmarshal(res.Payload, verifyInfo)
	if err != nil {
		return err
	}

	vargs := [][]byte{[]byte(verifyInfo.Function), []byte(primaryTxID), []byte(primaryTxProof)}
	vres := gtxm.stub.InvokeChaincode(verifyInfo.Contract, vargs, gtxm.stub.GetChannelID())
	if vres.Status != shim.OK {
		return errors.New(vres.Message)
	}

	ok, err := strconv.ParseBool(string(vres.Payload))
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("primary tx verify failed")
	}
	return nil
}

func (gtxm *GlobalTransactionManagerImpl) PrepareBranchTransaction(primaryNetwork string, primaryChain string, primaryTxID string, globalTxQueryContract string, globalTxQueryFunc string) error {
	network, err := gtxm.stub.GetState(sidemesh.Prefix + "NetworkID")
	if err != nil {
		return err
	}
	globalTxID := &pb.TransactionID{Uri: &pb.URI{Network: primaryNetwork, Chain: primaryChain}, Id: primaryTxID}
	globalTxQuery := &pb.Invocation{Contract: globalTxQueryContract, Func: globalTxQueryFunc, Args: []string{primaryTxID}}

	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	funcName := frame.Function[strings.LastIndex(frame.Function, ".")+1:]

	bidKey := sidemesh.Prefix + string(network) + gtxm.stub.GetChannelID() + gtxm.stub.GetTxID()
	if writeKeySet, ok := gtxm.lockMgr.WriteKeySet[bidKey]; ok {
		var buf bytes.Buffer
		err = gob.NewEncoder(&buf).Encode(&writeKeySet)
		if err != nil {
			return err
		}
		fmt.Println("wset")
		fmt.Println(buf.Bytes())
		err = gtxm.stub.PutState(bidKey+":wset", buf.Bytes())
		if err != nil {
			return err
		}
	}

	branchConfirmTx := &pb.BranchTransaction{TxId: &pb.TransactionID{Uri: &pb.URI{Network: string(network), Chain: gtxm.stub.GetChannelID()}}, Invocation: &pb.Invocation{Func: "Confirm" + funcName, Args: []string{gtxm.stub.GetTxID()}}}
	branchTxPreparedEvent := &pb.BranchTransactionPreparedEvent{PrimaryPrepareTxId: globalTxID, GlobalTxStatusQuery: globalTxQuery, ConfirmTx: branchConfirmTx}
	branchTxPreparedEventBytes, err := proto.Marshal(branchTxPreparedEvent)
	if err != nil {
		return err
	}
	err = gtxm.stub.SetEvent(sidemesh.Prefix+"BRANCH_TRANSACTION_PREPARED_EVENT", branchTxPreparedEventBytes)
	if err != nil {
		return err
	}
	return nil
}

func (gtxm *GlobalTransactionManagerImpl) ConfirmPrimaryTransaction(primaryPrepareTxID string, branchTxRes [][]string) error {
	network, err := gtxm.stub.GetState(sidemesh.Prefix + "NetworkID")
	if err != nil {
		return err
	}
	xidKey := sidemesh.Prefix + string(network) + gtxm.stub.GetChannelID() + primaryPrepareTxID
	globalChainTxStatusBytes, err := gtxm.stub.GetState(xidKey + ":status")
	if err != nil {
		return err
	}

	if globalChainTxStatusBytes == nil {
		return fmt.Errorf("not found primary prepare tx status by xid %s", xidKey)
	}

	globalTxStatus := &pb.GlobalTransactionStatus{}
	err = proto.Unmarshal(globalChainTxStatusBytes, globalTxStatus)
	if err != nil {
		return err
	}

	if globalTxStatus.Status != pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_PREPARED {
		return errors.New("duplicate confirm global transaction request")
	}

	globalTx := &pb.GlobalTransaction{}
	globalChainTxBytes, err := gtxm.stub.GetState(xidKey)
	err = proto.Unmarshal(globalChainTxBytes, globalTx)
	if err != nil {
		return err
	}

	if len(globalTx.BranchPrepareTxs) != len(branchTxRes) {
		return errors.New("dependent transaction result count not match branch prepare transaction count")
	}

	var numVerified uint32 = 0
	for i, depTxRes := range branchTxRes {
		globalTx.BranchPrepareTxs[i].TxId.Id = depTxRes[2]
		globalTx.BranchPrepareTxs[i].Proof = depTxRes[3]
		fmt.Println(depTxRes[0])
		fmt.Println(depTxRes[1])
		fmt.Println(depTxRes[2])
		fmt.Println(depTxRes[3])
		args := [][]byte{[]byte("Resolve"), []byte(depTxRes[0]), []byte(depTxRes[1])}
		res := gtxm.stub.InvokeChaincode("verifyregistry", args, gtxm.stub.GetChannelID())
		fmt.Println("1")
		if res.Status != shim.OK {
			fmt.Println(res.Message)
			return errors.New(res.Message)
		}

		fmt.Println("3")
		verifyInfo := &sidemesh.VerifyInfo{}
		fmt.Println(res.Payload)
		err := json.Unmarshal(res.Payload, verifyInfo)
		if err != nil {
			return err
		}

		fmt.Println("5")
		fmt.Println(verifyInfo.Function)
		fmt.Println(verifyInfo.Contract)
		if depTxRes[3] != "" {
			args := [][]byte{[]byte(verifyInfo.Function), []byte(depTxRes[2]), []byte(depTxRes[3])}
			res := gtxm.stub.InvokeChaincode(verifyInfo.Contract, args, gtxm.stub.GetChannelID())
			if res.Status != shim.OK {
				fmt.Println(res.Message)
				return errors.New(res.Message)
			}

			ok, err := strconv.ParseBool(string(res.Payload))
			if err != nil {
				fmt.Println(err)
				return err
			}
			if ok {
				numVerified++
			}
		}
		fmt.Println("6")
	}

	if int(numVerified) == len(globalTx.BranchPrepareTxs) {
		globalTxStatus.Status = pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED
	} else {
		globalTxStatus.Status = pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_CANCELED
	}

	fmt.Println("7")
	//waitingConfirmTxResponse := gtxm.stub.InvokeChaincode("qscc", [][]byte{[]byte("GetTransactionByID"), []byte(gtxm.stub.GetChannelID()), []byte(primaryPrepareTxID)}, gtxm.stub.GetChannelID())
	//if waitingConfirmTxResponse.Status != shim.OK {
	//	return errors.New(waitingConfirmTxResponse.Message)
	//}
	//
	//fmt.Println("8")
	//tx := &peer.ProcessedTransaction{}
	//err = proto.Unmarshal(waitingConfirmTxResponse.Payload, tx)
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Println("9")
	//if tx.ValidationCode != 0 {
	//	return errors.New("invalid tx validation code")
	//}
	//
	//fmt.Println("10")
	//chaincodeAction, err := protoutil.GetActionFromEnvelopeMsg(tx.TransactionEnvelope)
	//if err != nil {
	//	return err
	//}
	//
	//fmt.Println("11")
	//txRWSet := &rwsetutil.TxRwSet{}
	//if err = txRWSet.FromProtoBytes(chaincodeAction.Results); err != nil {
	//	return err
	//}

	fmt.Println("12")
	wsetBytes, err := gtxm.stub.GetState(xidKey + ":wset")
	if err != nil {
		return err
	}

	if wsetBytes != nil {
		fmt.Println("15")
		var wset []string
		err := gob.NewDecoder(bytes.NewBuffer(wsetBytes)).Decode(&wset)
		if err != nil {
			return err
		}
		fmt.Println("16")
		for _, wKey := range wset {
			l := &pb.Lock{}
			v, err := gtxm.stub.GetState(wKey)
			err = proto.Unmarshal(v, l)
			if err != nil {
				return fmt.Errorf("wkey %s lock not exist", wKey)
			}
			// TODO: check lock if expired
			if globalTxStatus.Status == pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED {
				err = gtxm.stub.PutState(wKey, l.UpdatingState)
				if err != nil {
					return err
				}
			} else if globalTxStatus.Status == pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_CANCELED {
				err = gtxm.stub.PutState(wKey, l.PrevState)
				if err != nil {
					return err
				}
			}
		}
	}
	//for _, ns := range txRWSet.NsRwSets {
	//	if ns.KvRwSet != nil && len(ns.KvRwSet.Writes) > 0 {
	//		for _, write := range ns.KvRwSet.Writes {
	//			l := &pb.Lock{}
	//			err = proto.Unmarshal(write.Value, l)
	//			if err != nil {
	//				continue
	//			}
	//			// TODO: check lock if expired
	//			if int(numVerified) == len(globalTx.BranchPrepareTxs) {
	//				globalTxStatus.Status = pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED
	//				err = gtxm.stub.PutState(write.Key, l.UpdatingState)
	//				if err != nil {
	//					return err
	//				}
	//			} else {
	//				globalTxStatus.Status = pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_CANCELED
	//				err = gtxm.stub.PutState(write.Key, l.PrevState)
	//				if err != nil {
	//					return err
	//				}
	//			}
	//
	//		}
	//	}

	// cannot implement private data and metadata lock
	//for _, c := range ns.CollHashedRwSets {
	//	if c.HashedRwSet != nil && len(c.HashedRwSet.HashedWrites) > 0 {
	//		for _, write := range c.HashedRwSet.HashedWrites {
	//			lock := &pb.Lock{}
	//			err = proto.Unmarshal(write.ValueHash, lock)
	//			if err != nil {
	//				return false, err
	//			}
	//			if crossChainTxStatus.Status == pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED {
	//				err = Stub.PutPrivateData(write, lock.UpdatingState)
	//				if err != nil {
	//					return false, err
	//				}
	//			} else {
	//				err = Stub.PutState(write.Key, lock.PrevState)
	//				if err != nil {
	//					return false, err
	//				}
	//			}
	//		}
	//	}
	//
	//	// private metadata updates
	//	if c.HashedRwSet != nil && len(c.HashedRwSet.MetadataWrites) > 0 {
	//		return true
	//	}
	//}

	//if ns.KvRwSet != nil && len(ns.KvRwSet.MetadataWrites) > 0 {
	//	for _, write := range ns.KvRwSet.MetadataWrites {
	//		lock := &pb.Lock{}
	//		err = proto.Unmarshal(write.Entries, lock)
	//		if err != nil {
	//			return false, err
	//		}
	//		if crossChainTxStatus.Status == pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED {
	//			err = Stub.PutState(write.Key, lock.UpdatingState)
	//			if err != nil {
	//				return false, err
	//			}
	//		} else {
	//			err = Stub.PutState(write.Key, lock.PrevState)
	//			if err != nil {
	//				return false, err
	//			}
	//		}
	//	}
	//}
	//}

	fmt.Println("13")
	for _, branchPrepareTx := range globalTx.BranchPrepareTxs {
		branchConfirmTx := &pb.BranchTransaction{TxId: &pb.TransactionID{Uri: &pb.URI{Network: branchPrepareTx.TxId.Uri.Network, Chain: branchPrepareTx.TxId.Uri.Chain}}, Invocation: &pb.Invocation{Contract: branchPrepareTx.Invocation.Contract, Func: "Confirm" + branchPrepareTx.Invocation.Func, Args: []string{branchPrepareTx.TxId.Id, strconv.Itoa(int(globalTxStatus.Status)), string(network), gtxm.stub.GetChannelID(), gtxm.stub.GetTxID()}}}
		globalTx.BranchConfirmTxs = append(globalTx.BranchConfirmTxs, branchConfirmTx)
	}

	fmt.Println("14")
	globalTxStatusBytes, err := proto.Marshal(globalTxStatus)
	if err != nil {
		return err
	}
	err = gtxm.stub.PutState(xidKey+":status", globalTxStatusBytes)
	if err != nil {
		return err
	}

	globalTxBytes, err := proto.Marshal(globalTx)
	if err != nil {
		return err
	}
	err = gtxm.stub.PutState(xidKey, globalTxBytes)
	if err != nil {
		return err
	}

	primaryTxConfirmedEvent := &pb.PrimaryTransactionConfirmedEvent{PrimaryConfirmTxId: &pb.TransactionID{Uri: &pb.URI{Network: string(network), Chain: gtxm.stub.GetChannelID()}, Id: gtxm.stub.GetTxID()}, BranchConfirmTxs: globalTx.BranchConfirmTxs}
	primaryTxConfirmedEventBytes, err := proto.Marshal(primaryTxConfirmedEvent)
	if err != nil {
		return err
	}
	err = gtxm.stub.SetEvent(sidemesh.Prefix+"PRIMARY_TRANSACTION_CONFIRMED_EVENT", primaryTxConfirmedEventBytes)
	if err != nil {
		return err
	}

	return nil
}

func (gtxm *GlobalTransactionManagerImpl) ConfirmBranchTransaction(branchPrepareTxID string, globalTxStatus int, primaryNetwork string, primaryChain string, primaryConfirmTxID string, primaryConfirmTxProof string) error {
	args := [][]byte{[]byte("Resolve"), []byte(primaryNetwork), []byte(primaryChain)}
	res := gtxm.stub.InvokeChaincode("verifyregistry", args, gtxm.stub.GetChannelID())
	if res.Status != shim.OK {
		return errors.New(res.Message)
	}

	verifyInfo := &sidemesh.VerifyInfo{}
	err := json.Unmarshal(res.Payload, verifyInfo)
	if err != nil {
		return err
	}

	if primaryConfirmTxProof == "" {
		return errors.New("primary confirm tx proof cannot empty")
	}

	vargs := [][]byte{[]byte(verifyInfo.Function), []byte(primaryConfirmTxID), []byte(primaryConfirmTxProof)}
	vres := gtxm.stub.InvokeChaincode(verifyInfo.Contract, vargs, gtxm.stub.GetChannelID())
	if vres.Status != shim.OK {
		return errors.New(vres.Message)
	}

	ok, err := strconv.ParseBool(string(vres.Payload))
	if err != nil {
		return err
	}

	if !ok {
		return errors.New("primary confirm tx verify failed")
	}

	//waitingConfirmTxResponse := gtxm.stub.InvokeChaincode("qscc", [][]byte{[]byte("GetTransactionByID"), []byte(gtxm.stub.GetChannelID()), []byte(branchPrepareTxID)}, gtxm.stub.GetChannelID())
	//if waitingConfirmTxResponse.Status != shim.OK {
	//	return errors.New(waitingConfirmTxResponse.Message)
	//}
	//
	//tx := &peer.ProcessedTransaction{}
	//err = proto.Unmarshal(waitingConfirmTxResponse.Payload, tx)
	//if err != nil {
	//	return err
	//}
	//
	//if tx.ValidationCode != 0 {
	//	return errors.New("invalid tx validation code")
	//}
	//
	//chaincodeAction, err := protoutil.GetActionFromEnvelopeMsg(tx.TransactionEnvelope)
	//if err != nil {
	//	return err
	//}
	//
	//txRWSet := &rwsetutil.TxRwSet{}
	//if err = txRWSet.FromProtoBytes(chaincodeAction.Results); err != nil {
	//	return err
	//}

	network, err := gtxm.stub.GetState(sidemesh.Prefix + "NetworkID")
	if err != nil {
		return err
	}
	bidKey := sidemesh.Prefix + string(network) + gtxm.stub.GetChannelID() + branchPrepareTxID
	wsetBytes, err := gtxm.stub.GetState(bidKey + ":wset")
	if err != nil {
		return err
	}

	if wsetBytes != nil {
		fmt.Println("13")
		var wset []string
		err := gob.NewDecoder(bytes.NewBuffer(wsetBytes)).Decode(&wset)
		if err != nil {
			return err
		}
		fmt.Println("14")
		for _, wKey := range wset {
			l := &pb.Lock{}
			v, err := gtxm.stub.GetState(wKey)
			err = proto.Unmarshal(v, l)
			if err != nil {
				return fmt.Errorf("wkey %s lock not exist", wKey)
			}
			// TODO: check lock if expired
			if globalTxStatus == int(pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED) {
				err = gtxm.stub.PutState(wKey, l.UpdatingState)
				if err != nil {
					return err
				}
			} else if globalTxStatus == int(pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_CANCELED) {
				err = gtxm.stub.PutState(wKey, l.PrevState)
				if err != nil {
					return err
				}
			}
		}
	}

	//for _, ns := range txRWSet.NsRwSets {
	//	if ns.KvRwSet != nil && len(ns.KvRwSet.Writes) > 0 {
	//		for _, write := range ns.KvRwSet.Writes {
	//			l := &pb.Lock{}
	//			err = proto.Unmarshal(write.Value, l)
	//			if err != nil {
	//				continue
	//			}
	//			if globalTxStatus == int(pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED) {
	//				err = gtxm.stub.PutState(write.Key, l.UpdatingState)
	//				if err != nil {
	//					return err
	//				}
	//			} else {
	//				err = gtxm.stub.PutState(write.Key, l.PrevState)
	//				if err != nil {
	//					return err
	//				}
	//			}
	//
	//		}
	//	}

	// cannot implement private data and metadata lock
	//for _, c := range ns.CollHashedRwSets {
	//	if c.HashedRwSet != nil && len(c.HashedRwSet.HashedWrites) > 0 {
	//		for _, write := range c.HashedRwSet.HashedWrites {
	//			lock := &pb.Lock{}
	//			err = proto.Unmarshal(write.ValueHash, lock)
	//			if err != nil {
	//				return false, err
	//			}
	//			if crossChainTxStatus.Status == pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED {
	//				err = Stub.PutPrivateData(write, lock.UpdatingState)
	//				if err != nil {
	//					return false, err
	//				}
	//			} else {
	//				err = Stub.PutState(write.Key, lock.PrevState)
	//				if err != nil {
	//					return false, err
	//				}
	//			}
	//		}
	//	}
	//
	//	// private metadata updates
	//	if c.HashedRwSet != nil && len(c.HashedRwSet.MetadataWrites) > 0 {
	//		return true
	//	}
	//}

	//if ns.KvRwSet != nil && len(ns.KvRwSet.MetadataWrites) > 0 {
	//	for _, write := range ns.KvRwSet.MetadataWrites {
	//		lock := &pb.Lock{}
	//		err = proto.Unmarshal(write.Entries, lock)
	//		if err != nil {
	//			return false, err
	//		}
	//		if crossChainTxStatus.Status == pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED {
	//			err = Stub.PutState(write.Key, lock.UpdatingState)
	//			if err != nil {
	//				return false, err
	//			}
	//		} else {
	//			err = Stub.PutState(write.Key, lock.PrevState)
	//			if err != nil {
	//				return false, err
	//			}
	//		}
	//	}
	//}
	//}

	return nil
}
