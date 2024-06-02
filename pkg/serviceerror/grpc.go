package serviceerror

import (
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/proto/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

func ConvertToGrpcError(serviceErr *ServiceError) error {
	st := status.New(codes.Unknown, serviceErr.Error())

	attrs := serviceErr.GetAttributes()
	strAttrs := make(map[string]string)
	for k, v := range attrs {
		strAttrs[k] = fmt.Sprintf("%v", v)
	}

	customErrorDetail := &common.CustomErrorDetail{
		Message:    serviceErr.Error(),
		Attributes: strAttrs,
	}

	detail, err := anypb.New(customErrorDetail)
	if err != nil {
		return st.Err()
	}

	if st, err = st.WithDetails(detail); err != nil {
		return st.Err()
	}

	return st.Err()
}

func ExtractFromGrpcError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	for _, detail := range st.Details() {
		anyDetail, ok := detail.(*anypb.Any)
		if !ok {
			continue
		}

		var customErrorDetail common.CustomErrorDetail
		if err = anyDetail.UnmarshalTo(&customErrorDetail); err != nil {
			continue
		}

		attrs := make(map[string]interface{})
		for k, v := range customErrorDetail.Attributes {
			attrs[k] = v
		}

		return New(ErrorMessage(customErrorDetail.Message), attrs)
	}

	return err
}
