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
	"github.com/mohsenabedy91/polyglot-sentences/proto/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
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

func (r *UserClient) Close() error {
	return r.conn.Close()
}

func (r *UserClient) GetByUUID(ctx context.Context, UserUUID string) (*domain.User, error) {
	req := userpb.GetByUUIDRequest{UserUuid: UserUUID}
	resp, err := r.userServiceClient.GetByUUID(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return nil, ExtractServiceErrorFromGrpcError(err)
	}

	return &domain.User{
		Base: domain.Base{
			ID:   uint(resp.GetId()),
			UUID: uuid.MustParse(resp.GetUuid()),
		},
		FirstName: resp.GetFirstName(),
		LastName:  resp.GetLastName(),
		Email:     resp.GetEmail(),
		Status:    domain.ToUserStatus(resp.Status),
	}, nil
}

func (r *UserClient) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	req := userpb.GetByEmailRequest{Email: email}
	resp, err := r.userServiceClient.GetByEmail(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return nil, ExtractServiceErrorFromGrpcError(err)
	}

	return &domain.User{
		Base: domain.Base{
			ID:   uint(resp.GetId()),
			UUID: uuid.MustParse(resp.GetUuid()),
		},
		FirstName: resp.GetFirstName(),
		LastName:  resp.GetLastName(),
		Email:     resp.GetEmail(),
		Password:  resp.GetPassword(),
		Status:    domain.ToUserStatus(resp.Status),
	}, nil
}

func (r *UserClient) IsEmailUnique(ctx context.Context, email string) error {
	req := userpb.IsEmailUniqueRequest{Email: email}
	_, err := r.userServiceClient.IsEmailUnique(ctx, &req)
	if err != nil {
		r.log.Error(logger.UserManagement, logger.API, err.Error(), map[logger.ExtraKey]interface{}{
			logger.RequestBody: &req,
		})
		return ExtractServiceErrorFromGrpcError(err)
	}
	return nil
}

func (r *UserClient) Create(ctx context.Context, userParam domain.User) error {
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
		return ExtractServiceErrorFromGrpcError(err)
	}

	return nil
}

func ExtractServiceErrorFromGrpcError(err error) error {
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

		return serviceerror.NewServiceError(serviceerror.ErrorMessage(customErrorDetail.Message), attrs)
	}

	return err
}
