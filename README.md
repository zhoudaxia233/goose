# goose

<a href="https://github.com/zhoudaxia233/goose"><img height="190px" src="logo.svg"></a>

**goose** is a lightweight web framework in Go.

>Note: Currently, goose is still not ready to release. You should not use it for your project since the APIs may change a lot. Also, there are still many features for me to implement...

<details>
<summary><strong>A hello world example</strong></summary>

```go
package main

import (
	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()

	g.GET("/", func(ctx *goose.Context) {
		ctx.String("Hello World!")
	})

	g.Run(":8080")
}
```

</details>

## Contents
- [goose](#goose)
	- [Contents](#contents)
	- [Features](#features)
		- [Dynamic Routing](#dynamic-routing)
		- [Router Group](#router-group)
		- [Middleware](#middleware)
		- [Static Files](#static-files)
		- [Templates](#templates)
	- [Acknowledgment](#acknowledgment)

## Features
### Dynamic Routing

<details>
<summary><strong>An example</strong></summary>

```go
package main

import (
	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()

	g.GET("/info/:name", func(ctx *goose.Context) {
		ctx.String("My name is %s", ctx.Param("name"))
	})

	g.Run(":8080")
}

```

</details>

### Router Group

<details>
<summary><strong>An example</strong></summary>

```go
package main

import (
	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()

	v1 := g.Group("v1")
	{
		v1.GET("/", func(ctx *goose.Context) {
			ctx.String("Page V1!")
		})

		v1.GET("/hello", func(ctx *goose.Context) {
			ctx.String("Hello V1!")
		})

		// goose also supports nested router group
		v2 := v1.Group("v2")
		{
			v2.GET("/hello", func(ctx *goose.Context) {
				ctx.String("Hello V2!")
			})
		}
	}

	g.Run(":8080")
}

```

</details>

### Middleware

<details>
<summary><strong>An example</strong></summary>

```go
package main

import (
	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()
	g.Use(func(ctx *goose.Context) {
		log.Println("here get executed before handling the request")
		ctx.Next()
		log.Println("here get executed after handling the request")
	})

	g.GET("/", func(ctx *goose.Context) {
		ctx.String("Hello World!")
	})

	v1 := g.Group("v1")
	v1.Use(func(ctx *goose.Context) {
		log.Println("before v1")
		ctx.Next()
		log.Println("after v1")
	})

	v1.GET("/hello", func(ctx *goose.Context) {
		ctx.String("Hello V1!")
	})

	g.Run(":8080")
}

```

</details>

### Static Files

<details>
<summary><strong>An example</strong></summary>

```go
package main

import (
	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()

	g.Static("/assets", "examples/static")
	g.StaticFile("/favicon.ico", "examples/favicon.ico")

	g.Run(":8080")
}

```

</details>

### Templates

<details>
<summary><strong>An example</strong></summary>

```go
package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()
	g.FuncMap(goose.X{
		"appendYear": func(s string) string {
			year := time.Now().Year()
			return strings.Join([]string{s, strconv.Itoa(year)}, " - ")
		},
	})
	g.Set("toUpper", strings.ToUpper)
	g.LoadHTMLGlob("testfiles/templates/*")

	g.GET("/", func(ctx *goose.Context) {
		ctx.HTML("hello.tmpl", goose.X{"name": "Goose"})
	})

	g.GET("/func", func(ctx *goose.Context) {
		ctx.HTML("funcmaps.tmpl", goose.X{"msg": "I love goose!"})
	})

	g.Run(":8080")
}

```

</details>

## Acknowledgment

1. I've got some useful design inspiration from [Gin](https://github.com/gin-gonic/gin).
2. The goose icon in the logo is made by [monkik](https://www.flaticon.com/authors/monkik) from [Flaticon](https://www.flaticon.com/).
