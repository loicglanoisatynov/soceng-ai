package sessions

import (
	"bytes"
	"io"
	"net/http"
)

func Start_challenge(r *http.Request) http.Response {
	var status_code int
	payload := []byte(`{"message": "default payload"}`)

	return http.Response{
		StatusCode: status_code,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBuffer(payload)),
	}
}
