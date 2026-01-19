package media

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
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
		media   *domainmedia.Media
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
		check   func(t *testing.T, media *domainmedia.Media)
	}{
		{
			name: "success create media",
			media: &domainmedia.Media{
				Name:      "test-image.jpg",
				Path:      "/storage/2025/12/19/test-image.jpg",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO media").
					WithArgs("test-image.jpg", "/storage/2025/12/19/test-image.jpg", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
			check: func(t *testing.T, media *domainmedia.Media) {
				assert.Equal(t, int64(1), media.ID)
				assert.Equal(t, "test-image.jpg", media.Name)
				assert.Equal(t, "/storage/2025/12/19/test-image.jpg", media.Path)
			},
		},
		{
			name: "error on database exec",
			media: &domainmedia.Media{
				Name:      "test-image.jpg",
				Path:      "/storage/2025/12/19/test-image.jpg",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO media").
					WithArgs("test-image.jpg", "/storage/2025/12/19/test-image.jpg", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "error on last insert id",
			media: &domainmedia.Media{
				Name:      "test-image.jpg",
				Path:      "/storage/2025/12/19/test-image.jpg",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO media").
					WithArgs("test-image.jpg", "/storage/2025/12/19/test-image.jpg", sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			result, err := repo.Create(context.Background(), tt.media)

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
		check   func(t *testing.T, media *domainmedia.Media)
	}{
		{
			name: "success get media by id",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "path", "created_at", "updated_at"}).
					AddRow(1, "test-image.jpg", "/storage/2025/12/19/test-image.jpg", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, path, created_at, updated_at").
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, media *domainmedia.Media) {
				assert.Equal(t, int64(1), media.ID)
				assert.Equal(t, "test-image.jpg", media.Name)
				assert.Equal(t, "/storage/2025/12/19/test-image.jpg", media.Path)
			},
		},
		{
			name: "media not found",
			id:   999,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, path, created_at, updated_at").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			check: func(t *testing.T, media *domainmedia.Media) {
				assert.Nil(t, media)
			},
		},
		{
			name: "database error",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, path, created_at, updated_at").
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
				if tt.name == "media not found" {
					assert.Equal(t, domainmedia.ErrMediaNotFound, err)
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
		media   *domainmedia.Media
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
		check   func(t *testing.T, media *domainmedia.Media)
	}{
		{
			name: "success update media",
			media: &domainmedia.Media{
				ID:        1,
				Name:      "updated-image.jpg",
				Path:      "/storage/2025/12/19/updated-image.jpg",
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE media").
					WithArgs("updated-image.jpg", "/storage/2025/12/19/updated-image.jpg", sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
			check: func(t *testing.T, media *domainmedia.Media) {
				assert.Equal(t, int64(1), media.ID)
				assert.Equal(t, "updated-image.jpg", media.Name)
				assert.Equal(t, "/storage/2025/12/19/updated-image.jpg", media.Path)
			},
		},
		{
			name: "error on database exec",
			media: &domainmedia.Media{
				ID:        1,
				Name:      "updated-image.jpg",
				Path:      "/storage/2025/12/19/updated-image.jpg",
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE media").
					WithArgs("updated-image.jpg", "/storage/2025/12/19/updated-image.jpg", sqlmock.AnyArg(), 1).
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

			result, err := repo.Update(context.Background(), tt.media)

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
			name: "success delete media",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM media").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "media not found",
			id:   999,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM media").
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "error on database exec",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM media").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "error on rows affected",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM media").
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
				if tt.name == "media not found" {
					assert.Equal(t, domainmedia.ErrMediaNotFound, err)
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
		check   func(t *testing.T, mediaList []*domainmedia.Media)
	}{
		{
			name:   "success list media",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "path", "created_at", "updated_at"}).
					AddRow(1, "image1.jpg", "/storage/2025/12/19/image1.jpg", time.Now(), time.Now()).
					AddRow(2, "image2.jpg", "/storage/2025/12/19/image2.jpg", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, path, created_at, updated_at").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, mediaList []*domainmedia.Media) {
				assert.Len(t, mediaList, 2)
				assert.Equal(t, int64(1), mediaList[0].ID)
				assert.Equal(t, "image1.jpg", mediaList[0].Name)
				assert.Equal(t, "/storage/2025/12/19/image1.jpg", mediaList[0].Path)
				assert.Equal(t, int64(2), mediaList[1].ID)
				assert.Equal(t, "image2.jpg", mediaList[1].Name)
				assert.Equal(t, "/storage/2025/12/19/image2.jpg", mediaList[1].Path)
			},
		},
		{
			name:   "success list media empty result",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "path", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, name, path, created_at, updated_at").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, mediaList []*domainmedia.Media) {
				if mediaList == nil {
					mediaList = []*domainmedia.Media{}
				}
				assert.Len(t, mediaList, 0)
			},
		},
		{
			name:   "error on database query",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, path, created_at, updated_at").
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
				rows := sqlmock.NewRows([]string{"id", "name", "path", "created_at", "updated_at"}).
					AddRow("invalid", "image1.jpg", "/storage/2025/12/19/image1.jpg", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, path, created_at, updated_at").
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
				rows := sqlmock.NewRows([]string{"id", "name", "path", "created_at", "updated_at"}).
					AddRow(1, "image1.jpg", "/storage/2025/12/19/image1.jpg", time.Now(), time.Now()).
					RowError(0, errors.New("row error"))
				mock.ExpectQuery("SELECT id, name, path, created_at, updated_at").
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
				if result == nil && tt.name == "success list media empty result" {
					result = []*domainmedia.Media{}
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
			name: "success count media",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).
					AddRow(42)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM media").
					WillReturnRows(rows)
			},
			want:    42,
			wantErr: false,
		},
		{
			name: "success count zero media",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).
					AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM media").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "error on database query",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM media").
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
