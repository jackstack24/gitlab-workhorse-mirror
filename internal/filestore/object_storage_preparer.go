package filestore

import (
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/config"
)

type ObjectStoragePreparer struct {
	objectStorageConfig config.ObjectStorageConfig
}

func NewObjectStoragePreparer(name string, c config.Config) UploadPreparer {
	cfg, err := c.FindObjectStorageConfig(name)

	if err != nil {
		cfg = &config.ObjectStorageConfig{Enabled: false}
	}

	return &ObjectStoragePreparer{objectStorageConfig: *cfg}
}

func (p *ObjectStoragePreparer) Prepare(a *api.Response) (*SaveFileOpts, UploadVerifier, error) {
	opts := GetOpts(a)
	opts.ObjectStorageConfig = p.objectStorageConfig

	return opts, nil, nil
}
