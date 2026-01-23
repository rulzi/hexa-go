package user

import (
	"database/sql"

	authadapter "github.com/rulzi/hexa-go/internal/adapters/auth"
	userdb "github.com/rulzi/hexa-go/internal/adapters/repository/user"
	userexternal "github.com/rulzi/hexa-go/internal/adapters/external/user"
	httpuser "github.com/rulzi/hexa-go/internal/adapters/http/user"
	"github.com/rulzi/hexa-go/internal/application/user/usecase"
	domainuser "github.com/rulzi/hexa-go/internal/domain/user"
)

// Container holds all user domain dependencies
type Container struct {
	Repo              domainuser.Repository
	Service           *domainuser.Service
	TokenGen          domainuser.TokenGenerator
	TokenValidator    domainuser.TokenValidator
	PasswordHasher    domainuser.PasswordHasher
	NotificationService domainuser.NotificationService
	CreateUseCase     *usecase.CreateUserUseCase
	GetUseCase        *usecase.GetUserUseCase
	ListUseCase       *usecase.ListUsersUseCase
	UpdateUseCase     *usecase.UpdateUserUseCase
	DeleteUseCase     *usecase.DeleteUserUseCase
	LoginUseCase      *usecase.LoginUseCase
	Handler           *httpuser.Handler
}

// NewContainer creates a new user domain container
func NewContainer(database *sql.DB, jwtSecret string, jwtExpiration int) *Container {
	// Initialize repository (driven adapter)
	userRepo := userdb.NewMySQLRepository(database)

	// Initialize auth adapters (driven adapters)
	jwtAdapter := authadapter.NewJWTAdapter(jwtSecret, jwtExpiration)
	passwordHasher := authadapter.NewBcryptPasswordHasher()

	// Initialize domain service
	userService := domainuser.NewService(userRepo, jwtAdapter, jwtAdapter, passwordHasher)

	// Initialize external service adapter
	notificationService := userexternal.NewEmailSenderImpl()

	// Initialize use cases (application layer)
	createUseCase := usecase.NewCreateUserUseCase(userRepo, passwordHasher, notificationService)
	getUseCase := usecase.NewGetUserUseCase(userRepo)
	listUseCase := usecase.NewListUsersUseCase(userRepo)
	updateUseCase := usecase.NewUpdateUserUseCase(userRepo, passwordHasher)
	deleteUseCase := usecase.NewDeleteUserUseCase(userRepo)
	loginUseCase := usecase.NewLoginUseCase(userRepo, passwordHasher, jwtAdapter)

	// Initialize HTTP handler (driving adapter)
	userHandler := httpuser.NewHandler(
		createUseCase,
		getUseCase,
		listUseCase,
		updateUseCase,
		deleteUseCase,
		loginUseCase,
	)

	return &Container{
		Repo:                userRepo,
		Service:             userService,
		TokenGen:            jwtAdapter,
		TokenValidator:      jwtAdapter,
		PasswordHasher:      passwordHasher,
		NotificationService: notificationService,
		CreateUseCase:       createUseCase,
		GetUseCase:          getUseCase,
		ListUseCase:         listUseCase,
		UpdateUseCase:       updateUseCase,
		DeleteUseCase:       deleteUseCase,
		LoginUseCase:        loginUseCase,
		Handler:             userHandler,
	}
}
