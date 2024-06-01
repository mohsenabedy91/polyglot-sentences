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

func NewUserClient(log logger.Logger, cfg config.UserManagement) *UserClient {
	target := fmt.Sprintf("%s:%s", cfg.URL, cfg.GRPCPort)
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

func (r UserClient) Close() error {
	return r.conn.Close()
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
				ID:   resp.GetId(),
				UUID: uuid.MustParse(resp.GetUUID()),
			},
			FirstName: resp.GetFirstName(),
			LastName:  resp.GetLastName(),
			Email:     resp.GetEmail(),
			Status:    domain.ToUserStatus(resp.GetStatus()),
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
				ID:   resp.GetId(),
				UUID: uuid.MustParse(resp.GetUUID()),
			},
			FirstName:          resp.GetFirstName(),
			LastName:           resp.GetLastName(),
			Email:              resp.GetEmail(),
			Password:           resp.GetPassword(),
			Status:             domain.ToUserStatus(resp.Status),
			WelcomeMessageSent: resp.GetWelcomeMessageSent(),
			GoogleID:           resp.GetGoogleId(),
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

func (r UserClient) Create(ctx context.Context, userParam domain.User) error {
	req := userpb.CreateRequest{
		FirstName: userParam.FirstName,
		LastName:  userParam.LastName,
		Email:     userParam.Email,
		Password:  userParam.Password,
	}
	_, err := r.userServiceClient.Create(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return serviceerror.ExtractFromGrpcError(err)
	}
	return nil
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
