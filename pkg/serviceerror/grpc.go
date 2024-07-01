package serviceerror

import (
	"github.com/mohsenabedy91/polyglot-sentences/proto/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var AnypbNew = anypb.New

func ConvertToGrpcError(serviceErr *ServiceError) error {
	st := status.New(codes.Unknown, serviceErr.Error())

	attrs := serviceErr.GetAttributes()
	strAttrs := make(map[string]interface{})
	for k, v := range attrs {
		strAttrs[k] = v
	}

	attributes, err := structpb.NewStruct(strAttrs)
	if err != nil {
		return st.Err()
	}

	customErrorDetail := &common.CustomErrorDetail{
		Message:    serviceErr.Error(),
		Attributes: attributes,
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
			for k, v := range customErrorDetail.Attributes.Fields {
				attrs[k] = v.AsInterface()
			}

			return New(ErrorMessage(customErrorDetail.Message), attrs)
		}
	}

	if st != nil {
		return st.Err()
	}

	return err
}
