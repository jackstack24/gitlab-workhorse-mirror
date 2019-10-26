package objectstore

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"gitlab.com/gitlab-org/labkit/log"
	"gitlab.com/gitlab-org/labkit/mask"
)

// ErrNotEnoughParts will be used when writing more than size * len(partURLs)
var ErrNotEnoughParts = errors.New("not enough Parts")

// Multipart represents a MultipartUpload on a S3 compatible Object Store service.
// It can be used as io.WriteCloser for uploading an object
type Multipart struct {
	// CompleteURL is a presigned URL for CompleteMultipartUpload
	CompleteURL string
	// AbortURL is a presigned URL for AbortMultipartUpload
	AbortURL string
	// DeleteURL is a presigned URL for RemoveObject
	DeleteURL string

	// By default, the MD5 sum of the file is expected to match the
	// ETag. If this is not true for the object storage provider or
	// transfer mode, this should be disabled by an upstream config.
	SkipETagVerify bool

	uploader
}

// NewMultipart provides Multipart pointer that can be used for uploading. Data written will be split buffered on disk up to size bytes
// then uploaded with S3 Upload Part. Once Multipart is Closed a final call to CompleteMultipartUpload will be sent.
// In case of any error a call to AbortMultipartUpload will be made to cleanup all the resources
func NewMultipart(ctx context.Context, partURLs []string, completeURL, abortURL, deleteURL string, putHeaders map[string]string, deadline time.Time, partSize int64, skipETagVerify bool) (*Multipart, error) {
	pr, pw := io.Pipe()
	uploadCtx, cancelFn := context.WithDeadline(ctx, deadline)
	m := &Multipart{
		CompleteURL:    completeURL,
		AbortURL:       abortURL,
		DeleteURL:      deleteURL,
		SkipETagVerify: skipETagVerify,
		uploader:       newUploader(uploadCtx, pw),
	}

	go m.trackUploadTime()
	go m.cleanup(ctx)

	objectStorageUploadsOpen.Inc()

	go func() {
		defer cancelFn()
		defer objectStorageUploadsOpen.Dec()
		defer func() {
			// This will be returned as error to the next write operation on the pipe
			pr.CloseWithError(m.uploadError)
		}()

		cmu := &CompleteMultipartUpload{}
		for i, partURL := range partURLs {
			src := io.LimitReader(pr, partSize)
			part, err := m.readAndUploadOnePart(partURL, putHeaders, src, i+1)
			if err != nil {
				m.uploadError = err
				return
			}
			if part == nil {
				break
			} else {
				cmu.Part = append(cmu.Part, part)
			}
		}

		n, err := io.Copy(ioutil.Discard, pr)
		if err != nil {
			m.uploadError = fmt.Errorf("drain pipe: %v", err)
			return
		}
		if n > 0 {
			m.uploadError = ErrNotEnoughParts
			return
		}

		if err := m.complete(cmu); err != nil {
			m.uploadError = err
			return
		}
	}()

	return m, nil
}

func (m *Multipart) trackUploadTime() {
	started := time.Now()
	<-m.ctx.Done()
	objectStorageUploadTime.Observe(time.Since(started).Seconds())
}

func (m *Multipart) cleanup(ctx context.Context) {
	// wait for the upload to finish
	<-m.ctx.Done()

	if m.uploadError != nil {
		objectStorageUploadRequestsRequestFailed.Inc()
		m.abort()
		return
	}

	// We have now successfully uploaded the file to object storage. Another
	// goroutine will hand off the object to gitlab-rails.
	<-ctx.Done()

	// gitlab-rails is now done with the object so it's time to delete it.
	m.delete()
}

func (m *Multipart) complete(cmu *CompleteMultipartUpload) error {
	body, err := xml.Marshal(cmu)
	if err != nil {
		return fmt.Errorf("marshal CompleteMultipartUpload request: %v", err)
	}

	req, err := http.NewRequest("POST", m.CompleteURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create CompleteMultipartUpload request: %v", err)
	}
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Type", "application/xml")
	req = req.WithContext(m.ctx)

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("CompleteMultipartUpload request %q: %v", mask.URL(m.CompleteURL), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CompleteMultipartUpload request %v returned: %s", mask.URL(m.CompleteURL), resp.Status)
	}

	result := &compoundCompleteMultipartUploadResult{}
	decoder := xml.NewDecoder(resp.Body)
	if err := decoder.Decode(&result); err != nil {
		return fmt.Errorf("decode CompleteMultipartUpload answer: %v", err)
	}

	if result.isError() {
		return result
	}

	if result.CompleteMultipartUploadResult == nil {
		return fmt.Errorf("empty CompleteMultipartUploadResult")
	}

	m.extractETag(result.ETag)

	return m.verifyETag(cmu)
}

func (m *Multipart) verifyETag(cmu *CompleteMultipartUpload) error {
	if m.SkipETagVerify {
		return nil
	}

	expectedChecksum, err := cmu.BuildMultipartUploadETag()
	if err != nil {
		return err
	}

	return compareMD5(expectedChecksum, m.etag)
}

func (m *Multipart) readAndUploadOnePart(partURL string, putHeaders map[string]string, src io.Reader, partNumber int) (*completeMultipartUploadPart, error) {
	file, err := ioutil.TempFile("", "part-buffer")
	if err != nil {
		return nil, fmt.Errorf("create temporary buffer file: %v", err)
	}
	defer func(path string) {
		if err := os.Remove(path); err != nil {
			log.WithError(err).WithField("file", path).Warning("Unable to delete temporary file")
		}
	}(file.Name())

	n, err := io.Copy(file, src)
	if err != nil {
		return nil, fmt.Errorf("write part %d to disk: %v", partNumber, err)
	}
	if n == 0 {
		return nil, nil
	}

	if _, err = file.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("rewind part %d temporary dump : %v", partNumber, err)
	}

	etag, err := m.uploadPart(partURL, putHeaders, file, n)
	if err != nil {
		return nil, fmt.Errorf("upload part %d: %v", partNumber, err)
	}
	return &completeMultipartUploadPart{PartNumber: partNumber, ETag: etag}, nil
}

func (m *Multipart) uploadPart(url string, headers map[string]string, body io.Reader, size int64) (string, error) {
	deadline, ok := m.ctx.Deadline()
	if !ok {
		return "", fmt.Errorf("missing deadline")
	}

	part, err := newObject(m.ctx, url, "", headers, deadline, size, false, m.SkipETagVerify)
	if err != nil {
		return "", err
	}

	_, err = io.CopyN(part, body, size)
	if err != nil {
		return "", err
	}

	err = part.Close()
	if err != nil {
		return "", err
	}

	return part.ETag(), nil
}

func (m *Multipart) delete() {
	m.syncAndDelete(m.DeleteURL)
}

func (m *Multipart) abort() {
	m.syncAndDelete(m.AbortURL)
}
