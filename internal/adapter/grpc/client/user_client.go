package client

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	userpb "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/grpc/proto/user"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClient struct {
	log               logger.Logger
	conn              *grpc.ClientConn
	userServiceClient userpb.UserServiceClient
}

func NewUserClient(log logger.Logger, conf config.UserManagement) *UserClient {
	target := fmt.Sprintf("%s:%s", conf.URL, conf.GRPCPort)
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(logger.Internal, logger.Startup, fmt.Sprintf("There is an error when run http: %v", err), nil)
		return nil
	}

	client := userpb.NewUserServiceClient(conn)
	return &UserClient{
		conn:              conn,
		log:               log,
		userServiceClient: client,
	}
}

func (r UserClient) Close() {
	if err := r.conn.Close(); err != nil {
		r.log.Error(logger.Internal, logger.Startup, fmt.Sprintf("Failed to close client connection: %v", err), nil)
	}
	return
}

func (r UserClient) GetByUUID(ctx context.Context, UserUUID string) (*domain.User, error) {
	req := userpb.GetByUUIDRequest{UserUUID: UserUUID}
	resp, err := r.userServiceClient.GetByUUID(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return nil, serviceerror.ExtractFromGrpcError(err)
	}

	if resp.String() != "" {
		return &domain.User{
			Base: domain.Base{
				ID:   resp.Id,
				UUID: uuid.MustParse(resp.UUID),
			},
			FirstName: resp.FirstName,
			LastName:  resp.LastName,
			Email:     resp.Email,
			Status:    domain.ToUserStatus(resp.Status),
		}, nil
	}

	return nil, nil
}

func (r UserClient) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	req := userpb.GetByEmailRequest{Email: email}
	resp, err := r.userServiceClient.GetByEmail(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return nil, serviceerror.ExtractFromGrpcError(err)
	}

	if resp.String() != "" {
		return &domain.User{
			Base: domain.Base{
				ID:   resp.Id,
				UUID: uuid.MustParse(resp.UUID),
			},
			FirstName:          resp.FirstName,
			LastName:           resp.LastName,
			Email:              resp.Email,
			Password:           resp.Password,
			Status:             domain.ToUserStatus(resp.Status),
			WelcomeMessageSent: resp.WelcomeMessageSent,
			GoogleID:           resp.GoogleId,
		}, nil
	}

	return nil, nil
}

func (r UserClient) IsEmailUnique(ctx context.Context, email string) error {
	req := userpb.IsEmailUniqueRequest{Email: email}
	_, err := r.userServiceClient.IsEmailUnique(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return serviceerror.ExtractFromGrpcError(err)
	}
	return nil
}

func (r UserClient) Create(ctx context.Context, userParam domain.User) (*domain.User, error) {
	req := userpb.CreateRequest{
		FirstName: userParam.FirstName,
		LastName:  userParam.LastName,
		Email:     userParam.Email,
		Password:  userParam.Password,
		Avatar:    userParam.Avatar,
		GoogleId:  userParam.GoogleID,
		Status:    userParam.Status.String(),
	}
	resp, err := r.userServiceClient.Create(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return nil, serviceerror.ExtractFromGrpcError(err)
	}

	if resp.String() != "" {
		return &domain.User{
			Base: domain.Base{
				ID:   resp.Id,
				UUID: uuid.MustParse(resp.UUID),
			},
			FirstName:          resp.FirstName,
			LastName:           resp.LastName,
			Email:              resp.Email,
			Password:           resp.Password,
			Status:             domain.ToUserStatus(resp.Status),
			WelcomeMessageSent: resp.WelcomeMessageSent,
			GoogleID:           resp.GoogleId,
		}, nil
	}

	return nil, nil
}

func (r UserClient) VerifiedEmail(ctx context.Context, email string) error {
	req := userpb.VerifiedEmailRequest{Email: email}
	_, err := r.userServiceClient.VerifiedEmail(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return serviceerror.ExtractFromGrpcError(err)
	}
	return nil
}

func (r UserClient) MarkWelcomeMessageSent(ctx context.Context, ID uint64) error {
	req := userpb.UpdateWelcomeMessageToSentRequest{UserId: ID}
	_, err := r.userServiceClient.MarkWelcomeMessageSent(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return serviceerror.ExtractFromGrpcError(err)
	}
	return nil
}

func (r UserClient) UpdateGoogleID(ctx context.Context, ID uint64, googleID string) error {
	req := userpb.UpdateGoogleIDRequest{UserId: ID, GoogleId: googleID}
	_, err := r.userServiceClient.UpdateGoogleID(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return serviceerror.ExtractFromGrpcError(err)
	}
	return nil
}

func (r UserClient) UpdateLastLoginTime(ctx context.Context, ID uint64) error {
	req := userpb.UpdateLastLoginTimeRequest{UserId: ID}
	_, err := r.userServiceClient.UpdateLastLoginTime(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return serviceerror.ExtractFromGrpcError(err)
	}
	return nil
}

func (r UserClient) UpdatePassword(ctx context.Context, ID uint64, password string) error {
	req := userpb.UpdatePasswordRequest{UserId: ID, Password: password}
	_, err := r.userServiceClient.UpdatePassword(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return serviceerror.ExtractFromGrpcError(err)
	}
	return nil
}
