package lock

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/protoutil"
	"github.com/zhigui-projects/sidemesh/pb"
)

type NeedReleaseExpiredLockError struct {
}

func (nrele *NeedReleaseExpiredLockError) Error() string {
	return "NeedReleaseExpiredLockError"
}

type ManagerImpl struct {
	Stub           shim.ChaincodeStubInterface
	ClientIdentity cid.ClientIdentity
}

func (lockManager *ManagerImpl) GetStateWithLock(key string) ([]byte, error) {
	for {
		existValue, err := lockManager.Stub.GetState(key)
		maybeLock := &pb.Lock{}
		err = proto.Unmarshal(existValue, maybeLock)
		if err != nil {
			return lockManager.Stub.GetState(key)
		}

		timeout, err := checkAndReleaseTimeoutLock(maybeLock, lockManager.Stub)
		if err != nil {
			return nil, err
		}

		if timeout {
			return nil, &NeedReleaseExpiredLockError{}
		} else {
			time.Sleep(2 * time.Second)
			continue
		}
	}
}

func (lockManager *ManagerImpl) PutStateWithLock(key string, value []byte, primaryPrepareTxId *pb.TransactionID) error {
	for {
		existValue, err := lockManager.Stub.GetState(key)
		maybeLock := &pb.Lock{}
		err = proto.Unmarshal(existValue, maybeLock)
		if err != nil {
			newLock := &pb.Lock{PrevState: existValue, UpdatingState: value, PrimaryPrepareTxId: primaryPrepareTxId}
			newLockBytes, err := proto.Marshal(newLock)
			if err != nil {
				return err
			}
			return lockManager.Stub.PutState(key, newLockBytes)
		}

		timeout, err := checkAndReleaseTimeoutLock(maybeLock, lockManager.Stub)
		if err != nil {
			return err
		}

		if timeout {
			newLock := &pb.Lock{PrevState: existValue, UpdatingState: value, PrimaryPrepareTxId: primaryPrepareTxId}
			newLockBytes, err := proto.Marshal(newLock)
			if err != nil {
				return err
			}
			return lockManager.Stub.PutState(key, newLockBytes)
		} else {
			time.Sleep(2 * time.Second)
			continue
		}
	}

}

func checkAndReleaseTimeoutLock(lock *pb.Lock, stub shim.ChaincodeStubInterface) (bool, error) {
	xid := lock.PrimaryPrepareTxId

	network, err := stub.GetState("NetworkID")
	if err != nil {
		return false, err
	}

	if string(network) != xid.Uri.Network {
		fmt.Println("this is a secondary lock")
		// TODO: current we cannot judge secondary lock, it need external mesher listening the primary prepare tx status.
		return false, nil
	}

	if stub.GetChannelID() != xid.Uri.Chain {
		fmt.Println("wrong channel")
		// TODO: current we cannot judge secondary lock, it need external mesher listening the primary prepare tx status.
		return false, nil
	}

	xidKey := xid.Uri.Network + xid.Uri.Chain + xid.Id
	crossChainTxStatusBytes, err := stub.GetState(xidKey + ":status")
	if err != nil {
		return false, err
	}

	if crossChainTxStatusBytes == nil {
		return false, errors.New("not found primary prepare tx status")
	}
	crossChainTxStatus := &pb.GlobalTransactionStatus{}
	err = proto.Unmarshal(crossChainTxStatusBytes, crossChainTxStatus)
	if err != nil {
		return false, err
	}

	if crossChainTxStatus.Status == pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED || crossChainTxStatus.Status == pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_CANCELED {
		fmt.Println("invalid tx status it should not have lock")
		err = releaseLock(stub, xid, crossChainTxStatus)
		if err != nil {
			return true, err
		}

	}

	crossChainTxBytes, err := stub.GetState(xidKey)
	if err != nil {
		return false, err
	}

	if crossChainTxBytes == nil {
		return false, errors.New("not found primary prepare tx")
	}
	crossChainTx := &pb.GlobalTransaction{}
	err = proto.Unmarshal(crossChainTxBytes, crossChainTx)
	if err != nil {
		return false, err
	}

	chainInfoResponse := stub.InvokeChaincode("qscc", [][]byte{[]byte("GetChainInfo"), []byte(xid.Uri.Chain)}, xid.Uri.Chain)
	chainInfo := &common.BlockchainInfo{}
	err = proto.Unmarshal(chainInfoResponse.Payload, chainInfo)
	if err != nil {
		return false, err
	}
	if chainInfo.Height > crossChainTx.TtlHeight {
		err = releaseLock(stub, xid, crossChainTxStatus)
		return true, err
	}

	now := time.Now().UTC()
	ttlTime, err := ptypes.Timestamp(crossChainTx.TtlTime)
	if err != nil {
		return false, err
	}
	if now.After(ttlTime) {
		err = releaseLock(stub, xid, crossChainTxStatus)
		return true, err
	}

	return false, nil
}

func releaseLock(stub shim.ChaincodeStubInterface, xid *pb.TransactionID, crossChainTxStatus *pb.GlobalTransactionStatus) error {
	txResponse := stub.InvokeChaincode("qscc", [][]byte{[]byte("GetTransactionByID"), []byte(xid.Uri.Chain), []byte(xid.Id)}, xid.Uri.Chain)
	if txResponse.Status != shim.OK {
		fmt.Println(txResponse.Message)
		return errors.New(txResponse.Message)
	}

	tx := &peer.ProcessedTransaction{}
	err := proto.Unmarshal(txResponse.Payload, tx)
	if err != nil {
		return err
	}

	if tx.ValidationCode != 0 {
		fmt.Println("invalid tx validate code it should not have lock")
		return errors.New("invalid primary prepare tx validation code")
	}

	chaincodeAction, err := protoutil.GetActionFromEnvelopeMsg(tx.TransactionEnvelope)
	if err != nil {
		return err
	}

	txRWSet := &rwsetutil.TxRwSet{}
	if err = txRWSet.FromProtoBytes(chaincodeAction.Results); err != nil {
		return err
	}

	for _, ns := range txRWSet.NsRwSets {
		if ns.KvRwSet != nil && len(ns.KvRwSet.Writes) > 0 {
			for _, write := range ns.KvRwSet.Writes {
				lock := &pb.Lock{}
				err = proto.Unmarshal(write.Value, lock)
				if err != nil {
					return err
				}
				if crossChainTxStatus.Status == pb.GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED {
					err = stub.PutState(write.Key, lock.UpdatingState)
					if err != nil {
						return err
					}
				} else {
					err = stub.PutState(write.Key, lock.PrevState)
					if err != nil {
						return err
					}
				}
			}
		}

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
	}
	return nil
}
