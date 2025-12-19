package media

import (
	"database/sql"

	mediadb "github.com/rulzi/hexa-go/internal/adapters/db/media"
	httpmedia "github.com/rulzi/hexa-go/internal/adapters/http/media"
	mediastorage "github.com/rulzi/hexa-go/internal/adapters/storage/media"
	"github.com/rulzi/hexa-go/internal/application/media/usecase"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// Container holds all media domain dependencies
type Container struct {
	Repo          domainmedia.Repository
	Storage       domainmedia.Storage
	Service       *domainmedia.Service
	CreateUseCase *usecase.CreateMediaUseCase
	GetUseCase    *usecase.GetMediaUseCase
	ListUseCase   *usecase.ListMediaUseCase
	UpdateUseCase *usecase.UpdateMediaUseCase
	DeleteUseCase *usecase.DeleteMediaUseCase
	Handler       *httpmedia.Handler
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

	// Initialize use cases (application layer)
	createMediaUseCase := usecase.NewCreateMediaUseCase(mediaRepo, mediaService, storage, baseURL)
	getMediaUseCase := usecase.NewGetMediaUseCase(mediaRepo, baseURL)
	listMediaUseCase := usecase.NewListMediaUseCase(mediaRepo, baseURL)
	updateMediaUseCase := usecase.NewUpdateMediaUseCase(mediaRepo, mediaService, storage, baseURL)
	deleteMediaUseCase := usecase.NewDeleteMediaUseCase(mediaRepo, storage)

	// Initialize HTTP handler (driving adapter)
	mediaHandler := httpmedia.NewHandler(
		createMediaUseCase,
		getMediaUseCase,
		listMediaUseCase,
		updateMediaUseCase,
		deleteMediaUseCase,
	)

	return &Container{
		Repo:          mediaRepo,
		Storage:       storage,
		Service:       mediaService,
		CreateUseCase: createMediaUseCase,
		GetUseCase:    getMediaUseCase,
		ListUseCase:   listMediaUseCase,
		UpdateUseCase: updateMediaUseCase,
		DeleteUseCase: deleteMediaUseCase,
		Handler:       mediaHandler,
	}, nil
}
