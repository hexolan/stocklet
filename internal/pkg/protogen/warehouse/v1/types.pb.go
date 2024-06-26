// Copyright (C) 2024 Declan Teevan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: stocklet/warehouse/v1/types.proto

package warehouse_v1

import (
	_ "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
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

type ProductStock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId string `protobuf:"bytes,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Quantity  int32  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (x *ProductStock) Reset() {
	*x = ProductStock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stocklet_warehouse_v1_types_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProductStock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductStock) ProtoMessage() {}

func (x *ProductStock) ProtoReflect() protoreflect.Message {
	mi := &file_stocklet_warehouse_v1_types_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductStock.ProtoReflect.Descriptor instead.
func (*ProductStock) Descriptor() ([]byte, []int) {
	return file_stocklet_warehouse_v1_types_proto_rawDescGZIP(), []int{0}
}

func (x *ProductStock) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *ProductStock) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type Reservation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id            string              `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	OrderId       string              `protobuf:"bytes,2,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	ReservedStock []*ReservationStock `protobuf:"bytes,3,rep,name=reserved_stock,json=reservedStock,proto3" json:"reserved_stock,omitempty"`
	CreatedAt     int64               `protobuf:"varint,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
}

func (x *Reservation) Reset() {
	*x = Reservation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stocklet_warehouse_v1_types_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reservation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reservation) ProtoMessage() {}

func (x *Reservation) ProtoReflect() protoreflect.Message {
	mi := &file_stocklet_warehouse_v1_types_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reservation.ProtoReflect.Descriptor instead.
func (*Reservation) Descriptor() ([]byte, []int) {
	return file_stocklet_warehouse_v1_types_proto_rawDescGZIP(), []int{1}
}

func (x *Reservation) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Reservation) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

func (x *Reservation) GetReservedStock() []*ReservationStock {
	if x != nil {
		return x.ReservedStock
	}
	return nil
}

func (x *Reservation) GetCreatedAt() int64 {
	if x != nil {
		return x.CreatedAt
	}
	return 0
}

type ReservationStock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProductId string `protobuf:"bytes,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Quantity  int32  `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (x *ReservationStock) Reset() {
	*x = ReservationStock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stocklet_warehouse_v1_types_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReservationStock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReservationStock) ProtoMessage() {}

func (x *ReservationStock) ProtoReflect() protoreflect.Message {
	mi := &file_stocklet_warehouse_v1_types_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReservationStock.ProtoReflect.Descriptor instead.
func (*ReservationStock) Descriptor() ([]byte, []int) {
	return file_stocklet_warehouse_v1_types_proto_rawDescGZIP(), []int{2}
}

func (x *ReservationStock) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *ReservationStock) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

var File_stocklet_warehouse_v1_types_proto protoreflect.FileDescriptor

var file_stocklet_warehouse_v1_types_proto_rawDesc = []byte{
	0x0a, 0x21, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x6c, 0x65, 0x74, 0x2f, 0x77, 0x61, 0x72, 0x65, 0x68,
	0x6f, 0x75, 0x73, 0x65, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x15, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x6c, 0x65, 0x74, 0x2e, 0x77, 0x61,
	0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x62, 0x75, 0x66, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5b, 0x0a, 0x0c, 0x50, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x26, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75,
	0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12,
	0x23, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x42, 0x07, 0xba, 0x48, 0x04, 0x1a, 0x02, 0x28, 0x00, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x22, 0xb9, 0x01, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x22, 0x0a,
	0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42,
	0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x4e, 0x0a, 0x0e, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x5f, 0x73, 0x74,
	0x6f, 0x63, 0x6b, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x73, 0x74, 0x6f, 0x63,
	0x6b, 0x6c, 0x65, 0x74, 0x2e, 0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x2e, 0x76,
	0x31, 0x2e, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x6f,
	0x63, 0x6b, 0x52, 0x0d, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x53, 0x74, 0x6f, 0x63,
	0x6b, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74,
	0x22, 0x5f, 0x0a, 0x10, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x12, 0x26, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xba, 0x48, 0x04, 0x72, 0x02, 0x10,
	0x01, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x23, 0x0a, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x42, 0x07,
	0xba, 0x48, 0x04, 0x1a, 0x02, 0x28, 0x00, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x42, 0x4d, 0x5a, 0x4b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x68, 0x65, 0x78, 0x6f, 0x6c, 0x61, 0x6e, 0x2f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x6c, 0x65, 0x74,
	0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2f, 0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65,
	0x2f, 0x76, 0x31, 0x3b, 0x77, 0x61, 0x72, 0x65, 0x68, 0x6f, 0x75, 0x73, 0x65, 0x5f, 0x76, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stocklet_warehouse_v1_types_proto_rawDescOnce sync.Once
	file_stocklet_warehouse_v1_types_proto_rawDescData = file_stocklet_warehouse_v1_types_proto_rawDesc
)

func file_stocklet_warehouse_v1_types_proto_rawDescGZIP() []byte {
	file_stocklet_warehouse_v1_types_proto_rawDescOnce.Do(func() {
		file_stocklet_warehouse_v1_types_proto_rawDescData = protoimpl.X.CompressGZIP(file_stocklet_warehouse_v1_types_proto_rawDescData)
	})
	return file_stocklet_warehouse_v1_types_proto_rawDescData
}

var file_stocklet_warehouse_v1_types_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_stocklet_warehouse_v1_types_proto_goTypes = []interface{}{
	(*ProductStock)(nil),     // 0: stocklet.warehouse.v1.ProductStock
	(*Reservation)(nil),      // 1: stocklet.warehouse.v1.Reservation
	(*ReservationStock)(nil), // 2: stocklet.warehouse.v1.ReservationStock
}
var file_stocklet_warehouse_v1_types_proto_depIdxs = []int32{
	2, // 0: stocklet.warehouse.v1.Reservation.reserved_stock:type_name -> stocklet.warehouse.v1.ReservationStock
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_stocklet_warehouse_v1_types_proto_init() }
func file_stocklet_warehouse_v1_types_proto_init() {
	if File_stocklet_warehouse_v1_types_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_stocklet_warehouse_v1_types_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProductStock); i {
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
		file_stocklet_warehouse_v1_types_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reservation); i {
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
		file_stocklet_warehouse_v1_types_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReservationStock); i {
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
			RawDescriptor: file_stocklet_warehouse_v1_types_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_stocklet_warehouse_v1_types_proto_goTypes,
		DependencyIndexes: file_stocklet_warehouse_v1_types_proto_depIdxs,
		MessageInfos:      file_stocklet_warehouse_v1_types_proto_msgTypes,
	}.Build()
	File_stocklet_warehouse_v1_types_proto = out.File
	file_stocklet_warehouse_v1_types_proto_rawDesc = nil
	file_stocklet_warehouse_v1_types_proto_goTypes = nil
	file_stocklet_warehouse_v1_types_proto_depIdxs = nil
}
