package di

import (
	"database/sql"

	userdb "github.com/rulzi/hexa-go/internal/adapters/db/user"
	userexternal "github.com/rulzi/hexa-go/internal/adapters/external/user"
	httpuser "github.com/rulzi/hexa-go/internal/adapters/http/user"
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
	Handler       *httpuser.Handler
}

// NewUserContainer creates a new user domain container
func NewUserContainer(database *sql.DB, jwtSecret string, jwtExpiration int) *UserContainer {
	// Initialize repository (driven adapter)
	userRepo := userdb.NewMySQLRepository(database)

	// Initialize domain service
	userService := domainuser.NewService(userRepo, jwtSecret, jwtExpiration)

	// Initialize external service adapter
	emailSender := userexternal.NewEmailSenderImpl()

	// Initialize use cases (application layer)
	createUseCase := appuser.NewCreateUserUseCase(userRepo, userService, emailSender)
	getUseCase := appuser.NewGetUserUseCase(userRepo)
	listUseCase := appuser.NewListUsersUseCase(userRepo)
	updateUseCase := appuser.NewUpdateUserUseCase(userRepo, userService)
	deleteUseCase := appuser.NewDeleteUserUseCase(userRepo)
	loginUseCase := appuser.NewLoginUseCase(userRepo, userService)

	// Initialize HTTP handler (driving adapter)
	userHandler := httpuser.NewHandler(
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
