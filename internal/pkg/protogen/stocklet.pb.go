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
// source: stocklet/stocklet.proto

// buf:lint:ignore PACKAGE_VERSION_SUFFIX

package protogen

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_stocklet_stocklet_proto protoreflect.FileDescriptor

var file_stocklet_stocklet_proto_rawDesc = []byte{
	0x0a, 0x17, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x6c, 0x65, 0x74, 0x2f, 0x73, 0x74, 0x6f, 0x63, 0x6b,
	0x6c, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x73, 0x74, 0x6f, 0x63, 0x6b,
	0x6c, 0x65, 0x74, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d,
	0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70, 0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e,
	0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x42, 0xd6, 0x01, 0x92, 0x41, 0x9f, 0x01, 0x12, 0x8e, 0x01, 0x0a, 0x08, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x6c, 0x65, 0x74, 0x22, 0x38, 0x0a, 0x11, 0x47, 0x69, 0x74, 0x48, 0x75,
	0x62, 0x20, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x12, 0x23, 0x68, 0x74,
	0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x68, 0x65, 0x78, 0x6f, 0x6c, 0x61, 0x6e, 0x2f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x6c, 0x65,
	0x74, 0x2a, 0x41, 0x0a, 0x08, 0x41, 0x47, 0x50, 0x4c, 0x2d, 0x33, 0x2e, 0x30, 0x12, 0x35, 0x68,
	0x74, 0x74, 0x70, 0x73, 0x3a, 0x2f, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x68, 0x65, 0x78, 0x6f, 0x6c, 0x61, 0x6e, 0x2f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x6c,
	0x65, 0x74, 0x2f, 0x62, 0x6c, 0x6f, 0x62, 0x2f, 0x6d, 0x61, 0x69, 0x6e, 0x2f, 0x4c, 0x49, 0x43,
	0x45, 0x4e, 0x53, 0x45, 0x32, 0x05, 0x30, 0x2e, 0x31, 0x2e, 0x30, 0x1a, 0x09, 0x6c, 0x6f, 0x63,
	0x61, 0x6c, 0x68, 0x6f, 0x73, 0x74, 0x2a, 0x01, 0x01, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x68, 0x65, 0x78, 0x6f, 0x6c, 0x61, 0x6e, 0x2f, 0x73, 0x74,
	0x6f, 0x63, 0x6b, 0x6c, 0x65, 0x74, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var file_stocklet_stocklet_proto_goTypes = []interface{}{}
var file_stocklet_stocklet_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_stocklet_stocklet_proto_init() }
func file_stocklet_stocklet_proto_init() {
	if File_stocklet_stocklet_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_stocklet_stocklet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_stocklet_stocklet_proto_goTypes,
		DependencyIndexes: file_stocklet_stocklet_proto_depIdxs,
	}.Build()
	File_stocklet_stocklet_proto = out.File
	file_stocklet_stocklet_proto_rawDesc = nil
	file_stocklet_stocklet_proto_goTypes = nil
	file_stocklet_stocklet_proto_depIdxs = nil
}
