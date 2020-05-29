package main

import (
	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()

	g.GET("/", func(ctx *goose.Context) {
		ctx.HTML("<h1>My name is goose!</h1>")
	})

	g.GET("/:name", func(ctx *goose.Context) {
		ctx.String("I love %s!", ctx.Param(":name"))
	})

	g.GET("/category/:category", func(ctx *goose.Context) {
		ctx.String("Category: %s", ctx.Param(":category"))
	})

	v1 := g.Group("/v1")
	{
		v1.GET("/", func(ctx *goose.Context) {
			ctx.HTML("<h1>V1 PAGE!</h1>")
		})

		v1.GET("/hello", func(ctx *goose.Context) {
			ctx.String("Hello V1!")
		})

		v2 := v1.Group("/v2")
		{
			v2.GET("/hello", func(ctx *goose.Context) {
				ctx.String("Hello V2!")
			})
		}
	}

	// g.DrawRoutingTree("GET")

	g.Run(":8080")
}
