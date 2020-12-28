package buildx

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

func (i Info) HTTPHandler() http.Handler {
	var body []byte
	body, _ = jsoniter.Marshal(i)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(body)
	})
}
