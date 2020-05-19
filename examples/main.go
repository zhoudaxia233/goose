package main

import (
	"fmt"
	"net/http"

	"github.com/zhoudaxia233/goose"
)

func main() {
	g := goose.New()

	g.GET("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	g.GET("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "My name is goose!")
	})

	g.Run(":8080")
}
