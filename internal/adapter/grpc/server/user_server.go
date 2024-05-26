package server

import (
	"context"
	"errors"
	"fmt"
	userpb "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/proto/user"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"github.com/mohsenabedy91/polyglot-sentences/proto/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
	userService port.UserService
	cfg         config.UserManagement
	trans       *translation.Translation
}

func NewUserGRPCServer(cfg config.UserManagement, userService port.UserService, trans *translation.Translation) Server {
	return Server{
		UnimplementedUserServiceServer: userpb.UnimplementedUserServiceServer{},
		userService:                    userService,
		cfg:                            cfg,
		trans:                          trans,
	}
}

func (r Server) StartUserGRPCServer() {
	address := fmt.Sprintf("%s:%s", r.cfg.URL, r.cfg.GRPCPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return
	}

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, r)

	if err = grpcServer.Serve(listener); err != nil {
		return
	}
}

func (r Server) GetByUUID(ctx context.Context, req *userpb.GetByUUIDRequest) (*userpb.UserResponse, error) {
	resp, err := r.userService.GetByUUID(ctx, req.GetUserUuid())
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &userpb.UserResponse{
		Id:        uint64(resp.ID),
		Uuid:      resp.UUID.String(),
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		Status:    resp.Status.String(),
	}, nil
}

func (r Server) GetByEmail(ctx context.Context, req *userpb.GetByEmailRequest) (*userpb.UserResponse, error) {
	resp, err := r.userService.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return &userpb.UserResponse{
		Id:        uint64(resp.ID),
		Uuid:      resp.UUID.String(),
		FirstName: resp.FirstName,
		LastName:  resp.LastName,
		Email:     resp.Email,
		Password:  resp.Password,
		Status:    resp.Status.String(),
	}, nil
}

func (r Server) IsEmailUnique(ctx context.Context, email *userpb.IsEmailUniqueRequest) (*emptypb.Empty, error) {
	err := r.userService.IsEmailUnique(ctx, email.GetEmail())
	if err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, ConvertServiceErrorToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}

func (r Server) Create(ctx context.Context, req *userpb.CreateRequest) (*emptypb.Empty, error) {
	err := r.userService.Create(ctx, domain.User{
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		Email:     req.GetEmail(),
		Password:  req.GetPassword(),
	})
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return nil, nil
}

func ConvertServiceErrorToGrpcError(serviceErr *serviceerror.ServiceError) error {
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
