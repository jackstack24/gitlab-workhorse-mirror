package filestore_test

import (
	"testing"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/config"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/filestore"

	"github.com/stretchr/testify/require"
)

func TestPrepareWithS3Config(t *testing.T) {
	os := make(map[string]*config.ObjectStorageConfig)
	cfg := &config.ObjectStorageConfig{
		Enabled:  true,
		Provider: "AWS",
		S3Config: config.S3Config{},
	}
	os["uploads"] = cfg

	c := config.Config{
		ObjectStorages: os,
	}

	r := &api.Response{}
	p := filestore.NewObjectStoragePreparer("uploads", c)
	opts, v, err := p.Prepare(r)

	require.NoError(t, err)
	require.True(t, opts.ObjectStorageConfig.Enabled)
	require.True(t, opts.ObjectStorageConfig.IsAWS())
	require.Equal(t, *cfg, opts.ObjectStorageConfig)
	require.Equal(t, nil, v)
}

func TestPrepareWithNoConfig(t *testing.T) {
	c := config.Config{}
	r := &api.Response{}
	p := filestore.NewObjectStoragePreparer("unknown", c)
	opts, v, err := p.Prepare(r)

	require.NoError(t, err)
	require.False(t, opts.ObjectStorageConfig.Enabled)
	require.Equal(t, nil, v)
}
