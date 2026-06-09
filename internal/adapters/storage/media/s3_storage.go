package media

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	domainmedia "github.com/rulzi/hexa-go/internal/domain/media"
)

// S3Config holds configuration for the S3 storage adapter.
type S3Config struct {
	Bucket       string
	Region       string
	Endpoint     string
	UsePathStyle bool
}

// S3StorageAdapter stores files in an S3-compatible object store.
type S3StorageAdapter struct {
	client *s3.Client
	bucket string
}

// NewS3StorageAdapter creates an S3 storage adapter using the AWS SDK default credential chain.
func NewS3StorageAdapter(cfg S3Config) (*S3StorageAdapter, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("S3 bucket is required")
	}
	if cfg.Region == "" {
		return nil, fmt.Errorf("S3 region is required")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(), awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.Endpoint != "" {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
		}
		o.UsePathStyle = cfg.UsePathStyle
	})

	return &S3StorageAdapter{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

// Save uploads a file and returns the object key.
func (s *S3StorageAdapter) Save(ctx context.Context, filename string, file io.Reader) (string, error) {
	key := generateStoragePath(filename)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return key, nil
}

// Delete removes an object by key. Idempotent when the object does not exist.
func (s *S3StorageAdapter) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(normalizeS3Key(path)),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// Get retrieves an object by key.
func (s *S3StorageAdapter) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	output, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(normalizeS3Key(path)),
	})
	if err != nil {
		if isS3NotFound(err) {
			return nil, domainmedia.NewMediaNotFound()
		}
		return nil, fmt.Errorf("failed to get file from S3: %w", err)
	}

	return output.Body, nil
}

func normalizeS3Key(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func isS3NotFound(err error) bool {
	var notFound *types.NotFound
	var noSuchKey *types.NoSuchKey
	return errors.As(err, &notFound) || errors.As(err, &noSuchKey)
}
