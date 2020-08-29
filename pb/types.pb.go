// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ChainType int32

const (
	ChainType_FABRIC     ChainType = 0
	ChainType_XUPERCHAIN ChainType = 1
	ChainType_BCOS       ChainType = 2
)

var ChainType_name = map[int32]string{
	0: "FABRIC",
	1: "XUPERCHAIN",
	2: "BCOS",
}

var ChainType_value = map[string]int32{
	"FABRIC":     0,
	"XUPERCHAIN": 1,
	"BCOS":       2,
}

func (x ChainType) String() string {
	return proto.EnumName(ChainType_name, int32(x))
}

func (ChainType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}

type GlobalTransactionStatusType int32

const (
	GlobalTransactionStatusType_PRIMARY_TRANSACTION_PREPARED  GlobalTransactionStatusType = 0
	GlobalTransactionStatusType_PRIMARY_TRANSACTION_COMMITTED GlobalTransactionStatusType = 1
	GlobalTransactionStatusType_PRIMARY_TRANSACTION_CANCELED  GlobalTransactionStatusType = 2
)

var GlobalTransactionStatusType_name = map[int32]string{
	0: "PRIMARY_TRANSACTION_PREPARED",
	1: "PRIMARY_TRANSACTION_COMMITTED",
	2: "PRIMARY_TRANSACTION_CANCELED",
}

var GlobalTransactionStatusType_value = map[string]int32{
	"PRIMARY_TRANSACTION_PREPARED":  0,
	"PRIMARY_TRANSACTION_COMMITTED": 1,
	"PRIMARY_TRANSACTION_CANCELED":  2,
}

func (x GlobalTransactionStatusType) String() string {
	return proto.EnumName(GlobalTransactionStatusType_name, int32(x))
}

func (GlobalTransactionStatusType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{1}
}

type BranchTransactionResponse_Status int32

const (
	BranchTransactionResponse_SUCCESS BranchTransactionResponse_Status = 0
	BranchTransactionResponse_FAILED  BranchTransactionResponse_Status = 1
)

var BranchTransactionResponse_Status_name = map[int32]string{
	0: "SUCCESS",
	1: "FAILED",
}

var BranchTransactionResponse_Status_value = map[string]int32{
	"SUCCESS": 0,
	"FAILED":  1,
}

func (x BranchTransactionResponse_Status) String() string {
	return proto.EnumName(BranchTransactionResponse_Status_name, int32(x))
}

func (BranchTransactionResponse_Status) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{4, 0}
}

type Lock struct {
	PrimaryPrepareTxId   *TransactionID `protobuf:"bytes,1,opt,name=primary_prepare_tx_id,json=primaryPrepareTxId,proto3" json:"primary_prepare_tx_id,omitempty"`
	PrevState            []byte         `protobuf:"bytes,2,opt,name=prev_state,json=prevState,proto3" json:"prev_state,omitempty"`
	UpdatingState        []byte         `protobuf:"bytes,3,opt,name=updating_state,json=updatingState,proto3" json:"updating_state,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *Lock) Reset()         { *m = Lock{} }
func (m *Lock) String() string { return proto.CompactTextString(m) }
func (*Lock) ProtoMessage()    {}
func (*Lock) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}

func (m *Lock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Lock.Unmarshal(m, b)
}
func (m *Lock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Lock.Marshal(b, m, deterministic)
}
func (m *Lock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Lock.Merge(m, src)
}
func (m *Lock) XXX_Size() int {
	return xxx_messageInfo_Lock.Size(m)
}
func (m *Lock) XXX_DiscardUnknown() {
	xxx_messageInfo_Lock.DiscardUnknown(m)
}

var xxx_messageInfo_Lock proto.InternalMessageInfo

func (m *Lock) GetPrimaryPrepareTxId() *TransactionID {
	if m != nil {
		return m.PrimaryPrepareTxId
	}
	return nil
}

func (m *Lock) GetPrevState() []byte {
	if m != nil {
		return m.PrevState
	}
	return nil
}

func (m *Lock) GetUpdatingState() []byte {
	if m != nil {
		return m.UpdatingState
	}
	return nil
}

type GlobalTransaction struct {
	PrimaryPrepareTxId   *TransactionID       `protobuf:"bytes,1,opt,name=primary_prepare_tx_id,json=primaryPrepareTxId,proto3" json:"primary_prepare_tx_id,omitempty"`
	PrimaryConfirmTx     *BranchTransaction   `protobuf:"bytes,2,opt,name=primary_confirm_tx,json=primaryConfirmTx,proto3" json:"primary_confirm_tx,omitempty"`
	BranchPrepareTxs     []*BranchTransaction `protobuf:"bytes,3,rep,name=branch_prepare_txs,json=branchPrepareTxs,proto3" json:"branch_prepare_txs,omitempty"`
	BranchConfirmTxs     []*BranchTransaction `protobuf:"bytes,4,rep,name=branch_confirm_txs,json=branchConfirmTxs,proto3" json:"branch_confirm_txs,omitempty"`
	TtlHeight            uint64               `protobuf:"varint,5,opt,name=ttl_height,json=ttlHeight,proto3" json:"ttl_height,omitempty"`
	TtlTime              *timestamp.Timestamp `protobuf:"bytes,6,opt,name=ttl_time,json=ttlTime,proto3" json:"ttl_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *GlobalTransaction) Reset()         { *m = GlobalTransaction{} }
func (m *GlobalTransaction) String() string { return proto.CompactTextString(m) }
func (*GlobalTransaction) ProtoMessage()    {}
func (*GlobalTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{1}
}

func (m *GlobalTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GlobalTransaction.Unmarshal(m, b)
}
func (m *GlobalTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GlobalTransaction.Marshal(b, m, deterministic)
}
func (m *GlobalTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GlobalTransaction.Merge(m, src)
}
func (m *GlobalTransaction) XXX_Size() int {
	return xxx_messageInfo_GlobalTransaction.Size(m)
}
func (m *GlobalTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_GlobalTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_GlobalTransaction proto.InternalMessageInfo

func (m *GlobalTransaction) GetPrimaryPrepareTxId() *TransactionID {
	if m != nil {
		return m.PrimaryPrepareTxId
	}
	return nil
}

func (m *GlobalTransaction) GetPrimaryConfirmTx() *BranchTransaction {
	if m != nil {
		return m.PrimaryConfirmTx
	}
	return nil
}

func (m *GlobalTransaction) GetBranchPrepareTxs() []*BranchTransaction {
	if m != nil {
		return m.BranchPrepareTxs
	}
	return nil
}

func (m *GlobalTransaction) GetBranchConfirmTxs() []*BranchTransaction {
	if m != nil {
		return m.BranchConfirmTxs
	}
	return nil
}

func (m *GlobalTransaction) GetTtlHeight() uint64 {
	if m != nil {
		return m.TtlHeight
	}
	return 0
}

func (m *GlobalTransaction) GetTtlTime() *timestamp.Timestamp {
	if m != nil {
		return m.TtlTime
	}
	return nil
}

type Invocation struct {
	Contract             string   `protobuf:"bytes,1,opt,name=contract,proto3" json:"contract,omitempty"`
	Func                 string   `protobuf:"bytes,2,opt,name=func,proto3" json:"func,omitempty"`
	Args                 []string `protobuf:"bytes,3,rep,name=args,proto3" json:"args,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Invocation) Reset()         { *m = Invocation{} }
func (m *Invocation) String() string { return proto.CompactTextString(m) }
func (*Invocation) ProtoMessage()    {}
func (*Invocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{2}
}

func (m *Invocation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Invocation.Unmarshal(m, b)
}
func (m *Invocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Invocation.Marshal(b, m, deterministic)
}
func (m *Invocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Invocation.Merge(m, src)
}
func (m *Invocation) XXX_Size() int {
	return xxx_messageInfo_Invocation.Size(m)
}
func (m *Invocation) XXX_DiscardUnknown() {
	xxx_messageInfo_Invocation.DiscardUnknown(m)
}

var xxx_messageInfo_Invocation proto.InternalMessageInfo

func (m *Invocation) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func (m *Invocation) GetFunc() string {
	if m != nil {
		return m.Func
	}
	return ""
}

func (m *Invocation) GetArgs() []string {
	if m != nil {
		return m.Args
	}
	return nil
}

type BranchTransaction struct {
	TxId                 *TransactionID `protobuf:"bytes,1,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	Invocation           *Invocation    `protobuf:"bytes,2,opt,name=invocation,proto3" json:"invocation,omitempty"`
	Proof                string         `protobuf:"bytes,3,opt,name=proof,proto3" json:"proof,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *BranchTransaction) Reset()         { *m = BranchTransaction{} }
func (m *BranchTransaction) String() string { return proto.CompactTextString(m) }
func (*BranchTransaction) ProtoMessage()    {}
func (*BranchTransaction) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{3}
}

func (m *BranchTransaction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BranchTransaction.Unmarshal(m, b)
}
func (m *BranchTransaction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BranchTransaction.Marshal(b, m, deterministic)
}
func (m *BranchTransaction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BranchTransaction.Merge(m, src)
}
func (m *BranchTransaction) XXX_Size() int {
	return xxx_messageInfo_BranchTransaction.Size(m)
}
func (m *BranchTransaction) XXX_DiscardUnknown() {
	xxx_messageInfo_BranchTransaction.DiscardUnknown(m)
}

var xxx_messageInfo_BranchTransaction proto.InternalMessageInfo

func (m *BranchTransaction) GetTxId() *TransactionID {
	if m != nil {
		return m.TxId
	}
	return nil
}

func (m *BranchTransaction) GetInvocation() *Invocation {
	if m != nil {
		return m.Invocation
	}
	return nil
}

func (m *BranchTransaction) GetProof() string {
	if m != nil {
		return m.Proof
	}
	return ""
}

type BranchTransactionResponse struct {
	TxId                 *TransactionID                   `protobuf:"bytes,1,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	Status               BranchTransactionResponse_Status `protobuf:"varint,2,opt,name=status,proto3,enum=pb.BranchTransactionResponse_Status" json:"status,omitempty"`
	Proof                []byte                           `protobuf:"bytes,3,opt,name=proof,proto3" json:"proof,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *BranchTransactionResponse) Reset()         { *m = BranchTransactionResponse{} }
func (m *BranchTransactionResponse) String() string { return proto.CompactTextString(m) }
func (*BranchTransactionResponse) ProtoMessage()    {}
func (*BranchTransactionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{4}
}

func (m *BranchTransactionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BranchTransactionResponse.Unmarshal(m, b)
}
func (m *BranchTransactionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BranchTransactionResponse.Marshal(b, m, deterministic)
}
func (m *BranchTransactionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BranchTransactionResponse.Merge(m, src)
}
func (m *BranchTransactionResponse) XXX_Size() int {
	return xxx_messageInfo_BranchTransactionResponse.Size(m)
}
func (m *BranchTransactionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BranchTransactionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BranchTransactionResponse proto.InternalMessageInfo

func (m *BranchTransactionResponse) GetTxId() *TransactionID {
	if m != nil {
		return m.TxId
	}
	return nil
}

func (m *BranchTransactionResponse) GetStatus() BranchTransactionResponse_Status {
	if m != nil {
		return m.Status
	}
	return BranchTransactionResponse_SUCCESS
}

func (m *BranchTransactionResponse) GetProof() []byte {
	if m != nil {
		return m.Proof
	}
	return nil
}

type PrimaryLockMeta struct {
	Xid                  string   `protobuf:"bytes,1,opt,name=xid,proto3" json:"xid,omitempty"`
	Chain                string   `protobuf:"bytes,2,opt,name=chain,proto3" json:"chain,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PrimaryLockMeta) Reset()         { *m = PrimaryLockMeta{} }
func (m *PrimaryLockMeta) String() string { return proto.CompactTextString(m) }
func (*PrimaryLockMeta) ProtoMessage()    {}
func (*PrimaryLockMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{5}
}

func (m *PrimaryLockMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrimaryLockMeta.Unmarshal(m, b)
}
func (m *PrimaryLockMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrimaryLockMeta.Marshal(b, m, deterministic)
}
func (m *PrimaryLockMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrimaryLockMeta.Merge(m, src)
}
func (m *PrimaryLockMeta) XXX_Size() int {
	return xxx_messageInfo_PrimaryLockMeta.Size(m)
}
func (m *PrimaryLockMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_PrimaryLockMeta.DiscardUnknown(m)
}

var xxx_messageInfo_PrimaryLockMeta proto.InternalMessageInfo

func (m *PrimaryLockMeta) GetXid() string {
	if m != nil {
		return m.Xid
	}
	return ""
}

func (m *PrimaryLockMeta) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

type URI struct {
	Network              string   `protobuf:"bytes,1,opt,name=network,proto3" json:"network,omitempty"`
	Chain                string   `protobuf:"bytes,2,opt,name=chain,proto3" json:"chain,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *URI) Reset()         { *m = URI{} }
func (m *URI) String() string { return proto.CompactTextString(m) }
func (*URI) ProtoMessage()    {}
func (*URI) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{6}
}

func (m *URI) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_URI.Unmarshal(m, b)
}
func (m *URI) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_URI.Marshal(b, m, deterministic)
}
func (m *URI) XXX_Merge(src proto.Message) {
	xxx_messageInfo_URI.Merge(m, src)
}
func (m *URI) XXX_Size() int {
	return xxx_messageInfo_URI.Size(m)
}
func (m *URI) XXX_DiscardUnknown() {
	xxx_messageInfo_URI.DiscardUnknown(m)
}

var xxx_messageInfo_URI proto.InternalMessageInfo

func (m *URI) GetNetwork() string {
	if m != nil {
		return m.Network
	}
	return ""
}

func (m *URI) GetChain() string {
	if m != nil {
		return m.Chain
	}
	return ""
}

type TransactionID struct {
	Uri                  *URI     `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
	Id                   string   `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TransactionID) Reset()         { *m = TransactionID{} }
func (m *TransactionID) String() string { return proto.CompactTextString(m) }
func (*TransactionID) ProtoMessage()    {}
func (*TransactionID) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{7}
}

func (m *TransactionID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TransactionID.Unmarshal(m, b)
}
func (m *TransactionID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TransactionID.Marshal(b, m, deterministic)
}
func (m *TransactionID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TransactionID.Merge(m, src)
}
func (m *TransactionID) XXX_Size() int {
	return xxx_messageInfo_TransactionID.Size(m)
}
func (m *TransactionID) XXX_DiscardUnknown() {
	xxx_messageInfo_TransactionID.DiscardUnknown(m)
}

var xxx_messageInfo_TransactionID proto.InternalMessageInfo

func (m *TransactionID) GetUri() *URI {
	if m != nil {
		return m.Uri
	}
	return nil
}

func (m *TransactionID) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type PrimaryTransactionPreparedEvent struct {
	PrimaryPrepareTxId   *TransactionID       `protobuf:"bytes,1,opt,name=primary_prepare_tx_id,json=primaryPrepareTxId,proto3" json:"primary_prepare_tx_id,omitempty"`
	PrimaryConfirmTx     *BranchTransaction   `protobuf:"bytes,2,opt,name=primary_confirm_tx,json=primaryConfirmTx,proto3" json:"primary_confirm_tx,omitempty"`
	GlobalTxStatusQuery  *Invocation          `protobuf:"bytes,3,opt,name=global_tx_status_query,json=globalTxStatusQuery,proto3" json:"global_tx_status_query,omitempty"`
	BranchPrepareTxs     []*BranchTransaction `protobuf:"bytes,4,rep,name=branch_prepare_txs,json=branchPrepareTxs,proto3" json:"branch_prepare_txs,omitempty"`
	BranchConfirmTxs     []*BranchTransaction `protobuf:"bytes,5,rep,name=branch_confirm_txs,json=branchConfirmTxs,proto3" json:"branch_confirm_txs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *PrimaryTransactionPreparedEvent) Reset()         { *m = PrimaryTransactionPreparedEvent{} }
func (m *PrimaryTransactionPreparedEvent) String() string { return proto.CompactTextString(m) }
func (*PrimaryTransactionPreparedEvent) ProtoMessage()    {}
func (*PrimaryTransactionPreparedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{8}
}

func (m *PrimaryTransactionPreparedEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrimaryTransactionPreparedEvent.Unmarshal(m, b)
}
func (m *PrimaryTransactionPreparedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrimaryTransactionPreparedEvent.Marshal(b, m, deterministic)
}
func (m *PrimaryTransactionPreparedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrimaryTransactionPreparedEvent.Merge(m, src)
}
func (m *PrimaryTransactionPreparedEvent) XXX_Size() int {
	return xxx_messageInfo_PrimaryTransactionPreparedEvent.Size(m)
}
func (m *PrimaryTransactionPreparedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_PrimaryTransactionPreparedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_PrimaryTransactionPreparedEvent proto.InternalMessageInfo

func (m *PrimaryTransactionPreparedEvent) GetPrimaryPrepareTxId() *TransactionID {
	if m != nil {
		return m.PrimaryPrepareTxId
	}
	return nil
}

func (m *PrimaryTransactionPreparedEvent) GetPrimaryConfirmTx() *BranchTransaction {
	if m != nil {
		return m.PrimaryConfirmTx
	}
	return nil
}

func (m *PrimaryTransactionPreparedEvent) GetGlobalTxStatusQuery() *Invocation {
	if m != nil {
		return m.GlobalTxStatusQuery
	}
	return nil
}

func (m *PrimaryTransactionPreparedEvent) GetBranchPrepareTxs() []*BranchTransaction {
	if m != nil {
		return m.BranchPrepareTxs
	}
	return nil
}

func (m *PrimaryTransactionPreparedEvent) GetBranchConfirmTxs() []*BranchTransaction {
	if m != nil {
		return m.BranchConfirmTxs
	}
	return nil
}

type PrimaryTransactionConfirmedEvent struct {
	PrimaryConfirmTxId   *TransactionID       `protobuf:"bytes,1,opt,name=primary_confirm_tx_id,json=primaryConfirmTxId,proto3" json:"primary_confirm_tx_id,omitempty"`
	BranchConfirmTxs     []*BranchTransaction `protobuf:"bytes,2,rep,name=branch_confirm_txs,json=branchConfirmTxs,proto3" json:"branch_confirm_txs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *PrimaryTransactionConfirmedEvent) Reset()         { *m = PrimaryTransactionConfirmedEvent{} }
func (m *PrimaryTransactionConfirmedEvent) String() string { return proto.CompactTextString(m) }
func (*PrimaryTransactionConfirmedEvent) ProtoMessage()    {}
func (*PrimaryTransactionConfirmedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{9}
}

func (m *PrimaryTransactionConfirmedEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PrimaryTransactionConfirmedEvent.Unmarshal(m, b)
}
func (m *PrimaryTransactionConfirmedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PrimaryTransactionConfirmedEvent.Marshal(b, m, deterministic)
}
func (m *PrimaryTransactionConfirmedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PrimaryTransactionConfirmedEvent.Merge(m, src)
}
func (m *PrimaryTransactionConfirmedEvent) XXX_Size() int {
	return xxx_messageInfo_PrimaryTransactionConfirmedEvent.Size(m)
}
func (m *PrimaryTransactionConfirmedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_PrimaryTransactionConfirmedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_PrimaryTransactionConfirmedEvent proto.InternalMessageInfo

func (m *PrimaryTransactionConfirmedEvent) GetPrimaryConfirmTxId() *TransactionID {
	if m != nil {
		return m.PrimaryConfirmTxId
	}
	return nil
}

func (m *PrimaryTransactionConfirmedEvent) GetBranchConfirmTxs() []*BranchTransaction {
	if m != nil {
		return m.BranchConfirmTxs
	}
	return nil
}

type BranchTransactionPreparedEvent struct {
	PrimaryPrepareTxId   *TransactionID `protobuf:"bytes,1,opt,name=primary_prepare_tx_id,json=primaryPrepareTxId,proto3" json:"primary_prepare_tx_id,omitempty"`
	GlobalTxStatusQuery  *Invocation    `protobuf:"bytes,2,opt,name=global_tx_status_query,json=globalTxStatusQuery,proto3" json:"global_tx_status_query,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *BranchTransactionPreparedEvent) Reset()         { *m = BranchTransactionPreparedEvent{} }
func (m *BranchTransactionPreparedEvent) String() string { return proto.CompactTextString(m) }
func (*BranchTransactionPreparedEvent) ProtoMessage()    {}
func (*BranchTransactionPreparedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{10}
}

func (m *BranchTransactionPreparedEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BranchTransactionPreparedEvent.Unmarshal(m, b)
}
func (m *BranchTransactionPreparedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BranchTransactionPreparedEvent.Marshal(b, m, deterministic)
}
func (m *BranchTransactionPreparedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BranchTransactionPreparedEvent.Merge(m, src)
}
func (m *BranchTransactionPreparedEvent) XXX_Size() int {
	return xxx_messageInfo_BranchTransactionPreparedEvent.Size(m)
}
func (m *BranchTransactionPreparedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_BranchTransactionPreparedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_BranchTransactionPreparedEvent proto.InternalMessageInfo

func (m *BranchTransactionPreparedEvent) GetPrimaryPrepareTxId() *TransactionID {
	if m != nil {
		return m.PrimaryPrepareTxId
	}
	return nil
}

func (m *BranchTransactionPreparedEvent) GetGlobalTxStatusQuery() *Invocation {
	if m != nil {
		return m.GlobalTxStatusQuery
	}
	return nil
}

type GlobalTransactionStatus struct {
	PrimaryPrepareTxId   *TransactionID              `protobuf:"bytes,1,opt,name=primary_prepare_tx_id,json=primaryPrepareTxId,proto3" json:"primary_prepare_tx_id,omitempty"`
	Status               GlobalTransactionStatusType `protobuf:"varint,2,opt,name=status,proto3,enum=pb.GlobalTransactionStatusType" json:"status,omitempty"`
	PrimaryConfirmTxId   *TransactionID              `protobuf:"bytes,3,opt,name=primary_confirm_tx_id,json=primaryConfirmTxId,proto3" json:"primary_confirm_tx_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *GlobalTransactionStatus) Reset()         { *m = GlobalTransactionStatus{} }
func (m *GlobalTransactionStatus) String() string { return proto.CompactTextString(m) }
func (*GlobalTransactionStatus) ProtoMessage()    {}
func (*GlobalTransactionStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{11}
}

func (m *GlobalTransactionStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GlobalTransactionStatus.Unmarshal(m, b)
}
func (m *GlobalTransactionStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GlobalTransactionStatus.Marshal(b, m, deterministic)
}
func (m *GlobalTransactionStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GlobalTransactionStatus.Merge(m, src)
}
func (m *GlobalTransactionStatus) XXX_Size() int {
	return xxx_messageInfo_GlobalTransactionStatus.Size(m)
}
func (m *GlobalTransactionStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_GlobalTransactionStatus.DiscardUnknown(m)
}

var xxx_messageInfo_GlobalTransactionStatus proto.InternalMessageInfo

func (m *GlobalTransactionStatus) GetPrimaryPrepareTxId() *TransactionID {
	if m != nil {
		return m.PrimaryPrepareTxId
	}
	return nil
}

func (m *GlobalTransactionStatus) GetStatus() GlobalTransactionStatusType {
	if m != nil {
		return m.Status
	}
	return GlobalTransactionStatusType_PRIMARY_TRANSACTION_PREPARED
}

func (m *GlobalTransactionStatus) GetPrimaryConfirmTxId() *TransactionID {
	if m != nil {
		return m.PrimaryConfirmTxId
	}
	return nil
}

type VerifyInfo struct {
	Contract             string   `protobuf:"bytes,1,opt,name=Contract,proto3" json:"Contract,omitempty"`
	Function             string   `protobuf:"bytes,2,opt,name=function,proto3" json:"function,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VerifyInfo) Reset()         { *m = VerifyInfo{} }
func (m *VerifyInfo) String() string { return proto.CompactTextString(m) }
func (*VerifyInfo) ProtoMessage()    {}
func (*VerifyInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{12}
}

func (m *VerifyInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VerifyInfo.Unmarshal(m, b)
}
func (m *VerifyInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VerifyInfo.Marshal(b, m, deterministic)
}
func (m *VerifyInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VerifyInfo.Merge(m, src)
}
func (m *VerifyInfo) XXX_Size() int {
	return xxx_messageInfo_VerifyInfo.Size(m)
}
func (m *VerifyInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_VerifyInfo.DiscardUnknown(m)
}

var xxx_messageInfo_VerifyInfo proto.InternalMessageInfo

func (m *VerifyInfo) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func (m *VerifyInfo) GetFunction() string {
	if m != nil {
		return m.Function
	}
	return ""
}

type ResourceRegisteredOrUpdatedEvent struct {
	Uri                  *URI      `protobuf:"bytes,1,opt,name=uri,proto3" json:"uri,omitempty"`
	Connection           []byte    `protobuf:"bytes,2,opt,name=connection,proto3" json:"connection,omitempty"`
	Type                 ChainType `protobuf:"varint,3,opt,name=type,proto3,enum=pb.ChainType" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *ResourceRegisteredOrUpdatedEvent) Reset()         { *m = ResourceRegisteredOrUpdatedEvent{} }
func (m *ResourceRegisteredOrUpdatedEvent) String() string { return proto.CompactTextString(m) }
func (*ResourceRegisteredOrUpdatedEvent) ProtoMessage()    {}
func (*ResourceRegisteredOrUpdatedEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{13}
}

func (m *ResourceRegisteredOrUpdatedEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ResourceRegisteredOrUpdatedEvent.Unmarshal(m, b)
}
func (m *ResourceRegisteredOrUpdatedEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ResourceRegisteredOrUpdatedEvent.Marshal(b, m, deterministic)
}
func (m *ResourceRegisteredOrUpdatedEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ResourceRegisteredOrUpdatedEvent.Merge(m, src)
}
func (m *ResourceRegisteredOrUpdatedEvent) XXX_Size() int {
	return xxx_messageInfo_ResourceRegisteredOrUpdatedEvent.Size(m)
}
func (m *ResourceRegisteredOrUpdatedEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_ResourceRegisteredOrUpdatedEvent.DiscardUnknown(m)
}

var xxx_messageInfo_ResourceRegisteredOrUpdatedEvent proto.InternalMessageInfo

func (m *ResourceRegisteredOrUpdatedEvent) GetUri() *URI {
	if m != nil {
		return m.Uri
	}
	return nil
}

func (m *ResourceRegisteredOrUpdatedEvent) GetConnection() []byte {
	if m != nil {
		return m.Connection
	}
	return nil
}

func (m *ResourceRegisteredOrUpdatedEvent) GetType() ChainType {
	if m != nil {
		return m.Type
	}
	return ChainType_FABRIC
}

func init() {
	proto.RegisterEnum("pb.ChainType", ChainType_name, ChainType_value)
	proto.RegisterEnum("pb.GlobalTransactionStatusType", GlobalTransactionStatusType_name, GlobalTransactionStatusType_value)
	proto.RegisterEnum("pb.BranchTransactionResponse_Status", BranchTransactionResponse_Status_name, BranchTransactionResponse_Status_value)
	proto.RegisterType((*Lock)(nil), "pb.Lock")
	proto.RegisterType((*GlobalTransaction)(nil), "pb.GlobalTransaction")
	proto.RegisterType((*Invocation)(nil), "pb.Invocation")
	proto.RegisterType((*BranchTransaction)(nil), "pb.BranchTransaction")
	proto.RegisterType((*BranchTransactionResponse)(nil), "pb.BranchTransactionResponse")
	proto.RegisterType((*PrimaryLockMeta)(nil), "pb.PrimaryLockMeta")
	proto.RegisterType((*URI)(nil), "pb.URI")
	proto.RegisterType((*TransactionID)(nil), "pb.TransactionID")
	proto.RegisterType((*PrimaryTransactionPreparedEvent)(nil), "pb.PrimaryTransactionPreparedEvent")
	proto.RegisterType((*PrimaryTransactionConfirmedEvent)(nil), "pb.PrimaryTransactionConfirmedEvent")
	proto.RegisterType((*BranchTransactionPreparedEvent)(nil), "pb.BranchTransactionPreparedEvent")
	proto.RegisterType((*GlobalTransactionStatus)(nil), "pb.GlobalTransactionStatus")
	proto.RegisterType((*VerifyInfo)(nil), "pb.VerifyInfo")
	proto.RegisterType((*ResourceRegisteredOrUpdatedEvent)(nil), "pb.ResourceRegisteredOrUpdatedEvent")
}

func init() { proto.RegisterFile("types.proto", fileDescriptor_d938547f84707355) }

var fileDescriptor_d938547f84707355 = []byte{
	// 911 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xcc, 0x56, 0xd1, 0x6e, 0xe3, 0x44,
	0x14, 0xad, 0x9d, 0x34, 0xad, 0x6f, 0xdb, 0x90, 0x0e, 0x2c, 0x64, 0xcb, 0xee, 0x36, 0xb5, 0x00,
	0x55, 0xfb, 0xe0, 0x15, 0x45, 0x2b, 0x04, 0xe2, 0x25, 0x75, 0x02, 0x6b, 0x69, 0xdb, 0x86, 0x89,
	0x83, 0xe0, 0xc9, 0x72, 0x9c, 0x49, 0x62, 0x6d, 0xe2, 0x31, 0x33, 0x93, 0x92, 0xf0, 0x04, 0x12,
	0x5f, 0xc0, 0x6f, 0x20, 0x7e, 0x80, 0x8f, 0xe1, 0x47, 0x90, 0x10, 0x9a, 0xf1, 0xd8, 0xcd, 0x6e,
	0x9a, 0xd2, 0x4a, 0x45, 0xe2, 0x29, 0xbe, 0xd7, 0xe7, 0x5e, 0x9f, 0xb9, 0xe7, 0xfa, 0x38, 0xb0,
	0x23, 0x16, 0x29, 0xe1, 0x4e, 0xca, 0xa8, 0xa0, 0xc8, 0x4c, 0xfb, 0x07, 0x87, 0x23, 0x4a, 0x47,
	0x13, 0xf2, 0x4c, 0x65, 0xfa, 0xb3, 0xe1, 0x33, 0x11, 0x4f, 0x09, 0x17, 0xe1, 0x34, 0xcd, 0x40,
	0xf6, 0xaf, 0x06, 0x94, 0x5f, 0xd2, 0xe8, 0x15, 0x6a, 0xc1, 0x83, 0x94, 0xc5, 0xd3, 0x90, 0x2d,
	0x82, 0x94, 0x91, 0x34, 0x64, 0x24, 0x10, 0xf3, 0x20, 0x1e, 0xd4, 0x8d, 0x86, 0x71, 0xbc, 0x73,
	0xb2, 0xef, 0xa4, 0x7d, 0xc7, 0x67, 0x61, 0xc2, 0xc3, 0x48, 0xc4, 0x34, 0xf1, 0x5a, 0x18, 0x69,
	0x7c, 0x27, 0x83, 0xfb, 0x73, 0x6f, 0x80, 0x1e, 0x03, 0xa4, 0x8c, 0x5c, 0x06, 0x5c, 0x84, 0x82,
	0xd4, 0xcd, 0x86, 0x71, 0xbc, 0x8b, 0x2d, 0x99, 0xe9, 0xca, 0x04, 0xfa, 0x10, 0xaa, 0xb3, 0x74,
	0x10, 0x8a, 0x38, 0x19, 0x69, 0x48, 0x49, 0x41, 0xf6, 0xf2, 0xac, 0x82, 0xd9, 0x7f, 0x9b, 0xb0,
	0xff, 0xd5, 0x84, 0xf6, 0xc3, 0xc9, 0xd2, 0x13, 0xef, 0x89, 0xa1, 0x0b, 0x79, 0x36, 0x88, 0x68,
	0x32, 0x8c, 0xd9, 0x34, 0x10, 0x73, 0xc5, 0x74, 0xe7, 0xe4, 0x81, 0x6c, 0x71, 0xca, 0xc2, 0x24,
	0x1a, 0x2f, 0x35, 0xc2, 0x35, 0x5d, 0xe0, 0x66, 0x78, 0x7f, 0x2e, 0x9b, 0xf4, 0x15, 0x6c, 0x89,
	0x09, 0xaf, 0x97, 0x1a, 0xa5, 0x1b, 0x9a, 0x64, 0x05, 0x05, 0x15, 0xbe, 0xd4, 0xe4, 0x8a, 0x08,
	0xaf, 0x97, 0x6f, 0xd1, 0xa4, 0x20, 0xc2, 0xe5, 0xc0, 0x85, 0x98, 0x04, 0x63, 0x12, 0x8f, 0xc6,
	0xa2, 0xbe, 0xd9, 0x30, 0x8e, 0xcb, 0xd8, 0x12, 0x62, 0xf2, 0x42, 0x25, 0xd0, 0x73, 0xd8, 0x96,
	0xb7, 0xa5, 0xea, 0xf5, 0x8a, 0x3a, 0xe3, 0x81, 0x93, 0xad, 0x84, 0x93, 0xaf, 0x84, 0xe3, 0xe7,
	0x2b, 0x81, 0xb7, 0x84, 0x98, 0xc8, 0xc8, 0xee, 0x00, 0x78, 0xc9, 0x25, 0x8d, 0x42, 0x35, 0xf8,
	0x03, 0xd8, 0x8e, 0x68, 0x22, 0x58, 0x18, 0x09, 0x35, 0x6b, 0x0b, 0x17, 0x31, 0x42, 0x50, 0x1e,
	0xce, 0x92, 0x48, 0x0d, 0xd0, 0xc2, 0xea, 0x5a, 0xe6, 0x42, 0x36, 0xca, 0xe6, 0x61, 0x61, 0x75,
	0x6d, 0xff, 0x6c, 0xc0, 0xfe, 0xca, 0x79, 0xd0, 0x47, 0xb0, 0xf9, 0x2f, 0x12, 0x96, 0x85, 0x14,
	0xcd, 0x01, 0x88, 0x0b, 0x3e, 0x5a, 0xac, 0xaa, 0x04, 0x5f, 0xb1, 0xc4, 0x4b, 0x08, 0xf4, 0x0e,
	0x6c, 0xa6, 0x8c, 0xd2, 0xa1, 0x5a, 0x2f, 0x0b, 0x67, 0x81, 0xfd, 0x87, 0x01, 0x0f, 0x57, 0x67,
	0x4a, 0x78, 0x4a, 0x13, 0x4e, 0x6e, 0xcd, 0xe5, 0x0b, 0xa8, 0xc8, 0xd5, 0x9d, 0x71, 0xc5, 0xa3,
	0x7a, 0xf2, 0xc1, 0xf5, 0x52, 0xe9, 0xb6, 0x4e, 0x57, 0x61, 0xb1, 0xae, 0x79, 0x9d, 0xd9, 0x6e,
	0xce, 0xec, 0x08, 0x2a, 0x19, 0x0e, 0xed, 0xc0, 0x56, 0xb7, 0xe7, 0xba, 0xed, 0x6e, 0xb7, 0xb6,
	0x81, 0x00, 0x2a, 0x5f, 0x36, 0xbd, 0x97, 0xed, 0x56, 0xcd, 0xb0, 0x3f, 0x83, 0xb7, 0x3a, 0xd9,
	0x1a, 0xca, 0xd7, 0xf5, 0x8c, 0x88, 0x10, 0xd5, 0xa0, 0x34, 0xd7, 0x7c, 0x2d, 0x2c, 0x2f, 0x65,
	0xf7, 0x68, 0x1c, 0xc6, 0x89, 0x96, 0x23, 0x0b, 0xec, 0xe7, 0x50, 0xea, 0x61, 0x0f, 0xd5, 0x61,
	0x2b, 0x21, 0xe2, 0x07, 0xca, 0x5e, 0xe9, 0x92, 0x3c, 0x5c, 0x53, 0xf6, 0x39, 0xec, 0xbd, 0x76,
	0x7e, 0xf4, 0x10, 0x4a, 0x33, 0x16, 0xeb, 0xf9, 0x6c, 0xc9, 0x63, 0xf7, 0xb0, 0x87, 0x65, 0x0e,
	0x55, 0xc1, 0x8c, 0x07, 0xba, 0xdc, 0x8c, 0x07, 0xf6, 0x5f, 0x26, 0x1c, 0x6a, 0xba, 0x4b, 0x3d,
	0xf4, 0xf2, 0x0f, 0xda, 0x97, 0x24, 0x11, 0xff, 0xaf, 0xf7, 0xf9, 0xdd, 0x91, 0xf2, 0x1b, 0xc9,
	0x20, 0x53, 0x2a, 0xf8, 0x7e, 0x46, 0xd8, 0x42, 0xc9, 0xb4, 0xba, 0x6b, 0x6f, 0x67, 0x68, 0x7f,
	0x9e, 0xa9, 0xf6, 0xb5, 0x84, 0xae, 0x31, 0x85, 0xf2, 0x7d, 0x98, 0xc2, 0xe6, 0x9d, 0x4c, 0xc1,
	0xfe, 0xdd, 0x80, 0xc6, 0xea, 0xf4, 0x35, 0xe0, 0x9a, 0xf1, 0x5f, 0x3d, 0xea, 0x56, 0xe3, 0x2f,
	0x9e, 0x94, 0x8d, 0xff, 0x1a, 0xbe, 0xe6, 0xdd, 0xf8, 0xfe, 0x66, 0xc0, 0x93, 0x15, 0xdc, 0x7f,
	0xb3, 0x2c, 0xeb, 0x74, 0x36, 0x6f, 0xad, 0xb3, 0xfd, 0xa7, 0x01, 0xef, 0xad, 0x7c, 0x9d, 0xf4,
	0xeb, 0x7b, 0x3f, 0x34, 0x3f, 0x7d, 0xc3, 0x62, 0x0e, 0x65, 0xd9, 0x9a, 0x47, 0xfa, 0x8b, 0x94,
	0x14, 0xee, 0xb2, 0x56, 0xd3, 0xd2, 0x1d, 0x34, 0xb5, 0x5b, 0x00, 0xdf, 0x10, 0x16, 0x0f, 0x17,
	0x5e, 0x32, 0xa4, 0xd2, 0xfd, 0xdd, 0x37, 0xdc, 0x3f, 0x8f, 0xe5, 0x3d, 0xe9, 0xf8, 0x85, 0x2b,
	0x5b, 0xb8, 0x88, 0xed, 0x9f, 0x0c, 0x68, 0x60, 0xc2, 0xe9, 0x8c, 0x45, 0x04, 0x93, 0x51, 0xcc,
	0x05, 0x61, 0x64, 0x70, 0xc1, 0x7a, 0xf2, 0x53, 0x9f, 0xcb, 0x7a, 0x83, 0xa5, 0x3c, 0x01, 0x88,
	0x68, 0x92, 0x90, 0xab, 0xee, 0xbb, 0x78, 0x29, 0x83, 0x8e, 0xa0, 0x2c, 0xff, 0xed, 0xa8, 0xa3,
	0x55, 0x4f, 0xf6, 0x64, 0xad, 0x2b, 0x7d, 0x4b, 0x0d, 0x44, 0xdd, 0x7a, 0xfa, 0x31, 0x58, 0x45,
	0x2a, 0x33, 0xd3, 0x53, 0xec, 0xb9, 0xb5, 0x0d, 0x54, 0x05, 0xf8, 0xb6, 0xd7, 0x69, 0x63, 0xf7,
	0x45, 0xd3, 0x3b, 0xaf, 0x19, 0x68, 0x1b, 0xca, 0xa7, 0xee, 0x45, 0xb7, 0x66, 0x3e, 0xfd, 0xc5,
	0x80, 0xf7, 0x6f, 0x98, 0x34, 0x6a, 0xc0, 0xa3, 0x0e, 0xf6, 0xce, 0x9a, 0xf8, 0xbb, 0xc0, 0xc7,
	0xcd, 0xf3, 0x6e, 0xd3, 0xf5, 0xbd, 0x8b, 0xf3, 0xa0, 0x83, 0xdb, 0x9d, 0x26, 0x6e, 0xb7, 0x6a,
	0x1b, 0xe8, 0x08, 0x1e, 0x5f, 0x87, 0x70, 0x2f, 0xce, 0xce, 0x3c, 0xdf, 0x97, 0x5e, 0xbe, 0xae,
	0x89, 0xdb, 0x3c, 0x77, 0xdb, 0xd2, 0xed, 0xcd, 0xd3, 0x47, 0x70, 0x10, 0xd1, 0xa9, 0xf3, 0xe3,
	0x38, 0x1e, 0xcd, 0x62, 0x27, 0x62, 0x94, 0xf3, 0x29, 0xe1, 0x63, 0xfd, 0xd1, 0xae, 0xa8, 0x9f,
	0x4f, 0xfe, 0x09, 0x00, 0x00, 0xff, 0xff, 0x58, 0x23, 0x7e, 0xc0, 0xef, 0x09, 0x00, 0x00,
}
