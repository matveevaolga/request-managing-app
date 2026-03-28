package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func ExecuteRequest(req *http.Request, handler http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

func CreateJSONBody(body interface{}) (*bytes.Reader, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBody), nil
}
