package user

import (
	"context"
	"time"

	"github.com/wojciechpawlinow/usermanagement/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, u *User, createdAt time.Time) error
	UpdateBasicFields(ctx context.Context, id domain.ID, fields map[string]any) error
	UpdateAddress(ctx context.Context, id domain.ID, addrType int, fields map[string]any) error
	InsertAddress(ctx context.Context, id domain.ID, addr *Address, createdAt time.Time) error
	Delete(ctx context.Context, id domain.ID) error
	GetByUUID(ctx context.Context, id domain.ID) (*User, error)
	Get(ctx context.Context, page, pageSize int) ([]*User, error)
}
