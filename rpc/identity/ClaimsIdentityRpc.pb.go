// Code generated by protoc-gen-go. DO NOT EDIT.
// source: identity/ClaimsIdentityRpc.proto

package identity // import "github.com/graphql-editor/azure-functions-golang-worker/rpc/identity"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import shared "github.com/graphql-editor/azure-functions-golang-worker/rpc/shared"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Light-weight representation of a .NET System.Security.Claims.ClaimsIdentity object.
// This is the same serialization as found in EasyAuth, and needs to be kept in sync with
// its ClaimsIdentitySlim definition, as seen in the WebJobs extension:
// https://github.com/Azure/azure-webjobs-sdk-extensions/blob/dev/src/WebJobs.Extensions.Http/ClaimsIdentitySlim.cs
type RpcClaimsIdentity struct {
	AuthenticationType   *shared.NullableString `protobuf:"bytes,1,opt,name=authentication_type,json=authenticationType,proto3" json:"authentication_type,omitempty"`
	NameClaimType        *shared.NullableString `protobuf:"bytes,2,opt,name=name_claim_type,json=nameClaimType,proto3" json:"name_claim_type,omitempty"`
	RoleClaimType        *shared.NullableString `protobuf:"bytes,3,opt,name=role_claim_type,json=roleClaimType,proto3" json:"role_claim_type,omitempty"`
	Claims               []*RpcClaim            `protobuf:"bytes,4,rep,name=claims,proto3" json:"claims,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *RpcClaimsIdentity) Reset()         { *m = RpcClaimsIdentity{} }
func (m *RpcClaimsIdentity) String() string { return proto.CompactTextString(m) }
func (*RpcClaimsIdentity) ProtoMessage()    {}
func (*RpcClaimsIdentity) Descriptor() ([]byte, []int) {
	return fileDescriptor_ClaimsIdentityRpc_98f493a5313d96ff, []int{0}
}
func (m *RpcClaimsIdentity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RpcClaimsIdentity.Unmarshal(m, b)
}
func (m *RpcClaimsIdentity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RpcClaimsIdentity.Marshal(b, m, deterministic)
}
func (dst *RpcClaimsIdentity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RpcClaimsIdentity.Merge(dst, src)
}
func (m *RpcClaimsIdentity) XXX_Size() int {
	return xxx_messageInfo_RpcClaimsIdentity.Size(m)
}
func (m *RpcClaimsIdentity) XXX_DiscardUnknown() {
	xxx_messageInfo_RpcClaimsIdentity.DiscardUnknown(m)
}

var xxx_messageInfo_RpcClaimsIdentity proto.InternalMessageInfo

func (m *RpcClaimsIdentity) GetAuthenticationType() *shared.NullableString {
	if m != nil {
		return m.AuthenticationType
	}
	return nil
}

func (m *RpcClaimsIdentity) GetNameClaimType() *shared.NullableString {
	if m != nil {
		return m.NameClaimType
	}
	return nil
}

func (m *RpcClaimsIdentity) GetRoleClaimType() *shared.NullableString {
	if m != nil {
		return m.RoleClaimType
	}
	return nil
}

func (m *RpcClaimsIdentity) GetClaims() []*RpcClaim {
	if m != nil {
		return m.Claims
	}
	return nil
}

// Light-weight representation of a .NET System.Security.Claims.Claim object.
// This is the same serialization as found in EasyAuth, and needs to be kept in sync with
// its ClaimSlim definition, as seen in the WebJobs extension:
// https://github.com/Azure/azure-webjobs-sdk-extensions/blob/dev/src/WebJobs.Extensions.Http/ClaimSlim.cs
type RpcClaim struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	Type                 string   `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RpcClaim) Reset()         { *m = RpcClaim{} }
func (m *RpcClaim) String() string { return proto.CompactTextString(m) }
func (*RpcClaim) ProtoMessage()    {}
func (*RpcClaim) Descriptor() ([]byte, []int) {
	return fileDescriptor_ClaimsIdentityRpc_98f493a5313d96ff, []int{1}
}
func (m *RpcClaim) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RpcClaim.Unmarshal(m, b)
}
func (m *RpcClaim) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RpcClaim.Marshal(b, m, deterministic)
}
func (dst *RpcClaim) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RpcClaim.Merge(dst, src)
}
func (m *RpcClaim) XXX_Size() int {
	return xxx_messageInfo_RpcClaim.Size(m)
}
func (m *RpcClaim) XXX_DiscardUnknown() {
	xxx_messageInfo_RpcClaim.DiscardUnknown(m)
}

var xxx_messageInfo_RpcClaim proto.InternalMessageInfo

func (m *RpcClaim) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func (m *RpcClaim) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func init() {
	proto.RegisterType((*RpcClaimsIdentity)(nil), "RpcClaimsIdentity")
	proto.RegisterType((*RpcClaim)(nil), "RpcClaim")
}

func init() {
	proto.RegisterFile("identity/ClaimsIdentityRpc.proto", fileDescriptor_ClaimsIdentityRpc_98f493a5313d96ff)
}

var fileDescriptor_ClaimsIdentityRpc_98f493a5313d96ff = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x91, 0xcf, 0x4a, 0x33, 0x31,
	0x14, 0xc5, 0xe9, 0xd7, 0x7e, 0xc5, 0xa6, 0x48, 0x71, 0x74, 0x51, 0xba, 0xaa, 0x5d, 0x15, 0x61,
	0x12, 0xa8, 0x82, 0x5b, 0x51, 0x37, 0x6e, 0x5c, 0x44, 0x57, 0x6e, 0x4a, 0x9a, 0xa6, 0x99, 0x60,
	0xfe, 0x79, 0x93, 0x28, 0xf5, 0x91, 0x7d, 0x0a, 0xc9, 0x4c, 0xa7, 0x52, 0xd0, 0x5d, 0x4e, 0xee,
	0xf9, 0x9d, 0xe4, 0x70, 0xd1, 0x54, 0xad, 0x85, 0x8d, 0x2a, 0x6e, 0xc9, 0x9d, 0x66, 0xca, 0x84,
	0x87, 0x9d, 0xa4, 0x9e, 0x63, 0x0f, 0x2e, 0xba, 0xc9, 0x24, 0x54, 0x0c, 0xc4, 0x9a, 0x3c, 0x26,
	0xad, 0xd9, 0x4a, 0x8b, 0xe7, 0xad, 0x17, 0xa1, 0x99, 0xcd, 0xbe, 0x3a, 0xe8, 0x84, 0x7a, 0x7e,
	0x88, 0x16, 0x37, 0xe8, 0x94, 0xa5, 0x58, 0x65, 0xc5, 0x59, 0x54, 0xce, 0x2e, 0xe3, 0xd6, 0x8b,
	0x71, 0x67, 0xda, 0x99, 0x0f, 0x17, 0x23, 0xdc, 0x06, 0x3d, 0x45, 0x50, 0x56, 0xd2, 0xe2, 0xd0,
	0x9b, 0xe3, 0x8b, 0x6b, 0x34, 0xb2, 0xcc, 0x88, 0x25, 0xcf, 0xc1, 0x0d, 0xfd, 0xef, 0x77, 0xfa,
	0x38, 0xfb, 0xea, 0xf7, 0x5b, 0x10, 0x9c, 0x3e, 0x00, 0xbb, 0x7f, 0x80, 0xd9, 0xf7, 0x03, 0x9e,
	0xa3, 0x7e, 0xcd, 0x84, 0x71, 0x6f, 0xda, 0x9d, 0x0f, 0x17, 0x03, 0xdc, 0xf6, 0xa2, 0xbb, 0xc1,
	0xec, 0x0a, 0x1d, 0xb5, 0x77, 0xc5, 0x19, 0xfa, 0xff, 0xce, 0x74, 0x6a, 0x4a, 0x0d, 0x68, 0x23,
	0x8a, 0x02, 0xf5, 0xf6, 0x7f, 0x1d, 0xd0, 0xfa, 0x7c, 0x0b, 0xe8, 0x82, 0x3b, 0x83, 0x8d, 0xe2,
	0xe0, 0x82, 0xdb, 0x44, 0xcc, 0x3e, 0x13, 0x08, 0xbc, 0x49, 0x96, 0xe7, 0xba, 0x01, 0x83, 0xe7,
	0xd8, 0x88, 0x10, 0x98, 0x14, 0xe1, 0xe5, 0x5e, 0xaa, 0x58, 0xa5, 0x15, 0xe6, 0xce, 0x10, 0x09,
	0xcc, 0x57, 0x6f, 0xba, 0x14, 0x6b, 0x15, 0x1d, 0x90, 0x9a, 0x2b, 0xf7, 0x5c, 0x29, 0x9d, 0x66,
	0x56, 0x96, 0x1f, 0x0e, 0x5e, 0x05, 0x10, 0xf0, 0x9c, 0xb4, 0xab, 0x5c, 0xf5, 0xeb, 0xed, 0x5c,
	0x7e, 0x07, 0x00, 0x00, 0xff, 0xff, 0x98, 0xeb, 0xd0, 0xa8, 0xdd, 0x01, 0x00, 0x00,
}
