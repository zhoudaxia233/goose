package goose

import (
	"log"
	"time"
)

// Time is a middleware used for logging elapsed time since the server starts
func Time(ctx *Context) {
	startTime := time.Now()
	ctx.Next()
	timeElapsed := time.Since(startTime)
	log.Printf("%v has passed...", timeElapsed)
}
