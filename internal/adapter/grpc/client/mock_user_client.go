package client

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type MockUserClient struct {
	mock.Mock
}

func (r *MockUserClient) GetByUUID(ctx context.Context, UserUUID string) (*domain.User, error) {
	args := r.Called(ctx, UserUUID)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *MockUserClient) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := r.Called(ctx, email)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *MockUserClient) IsEmailUnique(ctx context.Context, email string) error {
	args := r.Called(ctx, email)
	return args.Error(0)
}

func (r *MockUserClient) Create(ctx context.Context, userParam domain.User) (*domain.User, error) {
	args := r.Called(ctx, userParam)
	if args.Get(0) != nil {
		return args.Get(0).(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *MockUserClient) VerifiedEmail(ctx context.Context, email string) error {
	args := r.Called(ctx, email)
	return args.Error(0)
}

func (r *MockUserClient) MarkWelcomeMessageSent(ctx context.Context, ID uint64) error {
	args := r.Called(ctx, ID)
	return args.Error(0)
}

func (r *MockUserClient) UpdateGoogleID(ctx context.Context, ID uint64, googleID string) error {
	args := r.Called(ctx, ID, googleID)
	return args.Error(0)
}

func (r *MockUserClient) UpdateLastLoginTime(ctx context.Context, ID uint64) error {
	args := r.Called(ctx, ID)
	return args.Error(0)
}

func (r *MockUserClient) UpdatePassword(ctx context.Context, ID uint64, password string) error {
	args := r.Called(ctx, ID, password)
	return args.Error(0)
}
