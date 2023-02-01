// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate sh -c "protoc --go_out=paths=source_relative:. any.proto"

package opaqueany

import (
	"strings"

	proto "google.golang.org/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoregistry "google.golang.org/protobuf/reflect/protoregistry"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

// New marshals src into a new Any instance.
func New(src proto.Message) (*Any, error) {
	dst := new(Any)
	if err := dst.MarshalFrom(src); err != nil {
		return nil, err
	}
	return dst, nil
}

// MarshalFrom marshals src into dst as the underlying message
// using the provided marshal options.
//
// If no options are specified, call dst.MarshalFrom instead.
func MarshalFrom(dst *Any, src proto.Message, opts proto.MarshalOptions) error {
	const urlPrefix = "type.googleapis.com/"
	if src == nil {
		return protoimpl.X.NewError("invalid nil source message")
	}
	b, err := opts.Marshal(src)
	if err != nil {
		return err
	}
	dst.TypeUrl = urlPrefix + string(src.ProtoReflect().Descriptor().FullName())
	dst.Value = b
	return nil
}

// UnmarshalTo unmarshals the underlying message from src into dst
// using the provided unmarshal options.
// It reports an error if dst is not of the right message type.
//
// If no options are specified, call src.UnmarshalTo instead.
func UnmarshalTo(src *Any, dst proto.Message, opts proto.UnmarshalOptions) error {
	if src == nil {
		return protoimpl.X.NewError("invalid nil source message")
	}
	if !src.MessageIs(dst) {
		got := dst.ProtoReflect().Descriptor().FullName()
		want := src.MessageName()
		return protoimpl.X.NewError("mismatched message type: got %q, want %q", got, want)
	}
	return opts.Unmarshal(src.GetValue(), dst)
}

// UnmarshalNew unmarshals the underlying message from src into dst,
// which is newly created message using a type resolved from the type URL.
// The message type is resolved according to opt.Resolver,
// which should implement protoregistry.MessageTypeResolver.
// It reports an error if the underlying message type could not be resolved.
//
// If no options are specified, call src.UnmarshalNew instead.
func UnmarshalNew(src *Any, opts proto.UnmarshalOptions) (dst proto.Message, err error) {
	if src.GetTypeUrl() == "" {
		return nil, protoimpl.X.NewError("invalid empty type URL")
	}
	if opts.Resolver == nil {
		opts.Resolver = protoregistry.GlobalTypes
	}
	r, ok := opts.Resolver.(protoregistry.MessageTypeResolver)
	if !ok {
		return nil, protoregistry.NotFound
	}
	mt, err := r.FindMessageByURL(src.GetTypeUrl())
	if err != nil {
		if err == protoregistry.NotFound {
			return nil, err
		}
		return nil, protoimpl.X.NewError("could not resolve %q: %v", src.GetTypeUrl(), err)
	}
	dst = mt.New().Interface()
	return dst, opts.Unmarshal(src.GetValue(), dst)
}

// MessageIs reports whether the underlying message is of the same type as m.
func (x *Any) MessageIs(m proto.Message) bool {
	if m == nil {
		return false
	}
	url := x.GetTypeUrl()
	name := string(m.ProtoReflect().Descriptor().FullName())
	if !strings.HasSuffix(url, name) {
		return false
	}
	return len(url) == len(name) || url[len(url)-len(name)-1] == '/'
}

// MessageName reports the full name of the underlying message,
// returning an empty string if invalid.
func (x *Any) MessageName() protoreflect.FullName {
	url := x.GetTypeUrl()
	name := protoreflect.FullName(url)
	if i := strings.LastIndexByte(url, '/'); i >= 0 {
		name = name[i+len("/"):]
	}
	if !name.IsValid() {
		return ""
	}
	return name
}

// MarshalFrom marshals m into x as the underlying message.
func (x *Any) MarshalFrom(m proto.Message) error {
	return MarshalFrom(x, m, proto.MarshalOptions{})
}

// UnmarshalTo unmarshals the contents of the underlying message of x into m.
// It resets m before performing the unmarshal operation.
// It reports an error if m is not of the right message type.
func (x *Any) UnmarshalTo(m proto.Message) error {
	return UnmarshalTo(x, m, proto.UnmarshalOptions{})
}

// UnmarshalNew unmarshals the contents of the underlying message of x into
// a newly allocated message of the specified type.
// It reports an error if the underlying message type could not be resolved.
func (x *Any) UnmarshalNew() (proto.Message, error) {
	return UnmarshalNew(x, proto.UnmarshalOptions{})
}
