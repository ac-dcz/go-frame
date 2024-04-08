package geeweb

import (
	"log"
	"time"
)

func LoggerMiddleWare() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%s] %s in %v", c.Method, c.R.RequestURI, time.Since(t))
	}
}
