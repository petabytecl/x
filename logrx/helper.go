package logrx

import (
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"gitlab.endpoints.autentiax-ctl.cloud.goog/golang/x/tracex"
)

func WithRequest(l logr.Logger, r *http.Request) logr.Logger {
	headers := map[string]interface{}{}
	if ua := r.UserAgent(); len(ua) > 0 {
		headers["user-agent"] = ua
	}

	if cookie := r.Header.Get("Cookie"); cookie != "" {
		headers["cookie"] = cookie
	}

	if auth := r.Header.Get("Authorization"); auth != "" {
		headers["authorization"] = auth
	}

	for _, key := range []string{"Referer", "Origin", "Accept", "X-Request-ID", "If-None-Match",
		"X-Forwarded-For", "X-Forwarded-Proto", "Cache-Control", "Accept-Encoding", "Accept-Language", "If-Modified-Since"} {
		if value := r.Header.Get(key); len(value) > 0 {
			headers[strings.ToLower(key)] = value
		}
	}

	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}

	sc := tracex.GetSpan(r)

	ll := l.WithName("http_request").WithValues(
		"remote", r.RemoteAddr,
		"method", r.Method,
		"path", r.URL.EscapedPath(),
		"query", r.URL.RawQuery,
		"scheme", scheme,
		"host", r.Host,
		"headers", headers,
		"trace_id", sc.TraceID,
		"span_id", sc.SpanID,
	)

	return ll
}
