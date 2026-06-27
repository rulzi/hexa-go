package media

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/stretchr/testify/assert"
)

func TestNewS3StorageAdapter_MissingBucket(t *testing.T) {
	_, err := NewS3StorageAdapter(S3Config{Region: "us-east-1"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "bucket")
}

func TestNewS3StorageAdapter_MissingRegion(t *testing.T) {
	_, err := NewS3StorageAdapter(S3Config{Bucket: "my-bucket"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "region")
}

func TestNormalizeS3Key(t *testing.T) {
	assert.Equal(t, "2025/01/01/file.jpg", normalizeS3Key(`2025\01\01\file.jpg`))
}

func TestIsS3NotFound(t *testing.T) {
	assert.True(t, isS3NotFound(&types.NotFound{}))
	assert.True(t, isS3NotFound(&types.NoSuchKey{}))
	assert.False(t, isS3NotFound(errors.New("other error")))
}
