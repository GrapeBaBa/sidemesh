package lock

import (
	"errors"
	"fmt"
	"github.com/zhigui-projects/sidemesh"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/common"
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

		timeout, err := checkTimeoutLock(maybeLock, lockManager.Stub)
		if err != nil {
			return nil, err
		}

		if timeout {
			//TODO:IF need more info to notify caller to release the lock
			return nil, &NeedReleaseExpiredLockError{}
		} else {
			time.Sleep(2 * time.Second)
			continue
		}
	}
}

func (lockManager *ManagerImpl) PutStateWithPrimaryLock(key string, value []byte) error {
	network, err := lockManager.Stub.GetState(sidemesh.Prefix + "NetworkID")
	if err != nil {
		return err
	}
	primaryPrepareTxId := &pb.TransactionID{Uri: &pb.URI{Network: string(network), Chain: lockManager.Stub.GetChannelID()}, Id: lockManager.Stub.GetTxID()}
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

		timeout, err := checkTimeoutLock(maybeLock, lockManager.Stub)
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

func (lockManager *ManagerImpl) PutStateWithBranchLock(key string, value []byte, primaryNetwork string, primaryChain string, primaryTxID string) error {
	primaryPrepareTxId := &pb.TransactionID{Uri: &pb.URI{Network: primaryNetwork, Chain: primaryChain}, Id: primaryTxID}
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

		timeout, err := checkTimeoutLock(maybeLock, lockManager.Stub)
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

func checkTimeoutLock(lock *pb.Lock, stub shim.ChaincodeStubInterface) (bool, error) {
	xid := lock.PrimaryPrepareTxId

	network, err := stub.GetState(sidemesh.Prefix + "NetworkID")
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

	xidKey := sidemesh.Prefix + xid.Uri.Network + xid.Uri.Chain + xid.Id
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
		return false, err
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

	// This may be not consistent in multiple peers endorse phase
	chainInfoResponse := stub.InvokeChaincode("qscc", [][]byte{[]byte("GetChainInfo"), []byte(xid.Uri.Chain)}, xid.Uri.Chain)
	chainInfo := &common.BlockchainInfo{}
	err = proto.Unmarshal(chainInfoResponse.Payload, chainInfo)
	if err != nil {
		return false, err
	}
	if chainInfo.Height > crossChainTx.TtlHeight {
		return true, nil
	}

	// This may be not consistent in multiple peers endorse phase
	now := time.Now().UTC()
	ttlTime, err := ptypes.Timestamp(crossChainTx.TtlTime)
	if err != nil {
		return false, err
	}
	if now.After(ttlTime) {
		return true, nil
	}

	return false, nil
}
