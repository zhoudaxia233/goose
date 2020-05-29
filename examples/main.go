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

	// g.DrawRoutingTree("GET")

	g.Run(":8080")
}
