package goose

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"strings"
	"testing"
)

func sendRequest(mux http.Handler, method string, path string) (w *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, path, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	return
}

func testResponseBody(t *testing.T, url string, wantStatusCode int, wantBody string) {
	// https://stackoverflow.com/questions/12122159/how-to-do-a-https-request-with-bad-certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	resp, err := client.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, ioerr := ioutil.ReadAll(resp.Body)
	if ioerr != nil {
		t.Fatal(ioerr)
	}

	if resp.StatusCode != wantStatusCode {
		t.Errorf("Status Code: Want %d, Got %d", wantStatusCode, resp.StatusCode)
	}

	if string(body) != wantBody {
		t.Errorf("Text doesn't match. Want: %s, Got: %s", wantBody, string(body))
	}
}

func testResponseBodyLocal(t *testing.T, g *Goose, pattern string, wantStatusCode int, wantBody string) {
	ts := httptest.NewServer(g)
	defer ts.Close()

	resp, err := http.Get(joinURL(ts.URL, pattern))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, ioerr := ioutil.ReadAll(resp.Body)
	if ioerr != nil {
		t.Fatal(ioerr)
	}

	if resp.StatusCode != wantStatusCode {
		t.Errorf("Status Code: Want %d, Got %d", wantStatusCode, resp.StatusCode)
	}

	if string(body) != wantBody {
		t.Errorf("Text doesn't match. Want: %s, Got: %s", wantBody, string(body))
	}
}

// https://stackoverflow.com/a/34668130
func joinURL(base, pattern string) string {
	u, _ := url.Parse(base)
	u.Path = path.Join(u.Path, pattern)

	if len(pattern) > 1 && strings.HasSuffix(pattern, "/") {
		return u.String() + "/"
	}
	return u.String()
}
