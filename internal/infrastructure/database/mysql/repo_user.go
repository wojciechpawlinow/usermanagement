package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/wojciechpawlinow/usermanagement/internal/domain"
	"github.com/wojciechpawlinow/usermanagement/internal/domain/user"
	"github.com/wojciechpawlinow/usermanagement/internal/infrastructure/database/mysql/entity"
)

const (
	duplicatedEntry = 1062
)

type userRepository struct {
	dbRead  *sql.DB
	dbWrite *sql.DB
}

var _ user.Repository = (*userRepository)(nil)

func NewUserRepository(dbRead, dbWrite *sql.DB) *userRepository {
	return &userRepository{
		dbRead:  dbRead,
		dbWrite: dbWrite,
	}
}

func (r *userRepository) Create(ctx context.Context, u *user.User, createdAt time.Time) error {
	tx, err := r.dbWrite.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queryUser := `
		INSERT INTO users (uuid, email, password, created_at, updated_at, first_name, last_name, phone_number)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.dbWrite.ExecContext(ctx, queryUser, u.ID.String(), u.Email, u.Password, createdAt, nil, u.FirstName, u.LastName, u.PhoneNumber)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == duplicatedEntry { // duplicated entry
				return user.ErrEmailAlreadyExists
			}
		}

		return err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	addrQuery := `
		INSERT INTO addresses (user_id, type, street, city, state, postal_code, country, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	for _, addr := range u.Addresses {
		_, err = tx.ExecContext(ctx, addrQuery, userID, addr.Type, addr.Street, addr.City, addr.State, addr.PostalCode, addr.Country, createdAt, nil)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) {
				if mysqlErr.Number == duplicatedEntry {
					return user.ErrAddressAlreadyExists
				}
			}

			return err
		}
	}

	return tx.Commit()
}

func (r *userRepository) UpdateBasicFields(ctx context.Context, id domain.ID, fields map[string]any) error {
	existsQuery := "SELECT COUNT(1) FROM users WHERE uuid = ? AND deleted_at IS NULL "

	var exists int
	if err := r.dbRead.QueryRowContext(ctx, existsQuery, id.String()).Scan(&exists); err != nil {
		return fmt.Errorf("failed checking if user exists: %w", err)
	}

	if exists == 0 {
		return user.ErrNotFound
	}

	queryUser := "UPDATE users SET "
	var args []any
	i := 1

	for key, value := range fields {
		queryUser += key + " = ?"
		if i < len(fields) {
			queryUser += ", "
		}
		args = append(args, value)
		i++
	}

	queryUser += " WHERE uuid = ?"
	args = append(args, id.String())

	if _, err := r.dbWrite.ExecContext(ctx, queryUser, args...); err != nil {
		return fmt.Errorf("failed updating users: %w", err)
	}

	return nil
}

func (r *userRepository) UpdateAddress(ctx context.Context, id domain.ID, addrType int, fields map[string]any) error {
	existsQuery := "SELECT COUNT(1) FROM addresses WHERE user_id = (SELECT id FROM users WHERE uuid = ? AND deleted_at IS NULL) AND type = ?"

	var exists int
	_ = r.dbRead.QueryRowContext(ctx, existsQuery, id.String(), addrType).Scan(&exists)
	if exists == 0 {
		return user.ErrAddressNotFound
	}

	queryAddress := "UPDATE addresses SET "
	var args []any
	i := 1

	for key, value := range fields {
		queryAddress += key + " = ?"
		if i < len(fields) {
			queryAddress += ", "
		}
		args = append(args, value)
		i++
	}

	queryAddress += " WHERE user_id = (SELECT id FROM users WHERE uuid = ?) AND type = ?"
	args = append(args, id.String(), addrType)

	if _, err := r.dbWrite.ExecContext(ctx, queryAddress, args...); err != nil {
		return fmt.Errorf("failed updating address: %w", err)
	}

	return nil
}

func (r *userRepository) InsertAddress(ctx context.Context, id domain.ID, addr *user.Address, createdAt time.Time) error {
	query := `
		INSERT INTO addresses (user_id, type, street, city, state, postal_code, country, created_at, updated_at)
		VALUES ((SELECT id FROM users WHERE uuid = ? AND deleted_at IS NULL ), ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.dbWrite.ExecContext(ctx, query,
		id.String(),
		addr.Type,
		addr.Street,
		addr.City,
		addr.State,
		addr.PostalCode,
		addr.Country,
		createdAt,
		nil,
	)

	return err
}

func (r *userRepository) Delete(ctx context.Context, id domain.ID) error {
	tx, err := r.dbWrite.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queryID := "SELECT id FROM users WHERE uuid = ?"

	var userID int
	_ = r.dbRead.QueryRowContext(ctx, queryID, id.String()).Scan(&userID)
	if userID == 0 {
		return user.ErrNotFound
	}

	ts := time.Now()

	queryUser := "UPDATE users SET deleted_at = ? WHERE id = ? LIMIT 1"
	if _, err = r.dbWrite.ExecContext(ctx, queryUser, ts, userID); err != nil {
		return fmt.Errorf("failed deleting user: %w", err)
	}

	queryAddresses := "UPDATE addresses SET deleted_at = ? WHERE user_id = ?"
	if _, err = r.dbWrite.ExecContext(ctx, queryAddresses, ts, userID); err != nil {
		return fmt.Errorf("failed deleting addresses: %w", err)
	}

	return tx.Commit()
}

func (r *userRepository) GetByUUID(ctx context.Context, id domain.ID) (*user.User, error) {
	var dbUser entity.DbUser

	queryUser := "SELECT id, uuid, email, first_name, last_name, phone_number FROM users WHERE uuid = ? AND deleted_at IS NULL"

	row := r.dbRead.QueryRowContext(ctx, queryUser, id.String())
	err := row.Scan(&dbUser.ID, &dbUser.UUID, &dbUser.Email, &dbUser.FirstName, &dbUser.LastName, &dbUser.PhoneNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, fmt.Errorf("failed querying user: %w", err)
	}

	var domainAddresses []*user.Address

	queryAddresses := "SELECT type, street, city, state, postal_code, country FROM addresses WHERE user_id = ? AND deleted_at IS NULL"
	rows, err := r.dbRead.QueryContext(ctx, queryAddresses, dbUser.ID)
	if err != nil {
		return nil, fmt.Errorf("failed querying addresses: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var dbAddress entity.DbAddress
		if err := rows.Scan(&dbAddress.Type, &dbAddress.Street, &dbAddress.City, &dbAddress.State, &dbAddress.PostalCode, &dbAddress.Country); err != nil {
			return nil, fmt.Errorf("failed scanning addresses: %w", err)
		}

		domainAddress := &user.Address{
			Type:       user.AddressType(dbAddress.Type.Int64),
			Street:     dbAddress.Street.String,
			City:       dbAddress.City.String,
			State:      dbAddress.State.String,
			PostalCode: dbAddress.PostalCode.String,
			Country:    dbAddress.Country.String,
		}

		domainAddresses = append(domainAddresses, domainAddress)
	}

	userID, _ := domain.ParseID(dbUser.UUID.String)
	domainUser := &user.User{
		ID:          userID,
		Email:       dbUser.Email.String,
		Password:    "", // Password is not retrieved
		FirstName:   dbUser.FirstName.String,
		LastName:    dbUser.LastName.String,
		PhoneNumber: dbUser.PhoneNumber.String,
		Addresses:   domainAddresses,
	}

	return domainUser, nil
}

func (r *userRepository) Get(ctx context.Context, page, pageSize int) ([]*user.User, error) {
	offset := (page - 1) * pageSize

	var domainUsers []*user.User

	queryUsers := "SELECT id, uuid, email, first_name, last_name, phone_number FROM users WHERE deleted_at IS NULL LIMIT ? OFFSET ?"

	rowsUsers, err := r.dbRead.QueryContext(ctx, queryUsers, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed querying users: %w", err)
	}

	for rowsUsers.Next() {
		var dbUser entity.DbUser
		if err = rowsUsers.Scan(&dbUser.ID, &dbUser.UUID, &dbUser.Email, &dbUser.FirstName, &dbUser.LastName, &dbUser.PhoneNumber); err != nil {
			return nil, fmt.Errorf("failed scanning users: %w", err)
		}

		var domainAddresses []*user.Address

		queryAddresses := "SELECT type, street, city, state, postal_code, country FROM addresses WHERE user_id = ? AND deleted_at IS NULL "

		rowsAddresses, err := r.dbRead.QueryContext(ctx, queryAddresses, dbUser.ID.Int64)
		if err != nil {
			return nil, fmt.Errorf("failed querying addresses: %w", err)
		}

		for rowsAddresses.Next() {
			var dbAddress entity.DbAddress
			if err = rowsAddresses.Scan(&dbAddress.Type, &dbAddress.Street, &dbAddress.City, &dbAddress.State, &dbAddress.PostalCode, &dbAddress.Country); err != nil {
				return nil, fmt.Errorf("failed scanning addresses: %w", err)
			}

			domainAddress := &user.Address{
				Type:       user.AddressType(dbAddress.Type.Int64),
				Street:     dbAddress.Street.String,
				City:       dbAddress.City.String,
				State:      dbAddress.State.String,
				PostalCode: dbAddress.PostalCode.String,
				Country:    dbAddress.Country.String,
			}

			domainAddresses = append(domainAddresses, domainAddress)
		}

		userID, _ := domain.ParseID(dbUser.UUID.String)
		domainUser := &user.User{
			ID:          userID,
			Email:       dbUser.Email.String,
			Password:    "",
			FirstName:   dbUser.FirstName.String,
			LastName:    dbUser.LastName.String,
			PhoneNumber: dbUser.PhoneNumber.String,
			Addresses:   domainAddresses,
		}

		domainUsers = append(domainUsers, domainUser)
	}

	return domainUsers, nil
}
