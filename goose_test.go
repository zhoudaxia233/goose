package goose

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	baseURL := "http://127.0.0.1"
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
	g := New()
	g.GET("/:name", func(ctx *Context) {
		ctx.String("I love %s!", ctx.Param("name"))
	})
	g.GET("/category/:category", func(ctx *Context) {
		ctx.String("Category: %s", ctx.Param("category"))
	})

	want := "I love goose!"
	testResponseBodyLocal(t, g, "/goose", http.StatusOK, want)

	want = "Category: history"
	testResponseBodyLocal(t, g, "/category/history", http.StatusOK, want)
}

func TestRoutingWithAndWithoutTrailingSlash(t *testing.T) {
	g := New()

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

	want := "Information page"
	testResponseBodyLocal(t, g, "/info", http.StatusOK, want)

	want = "Trailing information page"
	testResponseBodyLocal(t, g, "/info/", http.StatusOK, want)

	want = "I love goose!"
	testResponseBodyLocal(t, g, "/animal/goose", http.StatusOK, want)

	want = "I love goose 3000 times!"
	testResponseBodyLocal(t, g, "/animal/goose/", http.StatusOK, want)
}

func TestRouterGroup(t *testing.T) {
	g := New()
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
	want := "Root page"
	testResponseBodyLocal(t, g, "/", http.StatusOK, want)

	want = "Group V1 is here!"
	testResponseBodyLocal(t, g, "/v1/", http.StatusOK, want)

	want = "Hello Group V1!"
	testResponseBodyLocal(t, g, "/v1/hello", http.StatusOK, want)

	want = "Hello Group V2!"
	testResponseBodyLocal(t, g, "/v1/v2/hello", http.StatusOK, want)
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
	testResponseBodyLocal(t, g, "/t", http.StatusOK, want)

	want = "I LOVE GOOSE! - 2020"
	testResponseBodyLocal(t, g, "/func", http.StatusOK, want)
}
