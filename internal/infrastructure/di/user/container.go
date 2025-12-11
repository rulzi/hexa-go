package user

import (
	"database/sql"

	userdb "github.com/rulzi/hexa-go/internal/adapters/db/user"
	userexternal "github.com/rulzi/hexa-go/internal/adapters/external/user"
	httpuser "github.com/rulzi/hexa-go/internal/adapters/http/user"
	"github.com/rulzi/hexa-go/internal/application/user/usecase"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// UserContainer holds all user domain dependencies
type UserContainer struct {
	Repo          domainuser.Repository
	Service       *domainuser.Service
	EmailSender   usecase.EmailSender
	CreateUseCase *usecase.CreateUserUseCase
	GetUseCase    *usecase.GetUserUseCase
	ListUseCase   *usecase.ListUsersUseCase
	UpdateUseCase *usecase.UpdateUserUseCase
	DeleteUseCase *usecase.DeleteUserUseCase
	LoginUseCase  *usecase.LoginUseCase
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
	createUseCase := usecase.NewCreateUserUseCase(userRepo, userService, emailSender)
	getUseCase := usecase.NewGetUserUseCase(userRepo)
	listUseCase := usecase.NewListUsersUseCase(userRepo)
	updateUseCase := usecase.NewUpdateUserUseCase(userRepo, userService)
	deleteUseCase := usecase.NewDeleteUserUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo, userService)

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
