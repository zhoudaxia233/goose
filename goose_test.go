package goose

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

const (
	baseURL = "http://127.0.0.1"
)

func TestRun(t *testing.T) {
	port := ":8080"
	g := New()
	go func() {
		g.GET("/", func(ctx *Context) {
			ctx.String("Main Page!")
		})
		g.GET("/goose", func(ctx *Context) {
			name := "goose"
			ctx.String("I love %s!", name)
		})
		if err := g.Run(port); err != nil {
			t.Fatal(err)
		}
	}()
	// wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(200 * time.Millisecond)

	url := baseURL + port
	testResponseBody(t, url, http.StatusOK, "Main Page!")

	url = baseURL + port + "/goose"
	testResponseBody(t, url, http.StatusOK, "I love goose!")
}

func TestColonWildcardRouting(t *testing.T) {
	port := ":8081"
	g := New()
	go func() {
		g.GET("/:name", func(ctx *Context) {
			ctx.String("I love %s!", ctx.Param(":name"))
		})
		g.GET("/category/:category", func(ctx *Context) {
			ctx.String("Category: %s", ctx.Param(":category"))
		})
		if err := g.Run(port); err != nil {
			t.Fatal(err)
		}
	}()
	// wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(200 * time.Millisecond)

	url := baseURL + port + "/goose"
	testResponseBody(t, url, http.StatusOK, "I love goose!")

	url = baseURL + port + "/category/history"
	testResponseBody(t, url, http.StatusOK, "Category: history")
}

func TestRoutingWithAndWithoutTrailingSlash(t *testing.T) {
	port := ":8085"
	g := New()
	go func() {
		g.GET("/info", func(ctx *Context) {
			ctx.String("Information page")
		})
		g.GET("/info/", func(ctx *Context) {
			ctx.String("Trailing information page")
		})

		g.GET("/animal/:name", func(ctx *Context) {
			ctx.String("I love %s!", ctx.Param(":name"))
		})
		g.GET("/animal/:name/", func(ctx *Context) {
			ctx.String("I love %s 3000 times!", ctx.Param(":name"))
		})

		if err := g.Run(port); err != nil {
			t.Fatal(err)
		}
	}()
	// wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(200 * time.Millisecond)

	url := baseURL + port + "/info"
	testResponseBody(t, url, http.StatusOK, "Information page")

	url = baseURL + port + "/info/"
	testResponseBody(t, url, http.StatusOK, "Trailing information page")

	url = baseURL + port + "/animal/goose"
	testResponseBody(t, url, http.StatusOK, "I love goose!")

	url = baseURL + port + "/animal/goose/"
	testResponseBody(t, url, http.StatusOK, "I love goose 3000 times!")
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
