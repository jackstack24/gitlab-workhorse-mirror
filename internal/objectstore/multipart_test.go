package objectstore_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/objectstore"
	"gitlab.com/gitlab-org/gitlab-workhorse/internal/objectstore/test"
)

func TestMultipartUploadWithUpcaseETags(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var putCnt, postCnt int

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()

		// Part upload request
		if r.Method == "PUT" {
			putCnt++

			w.Header().Set("ETag", strings.ToUpper(test.ObjectMD5))
		}

		// POST with CompleteMultipartUpload request
		if r.Method == "POST" {
			expectedETag := "6e6b164c392b04bfbb82368179d9ade2-1"
			completeBody := fmt.Sprintf(`<CompleteMultipartUploadResult>
                                                       <Bucket>test-bucket</Bucket>
			                               <ETag>%s</ETag>
  	                                             </CompleteMultipartUploadResult>`,
				strings.ToUpper(expectedETag))
			postCnt++

			w.Write([]byte(completeBody))
		}
	}))
	defer ts.Close()

	deadline := time.Now().Add(testTimeout)

	m, err := objectstore.NewMultipart(ctx,
		[]string{ts.URL},    // a single presigned part URL
		ts.URL,              // the complete multipart upload URL
		"",                  // no abort
		"",                  // no delete
		map[string]string{}, // no custom headers
		deadline,
		test.ObjectSize, // parts size equal to the whole content. Only 1 part
		false)
	require.NoError(t, err)

	_, err = m.Write([]byte(test.ObjectContent))
	require.NoError(t, err)
	require.NoError(t, m.Close())
	require.Equal(t, 1, putCnt, "1 part expected")
	require.Equal(t, 1, postCnt, "1 complete multipart upload expected")
}

func TestMultipartUploadWithEtagVerifyDisabled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var putCnt, postCnt int

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()

		// Part upload request
		if r.Method == "PUT" {
			putCnt++
			require.Equal(t, "/part", r.URL.Path)
			w.Header().Set("ETag", "not a valid MD5 hash")
		}

		// POST with CompleteMultipartUpload request
		if r.Method == "POST" {
			postCnt++
			require.Equal(t, "/complete", r.URL.Path)
			completeBody := `<CompleteMultipartUploadResult>
                             <Bucket>test-bucket</Bucket>
                             <ETag>not a valid MD5-based ETag"</ETag>
                             </CompleteMultipartUploadResult>`
			w.Write([]byte(completeBody))
		} else if r.Method == "DELETE" {
			require.Equal(t, "/abort", r.URL.Path)
			t.FailNow()
		}

	}))
	defer ts.Close()

	deadline := time.Now().Add(testTimeout)
	var headers = map[string]string{}

	m, err := objectstore.NewMultipart(ctx,
		[]string{ts.URL + "/part"},
		ts.URL+"/complete",
		ts.URL+"/abort",
		ts.URL+"/delete",
		headers,
		deadline,
		test.ObjectSize,
		true)
	require.NoError(t, err)

	_, err = m.Write([]byte(test.ObjectContent))
	require.NoError(t, err)
	require.NoError(t, m.Close())
	require.Equal(t, 1, putCnt, "1 part expected")
	require.Equal(t, 1, postCnt, "1 complete multipart upload expected")
}
