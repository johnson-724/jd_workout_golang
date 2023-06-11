package middleware

import (
	"fmt"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
	"time"
)

func FailResponseAlert() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v\n", r)
				log.Printf("Stack trace:\n%s", getStackTrace())
				sentry.CurrentHub().Recover(r)
				sentry.Flush(time.Second * 5)

				c.AbortWithStatus(500)
			}
		}()

		c.Next()

		status := c.Writer.Status()
		if status < 200 || status >= 300 {
			err := c.Err()
			if err != nil {
				sentry.CaptureException(err)
			} else {
				sentry.CaptureMessage(fmt.Sprintf("status code is %d, route: %s", c.Writer.Status(), c.Request.URL.Path))
			}
		}
	}
}

func getStackTrace() string {
	buf := make([]byte, 1<<16)
	stackSize := runtime.Stack(buf, true)
	return string(buf[:stackSize])
}
