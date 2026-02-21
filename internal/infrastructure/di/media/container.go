package media

import (
	"database/sql"

	mediadb "github.com/rulzi/hexa-go/internal/adapters/repository/media"
	httpmedia "github.com/rulzi/hexa-go/internal/adapters/http/media"
	mediastorage "github.com/rulzi/hexa-go/internal/adapters/storage/media"
	"github.com/rulzi/hexa-go/internal/application/media/usecase"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// Container holds all media domain dependencies
type Container struct {
	Repo      domainmedia.Repository
	Storage   domainmedia.Storage
	Service   *domainmedia.Service
	MediaUC   *usecase.MediaUseCase
	Handler   *httpmedia.Handler
}

// NewContainer creates a new media domain container
func NewContainer(database *sql.DB, storageBasePath string, baseURL string) (*Container, error) {
	// Initialize repository (driven adapter)
	mediaRepo := mediadb.NewMySQLRepository(database)

	// Initialize storage (driven adapter)
	storage, err := mediastorage.NewLocalStorage(storageBasePath)
	if err != nil {
		return nil, err
	}

	// Initialize domain service
	mediaService := domainmedia.NewService(mediaRepo)

	// Initialize use case (application layer)
	mediaUseCase := usecase.NewMediaUseCase(mediaRepo, mediaService, storage, baseURL)

	// Initialize HTTP handler (driving adapter)
	mediaHandler := httpmedia.NewHandler(mediaUseCase)

	return &Container{
		Repo:    mediaRepo,
		Storage: storage,
		Service: mediaService,
		MediaUC: mediaUseCase,
		Handler: mediaHandler,
	}, nil
}
