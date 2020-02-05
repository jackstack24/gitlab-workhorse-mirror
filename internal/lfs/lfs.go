/*
In this file we handle git lfs objects downloads and uploads
*/

package lfs

import (
	"fmt"
	"net/http"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/config"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/filestore"
)

type object struct {
	size int64
	oid  string
}

func (l *object) Verify(fh *filestore.FileHandler) error {
	if fh.Size != l.size {
		return fmt.Errorf("LFSObject: expected size %d, wrote %d", l.size, fh.Size)
	}

	if fh.SHA256() != l.oid {
		return fmt.Errorf("LFSObject: expected sha256 %s, got %s", l.oid, fh.SHA256())
	}

	return nil
}

type uploadPreparer struct {
	objectStorageConfig config.ObjectStorageConfig
}

func NewLfsUploadPreparer(c config.Config) filestore.UploadPreparer {
	cfg, err := c.FindObjectStorageConfig("lfs")

	if err != nil {
		cfg = &config.ObjectStorageConfig{Enabled: false}
	}

	return &uploadPreparer{objectStorageConfig: *cfg}
}

func (l *uploadPreparer) Prepare(a *api.Response) (*filestore.SaveFileOpts, filestore.UploadVerifier, error) {
	opts := filestore.GetOpts(a)
	opts.TempFilePrefix = a.LfsOid
	opts.ObjectStorageConfig = l.objectStorageConfig

	return opts, &object{oid: a.LfsOid, size: a.LfsSize}, nil
}

func PutStore(a *api.API, h http.Handler, p filestore.UploadPreparer) http.Handler {
	return filestore.BodyUploader(a, h, p)
}
