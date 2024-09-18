/*
Copyright 2024 Yuchen Cheng.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package go_grpc_protovalidate

import (
	"context"
	"errors"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type validator struct {
	*protovalidate.Validator
}

func (v *validator) Validate(ctx context.Context, m interface{}) error {
	msg := m.(proto.Message)
	err := v.Validator.Validate(msg)

	var valErr *protovalidate.ValidationError
	if errors.As(err, &valErr) {
		st := status.New(codes.InvalidArgument, err.Error())
		print(err.Error())

		violations := make([]*errdetails.BadRequest_FieldViolation, 0, len(valErr.Violations))
		for _, v := range valErr.Violations {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       v.GetFieldPath(),
				Description: v.GetMessage(),
			})
		}

		ds, err := st.WithDetails(
			&errdetails.ErrorInfo{Reason: "INVALID_ARGUMENT"},
			&errdetails.BadRequest{FieldViolations: violations},
		)
		if err != nil {
			return st.Err()
		}
		return ds.Err()
	}

	return nil
}

func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateOpts(opts)
	v := &validator{Validator: o.validator}

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := v.Validate(ctx, req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateOpts(opts)
	v := &validator{Validator: o.validator}

	return func(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return handler(srv, &serverStream{ServerStream: stream, ctx: stream.Context(), validator: v})
	}
}

type serverStream struct {
	grpc.ServerStream

	ctx       context.Context
	validator *validator
}

func (ss *serverStream) RecvMsg(m interface{}) error {
	if err := ss.ServerStream.RecvMsg(m); err != nil {
		return err
	}
	return ss.validator.Validate(ss.ctx, m)
}
