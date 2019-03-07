package sentry

import (
	"net/http"
	"reflect"

	"github.com/getsentry/raven-go"

	correlation "gitlab.com/gitlab-org/labkit/correlation/raven"
)

var ravenHeaderBlacklist = []string{
	"Authorization",
	"Private-Token",
}

func CaptureError(r *http.Request, err error) {
	client := raven.DefaultClient
	extra := raven.Extra{}

	interfaces := []raven.Interface{}
	if r != nil {
		cleanHeadersForRaven(r)
		interfaces = append(interfaces, raven.NewHttp(r))

		extra = correlation.SetExtra(r.Context(), extra)
	}

	exception := &raven.Exception{
		Stacktrace: raven.NewStacktrace(2, 3, nil),
		Value:      err.Error(),
		Type:       reflect.TypeOf(err).String(),
	}
	interfaces = append(interfaces, exception)

	packet := raven.NewPacketWithExtra(err.Error(), extra, interfaces...)
	client.Capture(packet, nil)
}

func cleanHeadersForRaven(r *http.Request) {
	if r == nil {
		return
	}

	for _, key := range ravenHeaderBlacklist {
		if r.Header.Get(key) != "" {
			r.Header.Set(key, "[redacted]")
		}
	}
}
