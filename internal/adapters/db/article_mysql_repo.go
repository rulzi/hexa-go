package db

import (
	"context"
	"database/sql"

	"github.com/rulzi/hexa-go/internal/domain/article"
)

// ArticleMySQLRepository is the MySQL implementation of article.Repository (driven adapter)
type ArticleMySQLRepository struct {
	db *sql.DB
}

// NewArticleMySQLRepository creates a new ArticleMySQLRepository
func NewArticleMySQLRepository(db *sql.DB) *ArticleMySQLRepository {
	return &ArticleMySQLRepository{db: db}
}

// Create creates a new article
func (r *ArticleMySQLRepository) Create(ctx context.Context, a *article.Article) (*article.Article, error) {
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
func (r *ArticleMySQLRepository) GetByID(ctx context.Context, id int64) (*article.Article, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM articles
		WHERE id = ?
	`

	a := &article.Article{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID,
		&a.Title,
		&a.Content,
		&a.AuthorID,
		&a.CreatedAt,
		&a.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, article.ErrArticleNotFound
	}
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Update updates an existing article
func (r *ArticleMySQLRepository) Update(ctx context.Context, a *article.Article) (*article.Article, error) {
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
func (r *ArticleMySQLRepository) Delete(ctx context.Context, id int64) error {
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
		return article.ErrArticleNotFound
	}

	return nil
}

// List retrieves all articles with pagination
func (r *ArticleMySQLRepository) List(ctx context.Context, limit, offset int) ([]*article.Article, error) {
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
	defer rows.Close()

	var articles []*article.Article
	for rows.Next() {
		a := &article.Article{}
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
func (r *ArticleMySQLRepository) ListByAuthor(ctx context.Context, authorID int64, limit, offset int) ([]*article.Article, error) {
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
	defer rows.Close()

	var articles []*article.Article
	for rows.Next() {
		a := &article.Article{}
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
func (r *ArticleMySQLRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM articles`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CountByAuthor returns the total number of articles by author
func (r *ArticleMySQLRepository) CountByAuthor(ctx context.Context, authorID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM articles WHERE author_id = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, authorID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// CreateTable creates the articles table if it doesn't exist
func (r *ArticleMySQLRepository) CreateTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS articles (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			author_id BIGINT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_author_id (author_id),
			INDEX idx_created_at (created_at),
			FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4
	`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

