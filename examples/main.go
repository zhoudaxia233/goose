package main

import (
	"fmt"

	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()

	g.GET("/", func(ctx *goose.Context) {
		fmt.Fprintf(ctx.ResponseWriter, "Hello World!")
	})

	g.GET("/info", func(ctx *goose.Context) {
		fmt.Fprintf(ctx.ResponseWriter, "My name is goose!")
	})

	g.Run(":8080")
}
