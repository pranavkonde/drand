//
// This protobuf file contains the services and message definitions of all
// methods used by drand nodes to produce distributed randomness.
//

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.3
// source: drand/protocol.proto

package drand

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type IdentityRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *Metadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *IdentityRequest) Reset() {
	*x = IdentityRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_drand_protocol_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdentityRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentityRequest) ProtoMessage() {}

func (x *IdentityRequest) ProtoReflect() protoreflect.Message {
	mi := &file_drand_protocol_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdentityRequest.ProtoReflect.Descriptor instead.
func (*IdentityRequest) Descriptor() ([]byte, []int) {
	return file_drand_protocol_proto_rawDescGZIP(), []int{0}
}

func (x *IdentityRequest) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type IdentityResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Key     []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	// Deprecated: Marked as deprecated in drand/protocol.proto.
	Tls bool `protobuf:"varint,3,opt,name=tls,proto3" json:"tls,omitempty"`
	// BLS signature over the identity to prove possession of the private key
	Signature []byte `protobuf:"bytes,4,opt,name=signature,proto3" json:"signature,omitempty"`
	// --------------
	Metadata *Metadata `protobuf:"bytes,5,opt,name=metadata,proto3" json:"metadata,omitempty"`
	// We need to specify the scheme name to make sure the key is getting probably decoded on the client side
	SchemeName string `protobuf:"bytes,6,opt,name=schemeName,proto3" json:"schemeName,omitempty"`
}

func (x *IdentityResponse) Reset() {
	*x = IdentityResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_drand_protocol_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IdentityResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IdentityResponse) ProtoMessage() {}

func (x *IdentityResponse) ProtoReflect() protoreflect.Message {
	mi := &file_drand_protocol_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IdentityResponse.ProtoReflect.Descriptor instead.
func (*IdentityResponse) Descriptor() ([]byte, []int) {
	return file_drand_protocol_proto_rawDescGZIP(), []int{1}
}

func (x *IdentityResponse) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *IdentityResponse) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

// Deprecated: Marked as deprecated in drand/protocol.proto.
func (x *IdentityResponse) GetTls() bool {
	if x != nil {
		return x.Tls
	}
	return false
}

func (x *IdentityResponse) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *IdentityResponse) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *IdentityResponse) GetSchemeName() string {
	if x != nil {
		return x.SchemeName
	}
	return ""
}

type PartialBeaconPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Round is the round for which the beacon will be created from the partial
	// signatures
	Round uint64 `protobuf:"varint,1,opt,name=round,proto3" json:"round,omitempty"`
	// signature of the previous round - could be removed at some point but now
	// is used to verify the signature even before accessing the store
	PreviousSignature []byte `protobuf:"bytes,2,opt,name=previous_signature,json=previousSignature,proto3" json:"previous_signature,omitempty"`
	// partial signature - a threshold of them needs to be aggregated to produce
	// the final beacon at the given round.
	PartialSig []byte `protobuf:"bytes,3,opt,name=partial_sig,json=partialSig,proto3" json:"partial_sig,omitempty"`
	Metadata *Metadata `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *PartialBeaconPacket) Reset() {
	*x = PartialBeaconPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_drand_protocol_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartialBeaconPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialBeaconPacket) ProtoMessage() {}

func (x *PartialBeaconPacket) ProtoReflect() protoreflect.Message {
	mi := &file_drand_protocol_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialBeaconPacket.ProtoReflect.Descriptor instead.
func (*PartialBeaconPacket) Descriptor() ([]byte, []int) {
	return file_drand_protocol_proto_rawDescGZIP(), []int{2}
}

func (x *PartialBeaconPacket) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *PartialBeaconPacket) GetPreviousSignature() []byte {
	if x != nil {
		return x.PreviousSignature
	}
	return nil
}

func (x *PartialBeaconPacket) GetPartialSig() []byte {
	if x != nil {
		return x.PartialSig
	}
	return nil
}

func (x *PartialBeaconPacket) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

// SyncRequest is from a node that needs to sync up with the current head of the
// chain
type SyncRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromRound uint64 `protobuf:"varint,1,opt,name=from_round,json=fromRound,proto3" json:"from_round,omitempty"`
	Metadata *Metadata `protobuf:"bytes,2,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *SyncRequest) Reset() {
	*x = SyncRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_drand_protocol_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncRequest) ProtoMessage() {}

func (x *SyncRequest) ProtoReflect() protoreflect.Message {
	mi := &file_drand_protocol_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncRequest.ProtoReflect.Descriptor instead.
func (*SyncRequest) Descriptor() ([]byte, []int) {
	return file_drand_protocol_proto_rawDescGZIP(), []int{3}
}

func (x *SyncRequest) GetFromRound() uint64 {
	if x != nil {
		return x.FromRound
	}
	return 0
}

func (x *SyncRequest) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type BeaconPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PreviousSignature []byte    `protobuf:"bytes,1,opt,name=previous_signature,json=previousSignature,proto3" json:"previous_signature,omitempty"`
	Round             uint64    `protobuf:"varint,2,opt,name=round,proto3" json:"round,omitempty"`
	Signature         []byte    `protobuf:"bytes,3,opt,name=signature,proto3" json:"signature,omitempty"`
	Metadata          *Metadata `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`
}

func (x *BeaconPacket) Reset() {
	*x = BeaconPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_drand_protocol_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BeaconPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BeaconPacket) ProtoMessage() {}

func (x *BeaconPacket) ProtoReflect() protoreflect.Message {
	mi := &file_drand_protocol_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BeaconPacket.ProtoReflect.Descriptor instead.
func (*BeaconPacket) Descriptor() ([]byte, []int) {
	return file_drand_protocol_proto_rawDescGZIP(), []int{4}
}

func (x *BeaconPacket) GetPreviousSignature() []byte {
	if x != nil {
		return x.PreviousSignature
	}
	return nil
}

func (x *BeaconPacket) GetRound() uint64 {
	if x != nil {
		return x.Round
	}
	return 0
}

func (x *BeaconPacket) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

func (x *BeaconPacket) GetMetadata() *Metadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

var File_drand_protocol_proto protoreflect.FileDescriptor

var file_drand_protocol_proto_rawDesc = []byte{
	0x0a, 0x14, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x1a, 0x12, 0x64,
	0x72, 0x61, 0x6e, 0x64, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x3e, 0x0a, 0x0f, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x2b, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0x22, 0xbf, 0x01, 0x0a, 0x10, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x14, 0x0a, 0x03, 0x74, 0x6c, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x42,
	0x02, 0x18, 0x01, 0x52, 0x03, 0x74, 0x6c, 0x73, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x2b, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64,
	0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x4e, 0x61, 0x6d,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0xa8, 0x01, 0x0a, 0x13, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x42,
	0x65, 0x61, 0x63, 0x6f, 0x6e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x72,
	0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e,
	0x64, 0x12, 0x2d, 0x0a, 0x12, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x73, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x70,
	0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x12, 0x1f, 0x0a, 0x0b, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x5f, 0x73, 0x69, 0x67, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x70, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x53, 0x69,
	0x67, 0x12, 0x2b, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x59,
	0x0a, 0x0b, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x09, 0x66, 0x72, 0x6f, 0x6d, 0x52, 0x6f, 0x75, 0x6e, 0x64, 0x12, 0x2b, 0x0a, 0x08,
	0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f,
	0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x22, 0x9e, 0x01, 0x0a, 0x0c, 0x42, 0x65,
	0x61, 0x63, 0x6f, 0x6e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12, 0x2d, 0x0a, 0x12, 0x70, 0x72,
	0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x11, 0x70, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73,
	0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x6f, 0x75,
	0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x2b, 0x0a,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0f, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x32, 0xf6, 0x01, 0x0a, 0x08, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x3e, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x49, 0x64,
	0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x12, 0x16, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x49,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x17,
	0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x0d, 0x50, 0x61, 0x72, 0x74, 0x69,
	0x61, 0x6c, 0x42, 0x65, 0x61, 0x63, 0x6f, 0x6e, 0x12, 0x1a, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64,
	0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x42, 0x65, 0x61, 0x63, 0x6f, 0x6e, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x1a, 0x0c, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x12, 0x36, 0x0a, 0x09, 0x53, 0x79, 0x6e, 0x63, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x12,
	0x12, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x42, 0x65, 0x61, 0x63,
	0x6f, 0x6e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x30, 0x01, 0x12, 0x37, 0x0a, 0x06, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x2e, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x64, 0x72, 0x61,
	0x6e, 0x64, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2f, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x2f, 0x76, 0x32,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x64, 0x72, 0x61, 0x6e, 0x64, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_drand_protocol_proto_rawDescOnce sync.Once
	file_drand_protocol_proto_rawDescData = file_drand_protocol_proto_rawDesc
)

func file_drand_protocol_proto_rawDescGZIP() []byte {
	file_drand_protocol_proto_rawDescOnce.Do(func() {
		file_drand_protocol_proto_rawDescData = protoimpl.X.CompressGZIP(file_drand_protocol_proto_rawDescData)
	})
	return file_drand_protocol_proto_rawDescData
}

var file_drand_protocol_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_drand_protocol_proto_goTypes = []interface{}{
	(*IdentityRequest)(nil),     // 0: drand.IdentityRequest
	(*IdentityResponse)(nil),    // 1: drand.IdentityResponse
	(*PartialBeaconPacket)(nil), // 2: drand.PartialBeaconPacket
	(*SyncRequest)(nil),         // 3: drand.SyncRequest
	(*BeaconPacket)(nil),        // 4: drand.BeaconPacket
	(*Metadata)(nil),            // 5: drand.Metadata
	(*StatusRequest)(nil),       // 6: drand.StatusRequest
	(*Empty)(nil),               // 7: drand.Empty
	(*StatusResponse)(nil),      // 8: drand.StatusResponse
}
var file_drand_protocol_proto_depIdxs = []int32{
	5, // 0: drand.IdentityRequest.metadata:type_name -> drand.Metadata
	5, // 1: drand.IdentityResponse.metadata:type_name -> drand.Metadata
	5, // 2: drand.PartialBeaconPacket.metadata:type_name -> drand.Metadata
	5, // 3: drand.SyncRequest.metadata:type_name -> drand.Metadata
	5, // 4: drand.BeaconPacket.metadata:type_name -> drand.Metadata
	0, // 5: drand.Protocol.GetIdentity:input_type -> drand.IdentityRequest
	2, // 6: drand.Protocol.PartialBeacon:input_type -> drand.PartialBeaconPacket
	3, // 7: drand.Protocol.SyncChain:input_type -> drand.SyncRequest
	6, // 8: drand.Protocol.Status:input_type -> drand.StatusRequest
	1, // 9: drand.Protocol.GetIdentity:output_type -> drand.IdentityResponse
	7, // 10: drand.Protocol.PartialBeacon:output_type -> drand.Empty
	4, // 11: drand.Protocol.SyncChain:output_type -> drand.BeaconPacket
	8, // 12: drand.Protocol.Status:output_type -> drand.StatusResponse
	9, // [9:13] is the sub-list for method output_type
	5, // [5:9] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_drand_protocol_proto_init() }
func file_drand_protocol_proto_init() {
	if File_drand_protocol_proto != nil {
		return
	}
	file_drand_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_drand_protocol_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdentityRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_drand_protocol_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IdentityResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_drand_protocol_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartialBeaconPacket); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_drand_protocol_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_drand_protocol_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BeaconPacket); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_drand_protocol_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_drand_protocol_proto_goTypes,
		DependencyIndexes: file_drand_protocol_proto_depIdxs,
		MessageInfos:      file_drand_protocol_proto_msgTypes,
	}.Build()
	File_drand_protocol_proto = out.File
	file_drand_protocol_proto_rawDesc = nil
	file_drand_protocol_proto_goTypes = nil
	file_drand_protocol_proto_depIdxs = nil
}
