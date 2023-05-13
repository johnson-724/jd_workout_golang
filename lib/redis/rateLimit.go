package redis

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/gin-gonic/gin"
	"context"
	"fmt"
	"os"
	"strconv"
)

func ApiRateLimit(c *gin.Context) error {
	limiter := redis_rate.NewLimiter(Connection)

	ip := c.ClientIP()
	ctx := context.Background()
	rate, err := strconv.Atoi(os.Getenv("API_RATE_PER_MINUTE"))

	if err != nil {
		return err
	}

	res, err := limiter.Allow(ctx, fmt.Sprintf("api_rate_%s", ip), redis_rate.PerMinute(rate))

	if err != nil {
		return err
	}

	if res.Allowed == 0 {
		return fmt.Errorf("too many requests")
	}

	return nil
}
