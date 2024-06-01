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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (r Server) StartUserGRPCServer() (*grpc.Server, error) {
	address := fmt.Sprintf("%s:%s", r.cfg.URL, r.cfg.GRPCPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	grpcServer := grpc.NewServer()

	userpb.RegisterUserServiceServer(grpcServer, r)

	if err = grpcServer.Serve(listener); err != nil {
		return nil, err
	}

	return grpcServer, nil
}

func (r Server) GetByUUID(ctx context.Context, req *userpb.GetByUUIDRequest) (*userpb.UserResponse, error) {
	resp, err := r.userService.GetByUUID(ctx, req.GetUserUUID())
	if err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if resp != nil {
		return &userpb.UserResponse{
			Id:        resp.ID,
			UUID:      resp.UUID.String(),
			FirstName: resp.FirstName,
			LastName:  resp.LastName,
			Email:     resp.Email,
			Status:    resp.Status.String(),
		}, nil
	}

	return nil, nil
}

func (r Server) GetByEmail(ctx context.Context, req *userpb.GetByEmailRequest) (*userpb.UserResponse, error) {
	resp, err := r.userService.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if resp != nil {
		return &userpb.UserResponse{
			Id:                 resp.ID,
			UUID:               resp.UUID.String(),
			FirstName:          resp.FirstName,
			LastName:           resp.LastName,
			Email:              resp.Email,
			Password:           resp.Password,
			Status:             resp.Status.String(),
			WelcomeMessageSent: resp.WelcomeMessageSent,
			GoogleId:           resp.GoogleID,
		}, nil
	}

	return nil, nil
}

func (r Server) IsEmailUnique(ctx context.Context, req *userpb.IsEmailUniqueRequest) (*emptypb.Empty, error) {
	err := r.userService.IsEmailUnique(ctx, req.GetEmail())
	if err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
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
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}

func (r Server) VerifiedEmail(ctx context.Context, req *userpb.VerifiedEmailRequest) (*emptypb.Empty, error) {
	err := r.userService.VerifiedEmail(ctx, req.GetEmail())
	if err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}

func (r Server) MarkWelcomeMessageSent(ctx context.Context, req *userpb.UpdateWelcomeMessageToSentRequest) (*emptypb.Empty, error) {
	err := r.userService.MarkWelcomeMessageSent(ctx, req.GetUserId())
	if err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}
