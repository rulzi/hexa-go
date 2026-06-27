package media

import (
	"database/sql"
	"fmt"

	httpmedia "github.com/rulzi/hexa-go/internal/adapters/http/media"
	mediadb "github.com/rulzi/hexa-go/internal/adapters/repository/media"
	mediastorage "github.com/rulzi/hexa-go/internal/adapters/storage/media"
	"github.com/rulzi/hexa-go/internal/application/media/usecase"
	mediaport "github.com/rulzi/hexa-go/internal/domain/media/port"
	"github.com/rulzi/hexa-go/internal/infrastructure/config"
	"github.com/rulzi/hexa-go/internal/infrastructure/logger"
)

// Container holds all media domain dependencies
type Container struct {
	Repo    mediaport.Repository
	Storage mediaport.StoragePort
	Handler *httpmedia.Handler
}

// NewContainer creates a new media domain container
func NewContainer(database *sql.DB, appLogger logger.Logger, storageCfg config.StorageConfig) (*Container, error) {
	mediaRepo := mediadb.NewMySQLRepository(database, appLogger)

	storage, err := newStorageAdapter(storageCfg, appLogger)
	if err != nil {
		return nil, err
	}

	mediaHandler := httpmedia.NewHandler(httpmedia.Deps{
		Create: usecase.NewCreateMedia(mediaRepo, storage, storageCfg.BaseURL),
		Get:    usecase.NewGetMedia(mediaRepo, storageCfg.BaseURL),
		List:   usecase.NewListMedia(mediaRepo, storageCfg.BaseURL),
		Update: usecase.NewUpdateMedia(mediaRepo, storage, storageCfg.BaseURL),
		Delete: usecase.NewDeleteMedia(mediaRepo, storage),
	}, appLogger)

	return &Container{
		Repo:    mediaRepo,
		Storage: storage,
		Handler: mediaHandler,
	}, nil
}

func newStorageAdapter(storageCfg config.StorageConfig, appLogger logger.Logger) (mediaport.StoragePort, error) {
	switch storageCfg.Driver {
	case "s3":
		return mediastorage.NewS3StorageAdapter(mediastorage.S3Config{
			Bucket:       storageCfg.S3Bucket,
			Region:       storageCfg.S3Region,
			Endpoint:     storageCfg.S3Endpoint,
			UsePathStyle: storageCfg.S3PathStyle,
		})
	case "local", "":
		return mediastorage.NewLocalStorageAdapter(storageCfg.BasePath, appLogger)
	default:
		return nil, fmt.Errorf("unsupported STORAGE_DRIVER %q", storageCfg.Driver)
	}
}
