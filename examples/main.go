package main

import (
	"log"

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

	v1 := g.Group("/v1")
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
