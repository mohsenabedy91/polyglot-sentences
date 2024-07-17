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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"net"
)

type Server struct {
	userpb.UnimplementedUserServiceServer
	conf        config.UserManagement
	userService port.UserService
	uowFactory  func() port.UserUnitOfWork
}

func NewUserGRPCServer(
	conf config.UserManagement,
	userService port.UserService,
	uowFactory func() port.UserUnitOfWork,
) *Server {
	return &Server{
		UnimplementedUserServiceServer: userpb.UnimplementedUserServiceServer{},
		conf:                           conf,
		userService:                    userService,
		uowFactory:                     uowFactory,
	}
}

func (r Server) StartUserGRPCServer() (*grpc.Server, error) {
	address := fmt.Sprintf("%s:%s", r.conf.URL, r.conf.GRPCPort)
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
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	resp, err := r.userService.GetByUUID(uowFactory, req.GetUserUUID())
	if err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			var se *serviceerror.ServiceError
			if errors.As(err, &se) {
				return nil, serviceerror.ConvertToGrpcError(se)
			}
		}
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := uowFactory.Commit(); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if resp != nil {
		return &userpb.UserResponse{
			Id:        resp.Base.ID,
			UUID:      resp.Base.UUID.String(),
			FirstName: resp.FirstName,
			LastName:  resp.LastName,
			Email:     resp.Email,
			Status:    resp.Status.String(),
		}, nil
	}

	return nil, nil
}

func (r Server) GetByEmail(ctx context.Context, req *userpb.GetByEmailRequest) (*userpb.UserResponse, error) {
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	resp, err := r.userService.GetByEmail(uowFactory, req.GetEmail())
	if err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			var se *serviceerror.ServiceError
			if errors.As(err, &se) {
				return nil, serviceerror.ConvertToGrpcError(se)
			}
		}
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := uowFactory.Commit(); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if resp != nil {
		return &userpb.UserResponse{
			Id:                 resp.Base.ID,
			UUID:               resp.Base.UUID.String(),
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
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := r.userService.IsEmailUnique(uowFactory, req.GetEmail()); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			var se *serviceerror.ServiceError
			if errors.As(err, &se) {
				return nil, serviceerror.ConvertToGrpcError(se)
			}
		}
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := uowFactory.Commit(); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}

func (r Server) Create(ctx context.Context, req *userpb.CreateRequest) (*userpb.UserResponse, error) {
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	resp, err := r.userService.Create(uowFactory, domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
		Avatar:    req.Avatar,
		GoogleID:  req.GoogleId,
		Status:    domain.ToUserStatus(req.Status),
	})
	if err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			var se *serviceerror.ServiceError
			if errors.As(err, &se) {
				return nil, serviceerror.ConvertToGrpcError(se)
			}
		}
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err = uowFactory.Commit(); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if resp != nil {
		return &userpb.UserResponse{
			Id:                 resp.Base.ID,
			UUID:               resp.Base.UUID.String(),
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

func (r Server) VerifiedEmail(ctx context.Context, req *userpb.VerifiedEmailRequest) (*emptypb.Empty, error) {
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := r.userService.VerifiedEmail(uowFactory, req.GetEmail()); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			var se *serviceerror.ServiceError
			if errors.As(err, &se) {
				return nil, serviceerror.ConvertToGrpcError(se)
			}
		}
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := uowFactory.Commit(); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}

func (r Server) MarkWelcomeMessageSent(ctx context.Context, req *userpb.UpdateWelcomeMessageToSentRequest) (*emptypb.Empty, error) {
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := r.userService.MarkWelcomeMessageSent(uowFactory, req.GetUserId()); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			var se *serviceerror.ServiceError
			if errors.As(err, &se) {
				return nil, serviceerror.ConvertToGrpcError(se)
			}
		}
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := uowFactory.Commit(); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}

func (r Server) UpdateGoogleID(ctx context.Context, req *userpb.UpdateGoogleIDRequest) (*emptypb.Empty, error) {
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := r.userService.UpdateGoogleID(uowFactory, req.GetUserId(), req.GetGoogleId()); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			var se *serviceerror.ServiceError
			if errors.As(err, &se) {
				return nil, serviceerror.ConvertToGrpcError(se)
			}
		}
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := uowFactory.Commit(); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}

func (r Server) UpdateLastLoginTime(ctx context.Context, req *userpb.UpdateLastLoginTimeRequest) (*emptypb.Empty, error) {
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := r.userService.UpdateLastLoginTime(uowFactory, req.GetUserId()); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			var se *serviceerror.ServiceError
			if errors.As(err, &se) {
				return nil, serviceerror.ConvertToGrpcError(se)
			}
		}
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := uowFactory.Commit(); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}

func (r Server) UpdatePassword(ctx context.Context, req *userpb.UpdatePasswordRequest) (*emptypb.Empty, error) {
	uowFactory := r.uowFactory()
	if err := uowFactory.BeginTx(ctx); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := r.userService.UpdatePassword(uowFactory, req.GetUserId(), req.GetPassword()); err != nil {
		if rErr := uowFactory.Rollback(); rErr != nil {
			var se *serviceerror.ServiceError
			if errors.As(err, &se) {
				return nil, serviceerror.ConvertToGrpcError(se)
			}
			return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
		}
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	if err := uowFactory.Commit(); err != nil {
		var se *serviceerror.ServiceError
		if errors.As(err, &se) {
			return nil, serviceerror.ConvertToGrpcError(se)
		}
		return nil, status.Errorf(codes.Internal, "unexpected error: %v", err)
	}

	return nil, nil
}
