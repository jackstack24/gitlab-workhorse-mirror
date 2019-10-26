package filestore

import (
	"time"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/api"
)

// DefaultObjectStoreTimeout is the timeout for ObjectStore upload operation
const DefaultObjectStoreTimeout = 4 * time.Hour

// SaveFileOpts represents all the options available for saving a file to object store
type SaveFileOpts struct {
	// TempFilePrefix is the prefix used to create temporary local file
	TempFilePrefix string
	// LocalTempPath is the directory where to write a local copy of the file
	LocalTempPath string
	// RemoteID is the remote ObjectID provided by GitLab
	RemoteID string
	// RemoteURL is the final URL of the file
	RemoteURL string
	// PresignedPut is a presigned S3 PutObject compatible URL
	PresignedPut string
	// PresignedDelete is a presigned S3 DeleteObject compatible URL.
	PresignedDelete string
	// HTTP headers to be sent along with PUT request
	PutHeaders map[string]string

	// Deadline it the S3 operation deadline, the upload will be aborted if not completed in time
	Deadline time.Time

	//MultipartUpload parameters
	// PartSize is the exact size of each uploaded part. Only the last one can be smaller
	PartSize int64
	// PresignedParts contains the presigned URLs for each part
	PresignedParts []string
	// PresignedCompleteMultipart is a presigned URL for CompleteMulipartUpload
	PresignedCompleteMultipart string
	// PresignedAbortMultipart is a presigned URL for AbortMultipartUpload
	PresignedAbortMultipart string
	// By default, the MD5 sum of the file is expected to match the
	// ETag. If this is not true for the object storage provider or
	// transfer mode, this should be disabled by an upstream config.
	SkipETagVerify bool
}

// IsLocal checks if the options require the writing of the file on disk
func (s *SaveFileOpts) IsLocal() bool {
	return s.LocalTempPath != ""
}

// IsRemote checks if the options requires a remote upload
func (s *SaveFileOpts) IsRemote() bool {
	return s.PresignedPut != "" || s.IsMultipart()
}

// IsMultipart checks if the options requires a Multipart upload
func (s *SaveFileOpts) IsMultipart() bool {
	return s.PartSize > 0
}

// GetOpts converts GitLab api.Response to a proper SaveFileOpts
func GetOpts(apiResponse *api.Response) *SaveFileOpts {
	timeout := time.Duration(apiResponse.RemoteObject.Timeout) * time.Second
	if timeout == 0 {
		timeout = DefaultObjectStoreTimeout
	}

	opts := SaveFileOpts{
		LocalTempPath:   apiResponse.TempPath,
		RemoteID:        apiResponse.RemoteObject.ID,
		RemoteURL:       apiResponse.RemoteObject.GetURL,
		PresignedPut:    apiResponse.RemoteObject.StoreURL,
		PresignedDelete: apiResponse.RemoteObject.DeleteURL,
		PutHeaders:      apiResponse.RemoteObject.PutHeaders,
		SkipETagVerify:  apiResponse.RemoteObject.SkipETagVerify,
		Deadline:        time.Now().Add(timeout),
	}

	// Backwards compatibility to ensure API servers that do not include the
	// CustomPutHeaders flag will default to the original content type.
	if !apiResponse.RemoteObject.CustomPutHeaders {
		opts.PutHeaders = make(map[string]string)
		opts.PutHeaders["Content-Type"] = "application/octet-stream"
	}

	if multiParams := apiResponse.RemoteObject.MultipartUpload; multiParams != nil {
		opts.PartSize = multiParams.PartSize
		opts.PresignedCompleteMultipart = multiParams.CompleteURL
		opts.PresignedAbortMultipart = multiParams.AbortURL
		opts.PresignedParts = append([]string(nil), multiParams.PartURLs...)
	}

	return &opts
}
