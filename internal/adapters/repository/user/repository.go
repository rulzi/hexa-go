package user

import (
	"context"
	"database/sql"

	adapterlogging "github.com/rulzi/hexa-go/internal/adapters/logging"
	userentity "github.com/rulzi/hexa-go/internal/domain/user/entity"
	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

var _ userport.Repository = (*MySQLRepository)(nil)

// MySQLRepository is the MySQL implementation of user.Repository (driven adapter)
type MySQLRepository struct {
	db         *sql.DB
	repoLogger *adapterlogging.RepoLogger
}

// NewMySQLRepository creates a new MySQLRepository
func NewMySQLRepository(db *sql.DB, appLogger logger.Logger) *MySQLRepository {
	return &MySQLRepository{
		db:         db,
		repoLogger: adapterlogging.NewRepoLogger(appLogger),
	}
}

// Create creates a new user
func (r *MySQLRepository) Create(ctx context.Context, u *userentity.User) (*userentity.User, error) {
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.Password, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	u.ID = id
	return u, nil
}

// GetByID retrieves a user by ID
func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*userentity.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	u := &userentity.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, userentity.NewUserNotFound()
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetByEmail retrieves a user by email
func (r *MySQLRepository) GetByEmail(ctx context.Context, email string) (*userentity.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	u := &userentity.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, userentity.NewUserNotFound()
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Update updates an existing user
func (r *MySQLRepository) Update(ctx context.Context, u *userentity.User) (*userentity.User, error) {
	query := `
		UPDATE users
		SET name = ?, email = ?, password = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, u.Name, u.Email, u.Password, u.UpdatedAt, u.ID)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Delete deletes a user by ID
func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return userentity.NewUserNotFound()
	}

	return nil
}

// List retrieves all users with pagination
func (r *MySQLRepository) List(ctx context.Context, limit, offset int) ([]*userentity.User, error) {
	query := `
		SELECT id, name, email, password, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.repoLogger.LogError(ctx, "failed to close rows: %v", err)
		}
	}()

	var users []*userentity.User
	for rows.Next() {
		u := &userentity.User{}
		err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Password,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Count returns the total number of users
func (r *MySQLRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM users`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
