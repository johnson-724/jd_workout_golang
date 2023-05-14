package middleware

import (
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func FailResponseAlert() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		status := c.Writer.Status()
		if status < 200 || status >= 300 {
			err := c.Err()
			if err != nil {
				sentry.CaptureException(err)
			} else {
				sentry.CaptureMessage("status code is not 2XX, route: " + c.Request.URL.Path)
			}
		}
	}
}
