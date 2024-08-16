package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/wojciechpawlinow/usermanagement/internal/domain"
	"github.com/wojciechpawlinow/usermanagement/internal/domain/user"
	"github.com/wojciechpawlinow/usermanagement/pkg/logger"
)

type UserPort interface {
	Create(ctx context.Context, dto *CreateUserDTO) error
	Update(ctx context.Context, userID string, dto *UpdateUserDTO) error
	Delete(ctx context.Context, userID string) error
	Get(ctx context.Context, page, pageSize int) ([]*user.User, error)
	GetByUUID(ctx context.Context, userID string) (*user.User, error)
}

type CreateUserDTO struct {
	ID          domain.ID
	Email       string
	Password    string
	FirstName   string
	LastName    string
	PhoneNumber string
	Addresses   []*CreateUserAddress
}

type CreateUserAddress struct {
	Type       int
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

type UpdateUserDTO struct {
	Password    *string
	FirstName   *string
	LastName    *string
	PhoneNumber *string
	Addresses   []*UpdateUserAddress
}

type UpdateUserAddress struct {
	Type       int
	Street     *string
	City       *string
	State      *string
	PostalCode *string
	Country    *string
}

type userService struct {
	userRepo     user.Repository
	timeProvider domain.TimeProvider
}

var _ UserPort = (*userService)(nil)

func NewUserService(userRepo user.Repository, timeProvider domain.TimeProvider) *userService {
	return &userService{
		userRepo:     userRepo,
		timeProvider: timeProvider,
	}
}

func (s *userService) Create(ctx context.Context, dto *CreateUserDTO) error {
	u := &user.User{
		ID:          dto.ID,
		Email:       dto.Email,
		Password:    dto.Password,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		PhoneNumber: dto.PhoneNumber,
		Addresses:   make([]*user.Address, 0, len(dto.Addresses)),
	}

	for _, addr := range dto.Addresses {
		u.Addresses = append(u.Addresses, &user.Address{
			Type:       user.AddressType(addr.Type),
			Street:     addr.Street,
			City:       addr.City,
			State:      addr.State,
			PostalCode: addr.PostalCode,
			Country:    addr.Country,
		})
	}

	if err := s.userRepo.Create(ctx, u, s.timeProvider.UtcNow()); err != nil {
		if errors.Is(err, user.ErrEmailAlreadyExists) {
			return user.ErrEmailAlreadyExists
		}
		err = fmt.Errorf("failed creating user: %w", err)
		logger.Debug(err)

		return err
	}

	return nil
}

func (s *userService) Update(ctx context.Context, userID string, dto *UpdateUserDTO) error {
	userFields := make(map[string]any)
	id, err := domain.ParseID(userID)
	if err != nil {
		return fmt.Errorf("failed parsing uuid: %w", err)
	}

	if dto.Password != nil {
		userFields["password"] = *dto.Password
	}
	if dto.FirstName != nil {
		userFields["first_name"] = *dto.FirstName
	}
	if dto.LastName != nil {
		userFields["last_name"] = *dto.LastName
	}
	if dto.PhoneNumber != nil {
		userFields["phone_number"] = *dto.PhoneNumber
	}

	if len(userFields) > 0 {
		if err = s.userRepo.UpdateBasicFields(ctx, id, userFields); err != nil {
			err = fmt.Errorf("failed updating user personal data: %w", err)
			logger.Debug(err)

			return err
		}
	}

	for _, addr := range dto.Addresses {
		addrFields := make(map[string]any)

		if addr.Street != nil {
			addrFields["street"] = *addr.Street
		}
		if addr.City != nil {
			addrFields["city"] = *addr.City
		}
		if addr.State != nil {
			addrFields["state"] = *addr.State
		}
		if addr.PostalCode != nil {
			addrFields["postal_code"] = *addr.PostalCode
		}
		if addr.Country != nil {
			addrFields["country"] = *addr.Country
		}

		if len(addrFields) > 0 {
			if err = s.userRepo.UpdateAddress(ctx, id, addr.Type, addrFields); err != nil {
				if errors.Is(err, user.ErrAddressNotFound) {
					newAddr := &user.Address{
						Type: user.AddressType(addr.Type),
					}
					if addr.Street != nil {
						newAddr.Street = *addr.Street
					}
					if addr.City != nil {
						newAddr.City = *addr.City
					}
					if addr.State != nil {
						newAddr.State = *addr.State
					}
					if addr.PostalCode != nil {
						newAddr.PostalCode = *addr.PostalCode
					}
					if addr.Country != nil {
						newAddr.Country = *addr.Country
					}

					if err = s.userRepo.InsertAddress(ctx, id, newAddr, s.timeProvider.UtcNow()); err != nil {
						err = fmt.Errorf("failed inserting additional address")
						logger.Debug(err)

						return err
					}
				} else {
					err = fmt.Errorf("failed updating user address data: %w", err)
					logger.Debug(err)

					return err
				}

				return nil
			}
		}
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, userID string) error {
	id, err := domain.ParseID(userID)
	if err != nil {
		return fmt.Errorf("failed parsing uuid: %w", err)
	}

	if err = s.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return err
		}

		err = fmt.Errorf("failed deleting user: %w", err)
		logger.Debug(err)

		return err
	}

	return nil
}

func (s *userService) Get(ctx context.Context, page, pageSize int) ([]*user.User, error) {
	return s.userRepo.Get(ctx, page, pageSize)
}

func (s *userService) GetByUUID(ctx context.Context, userID string) (*user.User, error) {
	id, err := domain.ParseID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed parsing uuid: %w", err)
	}

	return s.userRepo.GetByUUID(ctx, id)
}
