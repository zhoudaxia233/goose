package goose

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	g := New()
	go func() {
		g.GET("/goose", func(ctx *Context) {
			ctx.String("I love goose!")
		})
		if err := g.Run(":8080"); err != nil {
			t.Fatal(err)
		}
	}()
	// wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(200 * time.Millisecond)

	url := "http://127.0.0.1:8080" + "/goose"
	testResponseBody(t, url, http.StatusOK, "I love goose!")
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
