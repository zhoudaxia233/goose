package main

import (
	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()

	g.GET("/", func(ctx *goose.Context) {
		ctx.String("Hello World!")
	})

	g.GET("/info", func(ctx *goose.Context) {
		ctx.HTML("<h1>My name is goose!</h1>")
	})

	g.Run(":8080")
}
