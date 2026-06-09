package media

import (
	"database/sql"
	"fmt"

	httpmedia "github.com/rulzi/hexa-go/internal/adapters/http/media"
	mediadb "github.com/rulzi/hexa-go/internal/adapters/repository/media"
	mediastorage "github.com/rulzi/hexa-go/internal/adapters/storage/media"
	"github.com/rulzi/hexa-go/internal/application/media/usecase"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
	"github.com/rulzi/hexa-go/internal/infrastructure/config"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// Container holds all media domain dependencies
type Container struct {
	Repo      domainmedia.Repository
	Storage   domainmedia.StoragePort
	Service   *domainmedia.Service
	MediaUC   *usecase.MediaUseCase
	Handler   *httpmedia.Handler
}

// NewContainer creates a new media domain container
func NewContainer(database *sql.DB, appLogger logger.Logger, storageCfg config.StorageConfig) (*Container, error) {
	mediaRepo := mediadb.NewMySQLRepository(database)

	storage, err := newStorageAdapter(storageCfg)
	if err != nil {
		return nil, err
	}

	mediaService := domainmedia.NewService(mediaRepo)
	mediaUseCase := usecase.NewMediaUseCase(mediaRepo, mediaService, storage, storageCfg.BaseURL)
	mediaHandler := httpmedia.NewHandler(mediaUseCase, appLogger)

	return &Container{
		Repo:    mediaRepo,
		Storage: storage,
		Service: mediaService,
		MediaUC: mediaUseCase,
		Handler: mediaHandler,
	}, nil
}

func newStorageAdapter(storageCfg config.StorageConfig) (domainmedia.StoragePort, error) {
	switch storageCfg.Driver {
	case "s3":
		return mediastorage.NewS3StorageAdapter(mediastorage.S3Config{
			Bucket:       storageCfg.S3Bucket,
			Region:       storageCfg.S3Region,
			Endpoint:     storageCfg.S3Endpoint,
			UsePathStyle: storageCfg.S3PathStyle,
		})
	case "local", "":
		return mediastorage.NewLocalStorageAdapter(storageCfg.BasePath)
	default:
		return nil, fmt.Errorf("unsupported STORAGE_DRIVER %q", storageCfg.Driver)
	}
}
