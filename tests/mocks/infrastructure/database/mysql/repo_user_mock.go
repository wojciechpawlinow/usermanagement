package mysql

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/wojciechpawlinow/usermanagement/internal/domain"
	"github.com/wojciechpawlinow/usermanagement/internal/domain/user"
)

type UserRepositoryMock struct {
	mock.Mock
}

var _ user.Repository = (*UserRepositoryMock)(nil)

func (m *UserRepositoryMock) Create(ctx context.Context, u *user.User, createdAt time.Time) error {
	args := m.Called(ctx, u, createdAt)

	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateBasicFields(ctx context.Context, id domain.ID, fields map[string]any) error {
	args := m.Called(ctx, id, fields)

	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateAddress(ctx context.Context, id domain.ID, addrType int, fields map[string]any) error {
	args := m.Called(ctx, id, addrType, fields)

	return args.Error(0)
}

func (m *UserRepositoryMock) InsertAddress(ctx context.Context, id domain.ID, addr *user.Address, createdAt time.Time) error {
	args := m.Called(ctx, id, addr, createdAt)

	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(ctx context.Context, id domain.ID) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

func (m *UserRepositoryMock) GetByUUID(ctx context.Context, id domain.ID) (*user.User, error) {
	args := m.Called(ctx, id)

	if val, ok := args.Get(0).(*user.User); ok {
		return val, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *UserRepositoryMock) Get(ctx context.Context, page, pageSize int) ([]*user.User, error) {
	args := m.Called(ctx, page, pageSize)

	if val, ok := args.Get(0).([]*user.User); ok {
		return val, args.Error(1)
	}

	return nil, args.Error(1)
}
