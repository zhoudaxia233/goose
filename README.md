# goose

<a href="https://github.com/zhoudaxia233/goose"><img height="190px" src="logo.svg"></a>

**goose** is a lightweight web framework in Go.

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
		ctx.String("My name is %s", ctx.Param(":name"))
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

	v1 := g.Group("/v1")
	{
		v1.GET("/", func(ctx *goose.Context) {
			ctx.HTML("<h1>V1 PAGE!</h1>")
		})

		v1.GET("/hello", func(ctx *goose.Context) {
			ctx.String("Hello V1!")
		})

		// goose also supports nested router group
		v2 := v1.Group("/v2")
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