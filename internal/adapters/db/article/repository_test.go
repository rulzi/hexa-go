package article

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	domainarticle "github.com/rulzi/hexa-go/internal/domain/article"
	"github.com/stretchr/testify/assert"
)

func TestNewMySQLRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer func() {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatalf("Failed to close database connection: %v", err)
		}
	}()

	repo := NewMySQLRepository(db)
	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestMySQLRepository_Create(t *testing.T) {
	tests := []struct {
		name    string
		article *domainarticle.Article
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
		check   func(t *testing.T, article *domainarticle.Article)
	}{
		{
			name: "success create article",
			article: &domainarticle.Article{
				Title:     "Test Article",
				Content:   "This is a test article content",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO articles").
					WithArgs("Test Article", "This is a test article content", int64(1), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
			check: func(t *testing.T, article *domainarticle.Article) {
				assert.Equal(t, int64(1), article.ID)
				assert.Equal(t, "Test Article", article.Title)
				assert.Equal(t, "This is a test article content", article.Content)
				assert.Equal(t, int64(1), article.AuthorID)
			},
		},
		{
			name: "error on database exec",
			article: &domainarticle.Article{
				Title:     "Test Article",
				Content:   "This is a test article content",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO articles").
					WithArgs("Test Article", "This is a test article content", int64(1), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "error on last insert id",
			article: &domainarticle.Article{
				Title:     "Test Article",
				Content:   "This is a test article content",
				AuthorID:  1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO articles").
					WithArgs("Test Article", "This is a test article content", int64(1), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("last insert id error")))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Fatalf("Failed to close database connection: %v", err)
				}
			}()

			repo := NewMySQLRepository(db)
			tt.setup(mock)

			result, err := repo.Create(context.Background(), tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.check != nil {
					tt.check(t, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMySQLRepository_GetByID(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
		check   func(t *testing.T, article *domainarticle.Article)
	}{
		{
			name: "success get article by id",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow(1, "Test Article", "Test Content", 1, time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, article *domainarticle.Article) {
				assert.Equal(t, int64(1), article.ID)
				assert.Equal(t, "Test Article", article.Title)
				assert.Equal(t, "Test Content", article.Content)
				assert.Equal(t, int64(1), article.AuthorID)
			},
		},
		{
			name: "article not found",
			id:   999,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			check: func(t *testing.T, article *domainarticle.Article) {
				assert.Nil(t, article)
			},
		},
		{
			name: "database error",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Fatalf("Failed to close database connection: %v", err)
				}
			}()

			repo := NewMySQLRepository(db)
			tt.setup(mock)

			result, err := repo.GetByID(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "article not found" {
					assert.Equal(t, domainarticle.ErrArticleNotFound, err)
				}
				if tt.check != nil {
					tt.check(t, result)
				} else {
					assert.Nil(t, result)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.check != nil {
					tt.check(t, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMySQLRepository_Update(t *testing.T) {
	tests := []struct {
		name    string
		article *domainarticle.Article
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
		check   func(t *testing.T, article *domainarticle.Article)
	}{
		{
			name: "success update article",
			article: &domainarticle.Article{
				ID:        1,
				Title:     "Updated Article",
				Content:   "Updated Content",
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE articles").
					WithArgs("Updated Article", "Updated Content", sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
			check: func(t *testing.T, article *domainarticle.Article) {
				assert.Equal(t, int64(1), article.ID)
				assert.Equal(t, "Updated Article", article.Title)
				assert.Equal(t, "Updated Content", article.Content)
			},
		},
		{
			name: "error on database exec",
			article: &domainarticle.Article{
				ID:        1,
				Title:     "Updated Article",
				Content:   "Updated Content",
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE articles").
					WithArgs("Updated Article", "Updated Content", sqlmock.AnyArg(), 1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Fatalf("Failed to close database connection: %v", err)
				}
			}()

			repo := NewMySQLRepository(db)
			tt.setup(mock)

			result, err := repo.Update(context.Background(), tt.article)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.check != nil {
					tt.check(t, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMySQLRepository_Delete(t *testing.T) {
	tests := []struct {
		name    string
		id      int64
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
	}{
		{
			name: "success delete article",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM articles").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "article not found",
			id:   999,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM articles").
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "error on database exec",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM articles").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "error on rows affected",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM articles").
					WithArgs(1).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Fatalf("Failed to close database connection: %v", err)
				}
			}()

			repo := NewMySQLRepository(db)
			tt.setup(mock)

			err = repo.Delete(context.Background(), tt.id)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "article not found" {
					assert.Equal(t, domainarticle.ErrArticleNotFound, err)
				}
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMySQLRepository_List(t *testing.T) {
	tests := []struct {
		name    string
		limit   int
		offset  int
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
		check   func(t *testing.T, articles []*domainarticle.Article)
	}{
		{
			name:   "success list articles",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow(1, "Article 1", "Content 1", 1, time.Now(), time.Now()).
					AddRow(2, "Article 2", "Content 2", 1, time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, articles []*domainarticle.Article) {
				assert.Len(t, articles, 2)
				assert.Equal(t, int64(1), articles[0].ID)
				assert.Equal(t, "Article 1", articles[0].Title)
				assert.Equal(t, int64(2), articles[1].ID)
				assert.Equal(t, "Article 2", articles[1].Title)
			},
		},
		{
			name:   "success list articles empty result",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, articles []*domainarticle.Article) {
				if articles == nil {
					articles = []*domainarticle.Article{}
				}
				assert.Len(t, articles, 0)
			},
		},
		{
			name:   "error on database query",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(10, 0).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name:   "error on scan",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow("invalid", "Article 1", "Content 1", 1, time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:   "error on rows iteration",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow(1, "Article 1", "Content 1", 1, time.Now(), time.Now()).
					RowError(0, errors.New("row error"))
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Fatalf("Failed to close database connection: %v", err)
				}
			}()

			repo := NewMySQLRepository(db)
			tt.setup(mock)

			result, err := repo.List(context.Background(), tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				// For empty results, result might be nil or empty slice, both are acceptable
				if result == nil && tt.name == "success list articles empty result" {
					result = []*domainarticle.Article{}
				}
				if tt.check != nil {
					tt.check(t, result)
				} else {
					assert.NotNil(t, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMySQLRepository_ListByAuthor(t *testing.T) {
	tests := []struct {
		name     string
		authorID int64
		limit    int
		offset   int
		setup    func(mock sqlmock.Sqlmock)
		wantErr  bool
		check    func(t *testing.T, articles []*domainarticle.Article)
	}{
		{
			name:     "success list articles by author",
			authorID: 1,
			limit:    10,
			offset:   0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow(1, "Article 1", "Content 1", 1, time.Now(), time.Now()).
					AddRow(2, "Article 2", "Content 2", 1, time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(1, 10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, articles []*domainarticle.Article) {
				assert.Len(t, articles, 2)
				assert.Equal(t, int64(1), articles[0].ID)
				assert.Equal(t, "Article 1", articles[0].Title)
				assert.Equal(t, int64(1), articles[0].AuthorID)
				assert.Equal(t, int64(2), articles[1].ID)
				assert.Equal(t, "Article 2", articles[1].Title)
				assert.Equal(t, int64(1), articles[1].AuthorID)
			},
		},
		{
			name:     "success list articles by author empty result",
			authorID: 1,
			limit:    10,
			offset:   0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(1, 10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, articles []*domainarticle.Article) {
				if articles == nil {
					articles = []*domainarticle.Article{}
				}
				assert.Len(t, articles, 0)
			},
		},
		{
			name:     "error on database query",
			authorID: 1,
			limit:    10,
			offset:   0,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(1, 10, 0).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name:     "error on scan",
			authorID: 1,
			limit:    10,
			offset:   0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow("invalid", "Article 1", "Content 1", 1, time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(1, 10, 0).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:     "error on rows iteration",
			authorID: 1,
			limit:    10,
			offset:   0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at"}).
					AddRow(1, "Article 1", "Content 1", 1, time.Now(), time.Now()).
					RowError(0, errors.New("row error"))
				mock.ExpectQuery("SELECT id, title, content, author_id, created_at, updated_at").
					WithArgs(1, 10, 0).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Fatalf("Failed to close database connection: %v", err)
				}
			}()

			repo := NewMySQLRepository(db)
			tt.setup(mock)

			result, err := repo.ListByAuthor(context.Background(), tt.authorID, tt.limit, tt.offset)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				// For empty results, result might be nil or empty slice, both are acceptable
				if result == nil && tt.name == "success list articles by author empty result" {
					result = []*domainarticle.Article{}
				}
				if tt.check != nil {
					tt.check(t, result)
				} else {
					assert.NotNil(t, result)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMySQLRepository_Count(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(mock sqlmock.Sqlmock)
		want    int64
		wantErr bool
	}{
		{
			name: "success count articles",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).
					AddRow(42)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM articles").
					WillReturnRows(rows)
			},
			want:    42,
			wantErr: false,
		},
		{
			name: "success count zero articles",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).
					AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM articles").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "error on database query",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM articles").
					WillReturnError(errors.New("database error"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Fatalf("Failed to close database connection: %v", err)
				}
			}()

			repo := NewMySQLRepository(db)
			tt.setup(mock)

			result, err := repo.Count(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int64(0), result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMySQLRepository_CountByAuthor(t *testing.T) {
	tests := []struct {
		name     string
		authorID int64
		setup    func(mock sqlmock.Sqlmock)
		want     int64
		wantErr  bool
	}{
		{
			name:     "success count articles by author",
			authorID: 1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).
					AddRow(10)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM articles WHERE author_id = \\?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    10,
			wantErr: false,
		},
		{
			name:     "success count zero articles by author",
			authorID: 1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).
					AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM articles WHERE author_id = \\?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name:     "error on database query",
			authorID: 1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM articles WHERE author_id = \\?").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer func() {
				mock.ExpectClose()
				if err := db.Close(); err != nil {
					t.Fatalf("Failed to close database connection: %v", err)
				}
			}()

			repo := NewMySQLRepository(db)
			tt.setup(mock)

			result, err := repo.CountByAuthor(context.Background(), tt.authorID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int64(0), result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
