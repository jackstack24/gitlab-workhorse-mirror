package helper

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
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

func TestLogError(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		uri         string
		err         error
		logMatchers []string
	}{
		{
			name: "nil_request",
			err:  fmt.Errorf("crash"),
			logMatchers: []string{
				`level=error msg="unknown error" error=crash`,
			},
		},
		{
			name: "nil_request_nil_error",
			err:  nil,
			logMatchers: []string{
				`level=error msg="unknown error" error="<nil>"`,
			},
		},
		{
			name:   "basic_url",
			method: "GET",
			uri:    "http://localhost:3000/",
			err:    fmt.Errorf("error"),
			logMatchers: []string{
				`level=error msg=error correlation_id= error=error method=GET uri="http://localhost:3000/"`,
			},
		},
		{
			name:   "secret_url",
			method: "GET",
			uri:    "http://localhost:3000/path?certificate=123&sharedSecret=123&import_url=the_url&my_password_string=password",
			err:    fmt.Errorf("error"),
			logMatchers: []string{
				`level=error msg=error correlation_id= error=error method=GET uri="http://localhost:3000/path\?certificate=\[FILTERED\]&sharedSecret=\[FILTERED\]&import_url=\[FILTERED\]&my_password_string=\[FILTERED\]"`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			oldOut := logrus.StandardLogger().Out
			logrus.StandardLogger().Out = buf
			defer func() {
				logrus.StandardLogger().Out = oldOut
			}()

			var r *http.Request
			if tt.uri != "" {
				r = httptest.NewRequest(tt.method, tt.uri, nil)
			}

			LogError(r, tt.err)

			logString := buf.String()
			for _, v := range tt.logMatchers {
				require.Regexp(t, v, logString)
			}
		})
	}
}

func TestLogErrorWithFields(t *testing.T) {
	tests := []struct {
		name       string
		request    *http.Request
		err        error
		fields     map[string]interface{}
		logMatcher string
	}{
		{
			name:       "nil_request",
			err:        fmt.Errorf("crash"),
			fields:     map[string]interface{}{"extra_one": 123},
			logMatcher: `level=error msg="unknown error" error=crash extra_one=123`,
		},
		{
			name:       "nil_request_nil_error",
			err:        nil,
			fields:     map[string]interface{}{"extra_one": 123, "extra_two": "test"},
			logMatcher: `level=error msg="unknown error" error="<nil>" extra_one=123 extra_two=test`,
		},
		{
			name:       "basic_url",
			request:    httptest.NewRequest("GET", "http://localhost:3000/", nil),
			err:        fmt.Errorf("error"),
			fields:     map[string]interface{}{"extra_one": 123, "extra_two": "test"},
			logMatcher: `level=error msg=error correlation_id= error=error extra_one=123 extra_two=test method=GET uri="http://localhost:3000/`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}

			oldOut := logrus.StandardLogger().Out
			logrus.StandardLogger().Out = buf
			defer func() {
				logrus.StandardLogger().Out = oldOut
			}()

			LogErrorWithFields(tt.request, tt.err, tt.fields)

			logString := buf.String()
			require.Contains(t, logString, tt.logMatcher)
		})
	}
}
