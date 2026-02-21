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
	Repo                domainuser.Repository
	Service             *domainuser.Service
	TokenGen            domainuser.TokenGenerator
	TokenValidator      domainuser.TokenValidator
	PasswordHasher      domainuser.PasswordHasher
	NotificationService domainuser.NotificationService
	UserUC              *usecase.UserUseCase
	Handler             *httpuser.Handler
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

	// Initialize use case (application layer)
	userUseCase := usecase.NewUserUseCase(userRepo, passwordHasher, notificationService, jwtAdapter)

	// Initialize HTTP handler (driving adapter)
	userHandler := httpuser.NewHandler(userUseCase)

	return &Container{
		Repo:                userRepo,
		Service:             userService,
		TokenGen:            jwtAdapter,
		TokenValidator:      jwtAdapter,
		PasswordHasher:      passwordHasher,
		NotificationService: notificationService,
		UserUC:              userUseCase,
		Handler:             userHandler,
	}
}
