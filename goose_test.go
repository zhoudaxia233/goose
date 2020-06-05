package goose

import (
	"net/http"
	"strconv"
	"strings"
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
			ctx.String("I love %s!", ctx.Param("name"))
		})
		g.GET("/category/:category", func(ctx *Context) {
			ctx.String("Category: %s", ctx.Param("category"))
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
			ctx.String("I love %s!", ctx.Param("name"))
		})
		g.GET("/animal/:name/", func(ctx *Context) {
			ctx.String("I love %s 3000 times!", ctx.Param("name"))
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

func TestRouterGroup(t *testing.T) {
	port := ":8086"
	g := New()
	go func() {
		g.GET("/", func(ctx *Context) {
			ctx.String("Root page")
		})

		v1 := g.Group("v1")
		{
			v1.GET("/", func(ctx *Context) {
				ctx.String("Group V1 is here!")
			})

			v1.GET("/hello", func(ctx *Context) {
				ctx.String("Hello Group V1!")
			})

			v2 := v1.Group("v2")
			{
				v2.GET("/hello", func(ctx *Context) {
					ctx.String("Hello Group V2!")
				})
			}
		}

		if err := g.Run(port); err != nil {
			t.Fatal(err)
		}
	}()
	// wait for the goroutine to start and run the server
	// otherwise the main thread will complete
	time.Sleep(200 * time.Millisecond)

	url := baseURL + port + "/"
	testResponseBody(t, url, http.StatusOK, "Root page")

	url = baseURL + port + "/v1/"
	testResponseBody(t, url, http.StatusOK, "Group V1 is here!")

	url = baseURL + port + "/v1/hello"
	testResponseBody(t, url, http.StatusOK, "Hello Group V1!")

	url = baseURL + port + "/v1/v2/hello"
	testResponseBody(t, url, http.StatusOK, "Hello Group V2!")

}

func TestTemplate(t *testing.T) {
	g := New()
	g.FuncMap(X{
		"appendYear": func(s string) string {
			year := time.Now().Year()
			return strings.Join([]string{s, strconv.Itoa(year)}, " - ")
		},
	})
	g.Set("toUpper", strings.ToUpper)
	g.LoadHTMLGlob("testfiles/templates/*")

	g.GET("/t", func(ctx *Context) {
		ctx.HTML("hello.tmpl", X{"name": "Goose"})
	})

	g.GET("/func", func(ctx *Context) {
		ctx.HTML("funcmaps.tmpl", X{"msg": "I love goose!"})
	})

	want := "<h1>Hello Goose</h1>"
	testResponseBodyLocal(t, g, "/t", 200, want)

	want = "I LOVE GOOSE! - 2020"
	testResponseBodyLocal(t, g, "/func", 200, want)
}
