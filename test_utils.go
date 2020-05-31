package goose

import (
	"net/http"
	"net/http/httptest"
)

func sendRequest(mux http.Handler, method string, path string) (w *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return
}
