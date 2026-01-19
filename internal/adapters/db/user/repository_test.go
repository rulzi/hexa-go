package user

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
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
		user    *domainuser.User
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
		check   func(t *testing.T, user *domainuser.User)
	}{
		{
			name: "success create user",
			user: &domainuser.User{
				Name:      "John Doe",
				Email:     "john@example.com",
				Password:  "hashedpassword",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("John Doe", "john@example.com", "hashedpassword", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
			check: func(t *testing.T, user *domainuser.User) {
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john@example.com", user.Email)
			},
		},
		{
			name: "error on database exec",
			user: &domainuser.User{
				Name:      "John Doe",
				Email:     "john@example.com",
				Password:  "hashedpassword",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("John Doe", "john@example.com", "hashedpassword", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "error on last insert id",
			user: &domainuser.User{
				Name:      "John Doe",
				Email:     "john@example.com",
				Password:  "hashedpassword",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("INSERT INTO users").
					WithArgs("John Doe", "john@example.com", "hashedpassword", sqlmock.AnyArg(), sqlmock.AnyArg()).
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

			result, err := repo.Create(context.Background(), tt.user)

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
		check   func(t *testing.T, user *domainuser.User)
	}{
		{
			name: "success get user by id",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john@example.com", "hashedpassword", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
					WithArgs(1).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, user *domainuser.User) {
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john@example.com", user.Email)
				assert.Equal(t, "hashedpassword", user.Password)
			},
		},
		{
			name: "user not found",
			id:   999,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			check: func(t *testing.T, user *domainuser.User) {
				assert.Nil(t, user)
			},
		},
		{
			name: "database error",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
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
				if tt.name == "user not found" {
					assert.Equal(t, domainuser.ErrUserNotFound, err)
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

func TestMySQLRepository_GetByEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
		check   func(t *testing.T, user *domainuser.User)
	}{
		{
			name:  "success get user by email",
			email: "john@example.com",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john@example.com", "hashedpassword", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
					WithArgs("john@example.com").
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, user *domainuser.User) {
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, "John Doe", user.Name)
				assert.Equal(t, "john@example.com", user.Email)
				assert.Equal(t, "hashedpassword", user.Password)
			},
		},
		{
			name:  "user not found",
			email: "notfound@example.com",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
					WithArgs("notfound@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: true,
			check: func(t *testing.T, user *domainuser.User) {
				assert.Nil(t, user)
			},
		},
		{
			name:  "database error",
			email: "john@example.com",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
					WithArgs("john@example.com").
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

			result, err := repo.GetByEmail(context.Background(), tt.email)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.name == "user not found" {
					assert.Equal(t, domainuser.ErrUserNotFound, err)
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
		user    *domainuser.User
		setup   func(mock sqlmock.Sqlmock)
		wantErr bool
		check   func(t *testing.T, user *domainuser.User)
	}{
		{
			name: "success update user",
			user: &domainuser.User{
				ID:        1,
				Name:      "John Updated",
				Email:     "john.updated@example.com",
				Password:  "newhashedpassword",
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users").
					WithArgs("John Updated", "john.updated@example.com", "newhashedpassword", sqlmock.AnyArg(), 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
			check: func(t *testing.T, user *domainuser.User) {
				assert.Equal(t, int64(1), user.ID)
				assert.Equal(t, "John Updated", user.Name)
				assert.Equal(t, "john.updated@example.com", user.Email)
			},
		},
		{
			name: "error on database exec",
			user: &domainuser.User{
				ID:        1,
				Name:      "John Updated",
				Email:     "john.updated@example.com",
				Password:  "newhashedpassword",
				UpdatedAt: time.Now(),
			},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE users").
					WithArgs("John Updated", "john.updated@example.com", "newhashedpassword", sqlmock.AnyArg(), 1).
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

			result, err := repo.Update(context.Background(), tt.user)

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
			name: "success delete user",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users").
					WithArgs(1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "user not found",
			id:   999,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users").
					WithArgs(999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
		{
			name: "error on database exec",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users").
					WithArgs(1).
					WillReturnError(errors.New("database error"))
			},
			wantErr: true,
		},
		{
			name: "error on rows affected",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("DELETE FROM users").
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
				if tt.name == "user not found" {
					assert.Equal(t, domainuser.ErrUserNotFound, err)
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
		check   func(t *testing.T, users []*domainuser.User)
	}{
		{
			name:   "success list users",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john@example.com", "hashedpassword", time.Now(), time.Now()).
					AddRow(2, "Jane Doe", "jane@example.com", "hashedpassword2", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, users []*domainuser.User) {
				assert.Len(t, users, 2)
				assert.Equal(t, int64(1), users[0].ID)
				assert.Equal(t, "John Doe", users[0].Name)
				assert.Equal(t, int64(2), users[1].ID)
				assert.Equal(t, "Jane Doe", users[1].Name)
			},
		},
		{
			name:   "success list users empty result",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
			wantErr: false,
			check: func(t *testing.T, users []*domainuser.User) {
				if users == nil {
					users = []*domainuser.User{}
				}
				assert.Len(t, users, 0)
			},
		},
		{
			name:   "error on database query",
			limit:  10,
			offset: 0,
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
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
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow("invalid", "John Doe", "john@example.com", "hashedpassword", time.Now(), time.Now())
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
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
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john@example.com", "hashedpassword", time.Now(), time.Now()).
					RowError(0, errors.New("row error"))
				mock.ExpectQuery("SELECT id, name, email, password, created_at, updated_at").
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
				if result == nil && tt.name == "success list users empty result" {
					result = []*domainuser.User{}
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
			name: "success count users",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).
					AddRow(42)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
					WillReturnRows(rows)
			},
			want:    42,
			wantErr: false,
		},
		{
			name: "success count zero users",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).
					AddRow(0)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
					WillReturnRows(rows)
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "error on database query",
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users").
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
