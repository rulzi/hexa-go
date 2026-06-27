package user

import (
	"database/sql"

	authadapter "github.com/rulzi/hexa-go/internal/adapters/auth"
	userexternal "github.com/rulzi/hexa-go/internal/adapters/external/user"
	httpuser "github.com/rulzi/hexa-go/internal/adapters/http/user"
	userdb "github.com/rulzi/hexa-go/internal/adapters/repository/user"
	"github.com/rulzi/hexa-go/internal/application/user/usecase"
	userport "github.com/rulzi/hexa-go/internal/domain/user/port"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// Container holds all user domain dependencies
type Container struct {
	Repo                userport.Repository
	TokenGen            userport.TokenGenerator
	TokenValidator      userport.TokenValidator
	PasswordHasher      userport.PasswordHasher
	NotificationService userport.NotificationService
	Handler             *httpuser.Handler
}

// NewContainer creates a new user domain container
func NewContainer(database *sql.DB, appLogger logger.Logger, jwtSecret string, jwtExpiration int) *Container {
	userRepo := userdb.NewMySQLRepository(database, appLogger)

	jwtAdapter := authadapter.NewJWTAdapter(jwtSecret, jwtExpiration)
	passwordHasher := authadapter.NewBcryptPasswordHasher()

	notificationService := userexternal.NewEmailSenderImpl(appLogger)

	userHandler := httpuser.NewHandler(httpuser.Deps{
		Create: usecase.NewCreateUser(userRepo, passwordHasher, notificationService),
		Get:    usecase.NewGetUser(userRepo),
		List:   usecase.NewListUser(userRepo),
		Update: usecase.NewUpdateUser(userRepo, passwordHasher),
		Delete: usecase.NewDeleteUser(userRepo),
		Login:  usecase.NewLoginUser(userRepo, passwordHasher, jwtAdapter),
	}, appLogger)

	return &Container{
		Repo:                userRepo,
		TokenGen:            jwtAdapter,
		TokenValidator:      jwtAdapter,
		PasswordHasher:      passwordHasher,
		NotificationService: notificationService,
		Handler:             userHandler,
	}
}
