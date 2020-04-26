package filestore

import (
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/config"
)

type ObjectStoragePreparer struct {
	credentials config.ObjectStorageCredentials
}

func NewObjectStoragePreparer(c config.Config) UploadPreparer {
	creds := c.ObjectStorageCredentials

	if creds == nil {
		creds = &config.ObjectStorageCredentials{}
	}

	return &ObjectStoragePreparer{credentials: *creds}
}

func (p *ObjectStoragePreparer) Prepare(a *api.Response) (*SaveFileOpts, UploadVerifier, error) {
	opts := GetOpts(a)
	opts.ObjectStorageConfig.S3Credentials = p.credentials.S3Credentials

	return opts, nil, nil
}
