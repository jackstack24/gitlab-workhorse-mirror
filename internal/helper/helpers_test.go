package helper

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixRemoteAddr(t *testing.T) {
	testCases := []struct {
		initial   string
		forwarded string
		expected  string
	}{
		{initial: "@", forwarded: "", expected: "127.0.0.1:0"},
		{initial: "@", forwarded: "18.245.0.1", expected: "18.245.0.1:0"},
		{initial: "@", forwarded: "127.0.0.1", expected: "127.0.0.1:0"},
		{initial: "@", forwarded: "192.168.0.1", expected: "127.0.0.1:0"},
		{initial: "192.168.1.1:0", forwarded: "", expected: "192.168.1.1:0"},
		{initial: "192.168.1.1:0", forwarded: "18.245.0.1", expected: "18.245.0.1:0"},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest("POST", "unix:///tmp/test.socket/info/refs", nil)
		require.NoError(t, err)

		req.RemoteAddr = tc.initial

		if tc.forwarded != "" {
			req.Header.Add("X-Forwarded-For", tc.forwarded)
		}

		FixRemoteAddr(req)

		assert.Equal(t, tc.expected, req.RemoteAddr)
	}
}

func TestSetForwardedForGeneratesHeader(t *testing.T) {
	testCases := []struct {
		remoteAddr           string
		previousForwardedFor []string
		expected             string
	}{
		{
			"8.8.8.8:3000",
			nil,
			"8.8.8.8",
		},
		{
			"8.8.8.8:3000",
			[]string{"138.124.33.63, 151.146.211.237"},
			"138.124.33.63, 151.146.211.237, 8.8.8.8",
		},
		{
			"8.8.8.8:3000",
			[]string{"8.154.76.107", "115.206.118.179"},
			"8.154.76.107, 115.206.118.179, 8.8.8.8",
		},
	}
	for _, tc := range testCases {
		headers := http.Header{}
		originalRequest := http.Request{
			RemoteAddr: tc.remoteAddr,
		}

		if tc.previousForwardedFor != nil {
			originalRequest.Header = http.Header{
				"X-Forwarded-For": tc.previousForwardedFor,
			}
		}

		SetForwardedFor(&headers, &originalRequest)

		result := headers.Get("X-Forwarded-For")
		if result != tc.expected {
			t.Fatalf("Expected %v, got %v", tc.expected, result)
		}
	}
}

func TestReadRequestBody(t *testing.T) {
	data := []byte("123456")
	rw := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(data))

	result, err := ReadRequestBody(rw, req, 1000)
	assert.NoError(t, err)
	assert.Equal(t, data, result)
}

func TestReadRequestBodyLimit(t *testing.T) {
	data := []byte("123456")
	rw := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(data))

	_, err := ReadRequestBody(rw, req, 2)
	assert.Error(t, err)
}

func TestCloneRequestWithBody(t *testing.T) {
	input := []byte("test")
	newInput := []byte("new body")
	req, _ := http.NewRequest("POST", "/test", bytes.NewBuffer(input))
	newReq := CloneRequestWithNewBody(req, newInput)

	assert.NotEqual(t, req, newReq)
	assert.NotEqual(t, req.Body, newReq.Body)
	assert.NotEqual(t, len(newInput), newReq.ContentLength)

	var buffer bytes.Buffer
	io.Copy(&buffer, newReq.Body)
	assert.Equal(t, newInput, buffer.Bytes())
}

func TestApplicationJson(t *testing.T) {
	req, _ := http.NewRequest("POST", "/test", nil)
	req.Header.Set("Content-Type", "application/json")

	assert.True(t, IsApplicationJson(req), "expected to match 'application/json' as 'application/json'")

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	assert.True(t, IsApplicationJson(req), "expected to match 'application/json; charset=utf-8' as 'application/json'")

	req.Header.Set("Content-Type", "text/plain")
	assert.False(t, IsApplicationJson(req), "expected not to match 'text/plain' as 'application/json'")
}

func TestFail500WorksWithNils(t *testing.T) {
	body := bytes.NewBuffer(nil)
	w := httptest.NewRecorder()
	w.Body = body

	Fail500(w, nil, nil)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal server error\n", body.String())
}

func TestScrubURLParams(t *testing.T) {
	for before, expected := range map[string]string{
		"http://example.com":                                                "http://example.com",
		"http://example.com?foo=1":                                          "http://example.com?foo=1",
		"http://example.com?title=token":                                    "http://example.com?title=token",
		"http://example.com?authenticity_token=1":                           "http://example.com?authenticity_token=[FILTERED]",
		"http://example.com?private_token=1":                                "http://example.com?private_token=[FILTERED]",
		"http://example.com?rss_token=1":                                    "http://example.com?rss_token=[FILTERED]",
		"http://example.com?access_token=1":                                 "http://example.com?access_token=[FILTERED]",
		"http://example.com?refresh_token=1":                                "http://example.com?refresh_token=[FILTERED]",
		"http://example.com?foo&authenticity_token=blahblah&bar":            "http://example.com?foo&authenticity_token=[FILTERED]&bar",
		"http://example.com?private-token=1":                                "http://example.com?private-token=[FILTERED]",
		"http://example.com?foo&private-token=blahblah&bar":                 "http://example.com?foo&private-token=[FILTERED]&bar",
		"http://example.com?private-token=foo&authenticity_token=bar":       "http://example.com?private-token=[FILTERED]&authenticity_token=[FILTERED]",
		"https://example.com:8080?private-token=foo&authenticity_token=bar": "https://example.com:8080?private-token=[FILTERED]&authenticity_token=[FILTERED]",
		"/?private-token=foo&authenticity_token=bar":                        "/?private-token=[FILTERED]&authenticity_token=[FILTERED]",
		"?private-token=&authenticity_token=&bar":                           "?private-token=[FILTERED]&authenticity_token=[FILTERED]&bar",
		"?private-token=foo&authenticity_token=bar":                         "?private-token=[FILTERED]&authenticity_token=[FILTERED]",
		"?private_token=foo&authenticity-token=bar":                         "?private_token=[FILTERED]&authenticity-token=[FILTERED]",
		"?X-AMZ-Signature=foo":                                              "?X-AMZ-Signature=[FILTERED]",
		"?x-amz-signature=foo":                                              "?x-amz-signature=[FILTERED]",
		"?Signature=foo":                                                    "?Signature=[FILTERED]",
		"?confirmation_password=foo":                                        "?confirmation_password=[FILTERED]",
		"?pos_secret_number=foo":                                            "?pos_secret_number=[FILTERED]",
		"?sharedSecret=foo":                                                 "?sharedSecret=[FILTERED]",
		"?book_key=foo":                                                     "?book_key=[FILTERED]",
		"?certificate=foo":                                                  "?certificate=[FILTERED]",
		"?hook=foo":                                                         "?hook=[FILTERED]",
		"?import_url=foo":                                                   "?import_url=[FILTERED]",
		"?otp_attempt=foo":                                                  "?otp_attempt=[FILTERED]",
		"?sentry_dsn=foo":                                                   "?sentry_dsn=[FILTERED]",
		"?trace=foo":                                                        "?trace=[FILTERED]",
		"?variables=foo":                                                    "?variables=[FILTERED]",
		"?content=foo":                                                      "?content=[FILTERED]",
		"?content=e=mc2":                                                    "?content=[FILTERED]",
		"?formula=e=mc2":                                                    "?formula=e=mc2",
		"http://%41:8080/":                                                  "<invalid URL>",
	} {
		after := ScrubURLParams(before)
		assert.Equal(t, expected, after, "Scrubbing %q", before)
	}
}

func TestDecodeMd5Checksum(t *testing.T) {
	testCases := []struct {
		etag  string
		valid bool
	}{
		{etag: "64eba6b86e926702235156b6ebbe6932", valid: true},
		{etag: "64eba6b86e926702235156b6ebbe6932-1", valid: false},
		{etag: "64eba6b86e926702235156b6ebbe693", valid: false},
		{etag: "64eba6b86e926702235156b6ebbe69Z", valid: false},
		{etag: "64eba", valid: false},
	}

	for _, tc := range testCases {
		_, err := DecodeMd5Checksum(tc.etag)

		if tc.valid {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
