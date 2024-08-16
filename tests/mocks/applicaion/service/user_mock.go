package service

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/wojciechpawlinow/usermanagement/internal/application/service"
	"github.com/wojciechpawlinow/usermanagement/internal/domain/user"
)

type UserServiceMock struct {
	mock.Mock
}

var _ service.UserPort = (*UserServiceMock)(nil)

func (m *UserServiceMock) Create(ctx context.Context, dto *service.CreateUserDTO) error {
	args := m.Called(ctx, dto)

	return args.Error(0)
}

func (m *UserServiceMock) Update(ctx context.Context, userID string, dto *service.UpdateUserDTO) error {
	args := m.Called(ctx, userID, dto)

	return args.Error(0)
}

func (m *UserServiceMock) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)

	return args.Error(0)
}

func (m *UserServiceMock) Get(ctx context.Context, page, pageSize int) ([]*user.User, error) {
	args := m.Called(ctx, page, pageSize)

	if val, ok := (args.Get(0)).([]*user.User); ok {
		return val, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *UserServiceMock) GetByUUID(ctx context.Context, userID string) (*user.User, error) {
	args := m.Called(ctx, userID)

	if val, ok := (args.Get(0)).(*user.User); ok {
		return val, args.Error(1)
	}

	return nil, args.Error(1)
}
