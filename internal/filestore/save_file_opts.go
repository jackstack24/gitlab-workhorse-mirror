package filestore

import (
	"strings"
	"time"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/api"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/config"
)

// DefaultObjectStoreTimeout is the timeout for ObjectStore upload operation
const DefaultObjectStoreTimeout = 4 * time.Hour

type ObjectStorageConfig struct {
	Provider string

	S3Credentials config.S3Credentials `toml:"s3"`
	S3Config      config.S3Config
}

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
	// Whether to ignore Rails pre-signed URLs and have Workhorse directly access object storage provider
	UseWorkhorseClient bool
	// If UseWorkhorseClient is true, this is the temporary object name to store the file
	RemoteTempObjectID string
	// Workhorse object storage client (e.g. S3) parameters
	ObjectStorageConfig ObjectStorageConfig
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
}

// UseWorkhorseClientEnabled checks if the options require direct access to object storage
func (s *SaveFileOpts) UseWorkhorseClientEnabled() bool {
	return s.UseWorkhorseClient && s.ObjectStorageConfig.IsValid() && s.RemoteTempObjectID != ""
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
		LocalTempPath:      apiResponse.TempPath,
		RemoteID:           apiResponse.RemoteObject.ID,
		RemoteURL:          apiResponse.RemoteObject.GetURL,
		PresignedPut:       apiResponse.RemoteObject.StoreURL,
		PresignedDelete:    apiResponse.RemoteObject.DeleteURL,
		PutHeaders:         apiResponse.RemoteObject.PutHeaders,
		UseWorkhorseClient: apiResponse.RemoteObject.UseWorkhorseClient,
		RemoteTempObjectID: apiResponse.RemoteObject.RemoteTempObjectID,
		Deadline:           time.Now().Add(timeout),
	}

	objectStorageParams := apiResponse.RemoteObject.ObjectStorage
	if opts.UseWorkhorseClient && objectStorageParams != nil {
		opts.ObjectStorageConfig.Provider = objectStorageParams.Provider
		opts.ObjectStorageConfig.S3Config = objectStorageParams.S3Config
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

func (c *ObjectStorageConfig) IsAWS() bool {
	return strings.EqualFold(c.Provider, "AWS") || strings.EqualFold(c.Provider, "S3")
}

func (c *ObjectStorageConfig) IsValid() bool {
	return c.S3Config.Bucket != "" && c.S3Config.Region != ""
}
