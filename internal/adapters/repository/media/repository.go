package media

import (
	"context"
	"database/sql"

	adapterlogging "github.com/rulzi/hexa-go/internal/adapters/logging"
	mediaentity "github.com/rulzi/hexa-go/internal/domain/media/entity"
	mediaport "github.com/rulzi/hexa-go/internal/domain/media/port"
)

var _ mediaport.Repository = (*MySQLRepository)(nil)

// MySQLRepository is the MySQL implementation of media.Repository (driven adapter)
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository creates a new MySQLRepository
func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// Create creates a new media
func (r *MySQLRepository) Create(ctx context.Context, m *mediaentity.Media) (*mediaentity.Media, error) {
	query := `
		INSERT INTO media (name, path, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, m.Name, m.Path, m.CreatedAt, m.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	m.ID = id
	return m, nil
}

// GetByID retrieves a media by ID
func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*mediaentity.Media, error) {
	query := `
		SELECT id, name, path, created_at, updated_at
		FROM media
		WHERE id = ?
	`

	m := &mediaentity.Media{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&m.ID,
		&m.Name,
		&m.Path,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, mediaentity.NewMediaNotFound()
	}
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Update updates an existing media
func (r *MySQLRepository) Update(ctx context.Context, m *mediaentity.Media) (*mediaentity.Media, error) {
	query := `
		UPDATE media
		SET name = ?, path = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, m.Name, m.Path, m.UpdatedAt, m.ID)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Delete deletes a media by ID
func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM media WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return mediaentity.NewMediaNotFound()
	}

	return nil
}

// List retrieves all media with pagination
func (r *MySQLRepository) List(ctx context.Context, limit, offset int) ([]*mediaentity.Media, error) {
	query := `
		SELECT id, name, path, created_at, updated_at
		FROM media
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			adapterlogging.LogRepoError(ctx, "failed to close rows: %v", err)
		}
	}()

	var mediaList []*mediaentity.Media
	for rows.Next() {
		m := &mediaentity.Media{}
		err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.Path,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		mediaList = append(mediaList, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return mediaList, nil
}

// Count returns the total number of media
func (r *MySQLRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM media`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
