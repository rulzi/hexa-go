package di

import (
	"database/sql"

	"github.com/rulzi/hexa-go/internal/adapters/db"
	"github.com/rulzi/hexa-go/internal/adapters/external"
	"github.com/rulzi/hexa-go/internal/adapters/http"
	appuser "github.com/rulzi/hexa-go/internal/application/user"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// UserContainer holds all user domain dependencies
type UserContainer struct {
	Repo          domainuser.Repository
	Service       *domainuser.Service
	EmailSender   appuser.EmailSender
	CreateUseCase *appuser.CreateUserUseCase
	GetUseCase    *appuser.GetUserUseCase
	ListUseCase   *appuser.ListUsersUseCase
	UpdateUseCase *appuser.UpdateUserUseCase
	DeleteUseCase *appuser.DeleteUserUseCase
	LoginUseCase  *appuser.LoginUseCase
	Handler       *http.UserHandler
}

// NewUserContainer creates a new user domain container
func NewUserContainer(database *sql.DB, jwtSecret string, jwtExpiration int) *UserContainer {
	// Initialize repository (driven adapter)
	userRepo := db.NewUserMySQLRepository(database)

	// Initialize domain service
	userService := domainuser.NewService(userRepo, jwtSecret, jwtExpiration)

	// Initialize external service adapter
	emailSender := external.NewEmailSenderImpl()

	// Initialize use cases (application layer)
	createUseCase := appuser.NewCreateUserUseCase(userRepo, userService, emailSender)
	getUseCase := appuser.NewGetUserUseCase(userRepo)
	listUseCase := appuser.NewListUsersUseCase(userRepo)
	updateUseCase := appuser.NewUpdateUserUseCase(userRepo, userService)
	deleteUseCase := appuser.NewDeleteUserUseCase(userRepo)
	loginUseCase := appuser.NewLoginUseCase(userRepo, userService)

	// Initialize HTTP handler (driving adapter)
	userHandler := http.NewUserHandler(
		createUseCase,
		getUseCase,
		listUseCase,
		updateUseCase,
		deleteUseCase,
		loginUseCase,
	)

	return &UserContainer{
		Repo:          userRepo,
		Service:       userService,
		EmailSender:   emailSender,
		CreateUseCase: createUseCase,
		GetUseCase:    getUseCase,
		ListUseCase:   listUseCase,
		UpdateUseCase: updateUseCase,
		DeleteUseCase: deleteUseCase,
		LoginUseCase:  loginUseCase,
		Handler:       userHandler,
	}
}
