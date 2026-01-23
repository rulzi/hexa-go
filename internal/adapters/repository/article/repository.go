package article

import (
	"context"
	"database/sql"
	"log"

	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
)

// MySQLRepository is the MySQL implementation of article.Repository (driven adapter)
type MySQLRepository struct {
	db *sql.DB
}

// NewMySQLRepository creates a new MySQLRepository
func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{db: db}
}

// Create creates a new article
func (r *MySQLRepository) Create(ctx context.Context, a *domainarticle.Article) (*domainarticle.Article, error) {
	query := `
		INSERT INTO articles (title, content, author_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query, a.Title, a.Content, a.AuthorID, a.CreatedAt, a.UpdatedAt)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	a.ID = id
	return a, nil
}

// GetByID retrieves an article by ID
func (r *MySQLRepository) GetByID(ctx context.Context, id int64) (*domainarticle.Article, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM articles
		WHERE id = ?
	`

	a := &domainarticle.Article{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID,
		&a.Title,
		&a.Content,
		&a.AuthorID,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domainarticle.ErrArticleNotFound
	}
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Update updates an existing article
func (r *MySQLRepository) Update(ctx context.Context, a *domainarticle.Article) (*domainarticle.Article, error) {
	query := `
		UPDATE articles
		SET title = ?, content = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, a.Title, a.Content, a.UpdatedAt, a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Delete deletes an article by ID
func (r *MySQLRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM articles WHERE id = ?`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domainarticle.ErrArticleNotFound
	}

	return nil
}

// List retrieves all articles with pagination
func (r *MySQLRepository) List(ctx context.Context, limit, offset int) ([]*domainarticle.Article, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM articles
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	var articles []*domainarticle.Article
	for rows.Next() {
		a := &domainarticle.Article{}
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Content,
			&a.AuthorID,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

// ListByAuthor retrieves articles by author ID with pagination
func (r *MySQLRepository) ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*domainarticle.Article, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM articles
		WHERE author_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}()

	var articles []*domainarticle.Article
	for rows.Next() {
		a := &domainarticle.Article{}
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Content,
			&a.AuthorID,
			&a.CreatedAt,
			&a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}

// Count returns the total number of articles
func (r *MySQLRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM articles`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountByAuthor returns the total number of articles by author
func (r *MySQLRepository) CountByAuthor(ctx context.Context, authorID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM articles WHERE author_id = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, authorID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
