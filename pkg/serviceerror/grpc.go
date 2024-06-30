package serviceerror

import (
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/proto/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

var AnypbNew = anypb.New

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

	detail, err := AnypbNew(customErrorDetail)
	if err != nil {
		return st.Err()
	}

	st, _ = st.WithDetails(detail)

	return st.Err()
}

func ExtractFromGrpcError(err error) error {
	st, ok := status.FromError(err)
	if ok && st != nil {
		for _, detail := range st.Details() {
			anyDetail, _ := detail.(*anypb.Any)

			var customErrorDetail common.CustomErrorDetail
			_ = anyDetail.UnmarshalTo(&customErrorDetail)

			attrs := make(map[string]interface{})
			for k, v := range customErrorDetail.Attributes {
				attrs[k] = v
			}

			return New(ErrorMessage(customErrorDetail.Message), attrs)
		}
	}

	if st != nil {
		return st.Err()
	}

	return err
}
