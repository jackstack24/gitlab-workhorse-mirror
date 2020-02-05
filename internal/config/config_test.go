package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadObjectStorageConfig(t *testing.T) {
	config := `[object_storage]

[object_storage.uploads]
enabled = true
provider = "AWS"

[object_storage.uploads.s3]
aws_access_key_id = "minio"
aws_secret_access_key = "gdk-minio"
region = "gdk"
path_style = true
endpoint = 'http://127.0.0.1:9000'
bucket = "uploads"
`
	tmpFile, err := ioutil.TempFile(os.TempDir(), "test-")
	require.NoError(t, err)

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(config))
	require.NoError(t, err)

	cfg, err := LoadConfig(tmpFile.Name())
	require.NoError(t, err)

	require.Equal(t, 1, len(cfg.ObjectStorages), "Expected one (1) object storages")

	o := cfg.ObjectStorages["uploads"]
	require.True(t, o.IsAWS())
	require.True(t, o.IsValid())

	expected := ObjectStorageConfig{
		Enabled:  true,
		Provider: "AWS",
		S3Config: S3Config{
			AwsAccessKeyID:     "minio",
			AwsSecretAccessKey: "gdk-minio",
			Region:             "gdk",
			PathStyle:          true,
			Endpoint:           "http://127.0.0.1:9000",
			Bucket:             "uploads",
		},
	}

	require.Equal(t, expected, *o)

	_, err = cfg.FindObjectStorageConfig("nothing")
	require.Error(t, err)

	uploadConfig, err := cfg.FindObjectStorageConfig("uploads")
	require.NoError(t, err)
	require.Equal(t, expected, *uploadConfig)
}

func TestValidS3Config(t *testing.T) {
	cfg := ObjectStorageConfig{
		Enabled:  true,
		Provider: "AWS",
		S3Config: S3Config{
			AwsAccessKeyID:     "minio",
			AwsSecretAccessKey: "gdk-minio",
			Region:             "gdk",
			PathStyle:          true,
			Endpoint:           "http://127.0.0.1:9000",
			Bucket:             "uploads",
		},
	}

	require.True(t, cfg.IsValid())
	cfg.S3Config.Region = ""
	require.False(t, cfg.IsValid())

	cfg.S3Config.Region = "gdk"
	cfg.S3Config.Bucket = ""
	require.False(t, cfg.IsValid())
}

func TestDuplicateObjectStorageConfig(t *testing.T) {
	config := `[object_storage]

[object_storage.uploads]
enabled = true
provider = "AWS"

[object_storage.uploads.s3]
aws_access_key_id = "minio"
aws_secret_access_key = "gdk-minio"
region = "gdk"
path_style = true
endpoint = 'http://127.0.0.1:9000'
bucket = "uploads"

[object_storage.uploads]
enabled = true
provider = "S3"

[object_storage.uploads.s3]
aws_access_key_id = "minio"
aws_secret_access_key = "gdk-minio"
region = "gdk"
path_style = true
endpoint = 'http://127.0.0.1:9000'
bucket = "uploads2"

`
	tmpFile, err := ioutil.TempFile(os.TempDir(), "test-")
	require.NoError(t, err)

	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(config))
	require.NoError(t, err)

	_, err = LoadConfig(tmpFile.Name())
	require.Error(t, err)
	require.Contains(t, err.Error(), "Key 'object_storage.uploads' has already been defined.")
}
