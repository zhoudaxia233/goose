package goose

import (
	"log"
	"time"
)

// Time is a middleware used for logging elapsed time in each http request (per goose.Context)
func Time(ctx *Context) {
	startTime := time.Now()
	ctx.Next()
	timeElapsed := time.Since(startTime)
	log.Printf("%v has passed...", timeElapsed)
}
