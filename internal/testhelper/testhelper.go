package testhelper

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"regexp"
	"runtime"
	"strings"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"

	"gitlab.com/gitlab-org/labkit/log"

	"gitlab.com/gitlab-org/gitlab-workhorse/internal/secret"
)

func ConfigureSecret() {
	secret.SetPath(path.Join(RootDir(), "testdata/test-secret"))
}

var extractPatchSeriesMatcher = regexp.MustCompile(`^From (\w+)`)

// AssertPatchSeries takes a `git format-patch` blob, extracts the From xxxxx
// lines and compares the SHAs to expected list.
func AssertPatchSeries(t *testing.T, blob []byte, expected ...string) {
	var actual []string
	footer := make([]string, 3)

	scanner := bufio.NewScanner(bytes.NewReader(blob))

	for scanner.Scan() {
		line := scanner.Text()
		if matches := extractPatchSeriesMatcher.FindStringSubmatch(line); len(matches) == 2 {
			actual = append(actual, matches[1])
		}
		footer = []string{footer[1], footer[2], line}
	}

	if strings.Join(actual, "\n") != strings.Join(expected, "\n") {
		t.Fatalf("Patch series differs. Expected: %v. Got: %v", expected, actual)
	}

	// Check the last returned patch is complete
	// Don't assert on the final line, it is a git version
	if footer[0] != "-- " {
		t.Fatalf("Expected end of patch, found: \n\t%q", strings.Join(footer, "\n\t"))
	}
}

func AssertResponseCode(t *testing.T, response *httptest.ResponseRecorder, expectedCode int) {
	if response.Code != expectedCode {
		t.Fatalf("for HTTP request expected to get %d, got %d instead", expectedCode, response.Code)
	}
}

func AssertResponseBody(t *testing.T, response *httptest.ResponseRecorder, expectedBody string) {
	if response.Body.String() != expectedBody {
		t.Fatalf("for HTTP request expected to receive %q, got %q instead as body", expectedBody, response.Body.String())
	}
}

func AssertResponseBodyRegexp(t *testing.T, response *httptest.ResponseRecorder, expectedBody *regexp.Regexp) {
	if !expectedBody.MatchString(response.Body.String()) {
		t.Fatalf("for HTTP request expected to receive body matching %q, got %q instead", expectedBody.String(), response.Body.String())
	}
}

func AssertResponseWriterHeader(t *testing.T, w http.ResponseWriter, header string, expected ...string) {
	actual := w.Header()[http.CanonicalHeaderKey(header)]

	assertHeaderExists(t, header, actual, expected)
}

func AssertAbsentResponseWriterHeader(t *testing.T, w http.ResponseWriter, header string) {
	actual := w.Header()[http.CanonicalHeaderKey(header)]

	if len(actual) != 0 {
		t.Fatalf("for HTTP request expected not to receive the header %q, got %+v", header, actual)
	}
}

func AssertResponseHeader(t *testing.T, w interface{}, header string, expected ...string) {
	var actual []string

	header = http.CanonicalHeaderKey(header)

	if resp, ok := w.(*http.Response); ok {
		actual = resp.Header[header]
	} else if resp, ok := w.(http.ResponseWriter); ok {
		actual = resp.Header()[header]
	} else if resp, ok := w.(*httptest.ResponseRecorder); ok {
		actual = resp.Header()[header]
	} else {
		t.Fatalf("invalid type of w passed AssertResponseHeader")
	}

	assertHeaderExists(t, header, actual, expected)
}

func assertHeaderExists(t *testing.T, header string, actual, expected []string) {
	if len(expected) != len(actual) {
		t.Fatalf("for HTTP request expected to receive the header %q with %+v, got %+v", header, expected, actual)
	}

	for i, value := range expected {
		if value != actual[i] {
			t.Fatalf("for HTTP request expected to receive the header %q with %+v, got %+v", header, expected, actual)
		}
	}
}

func TestServerWithHandler(url *regexp.Regexp, handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logEntry := log.WithFields(log.Fields{
			"method": r.Method,
			"url":    r.URL,
			"action": "DENY",
		})

		if url != nil && !url.MatchString(r.URL.Path) {
			logEntry.Info("UPSTREAM")
			w.WriteHeader(404)
			return
		}

		if version := r.Header.Get("Gitlab-Workhorse"); version == "" {
			logEntry.Info("UPSTREAM")
			w.WriteHeader(403)
			return
		}

		handler(w, r)
	}))
}

var workhorseExecutables = []string{"gitlab-workhorse", "gitlab-zip-cat", "gitlab-zip-metadata"}

func BuildExecutables() error {
	rootDir := RootDir()

	for _, exe := range workhorseExecutables {
		if _, err := os.Stat(path.Join(rootDir, exe)); os.IsNotExist(err) {
			return fmt.Errorf("cannot find executable %s. Please run 'make prepare-tests'", exe)
		}
	}

	oldPath := os.Getenv("PATH")
	testPath := fmt.Sprintf("%s:%s", rootDir, oldPath)
	if err := os.Setenv("PATH", testPath); err != nil {
		return fmt.Errorf("failed to set PATH to %v", testPath)
	}

	return nil
}

func RootDir() string {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic(errors.New("RootDir: calling runtime.Caller failed"))
	}
	return path.Join(path.Dir(currentFile), "../..")
}

func LoadFile(t *testing.T, filePath string) string {
	content, err := ioutil.ReadFile(path.Join(RootDir(), filePath))
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func ParseJWT(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}

	ConfigureSecret()
	secretBytes, err := secret.Bytes()
	if err != nil {
		return nil, fmt.Errorf("read secret from file: %v", err)
	}

	return secretBytes, nil
}
