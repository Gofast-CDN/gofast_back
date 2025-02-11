package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gofast/database"

	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		ctx := context.Background()
		limit := 10
		window := 60

		count, err := database.RedisClient.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur Redis"})
			c.Abort()
			return
		}

		// expiration de la clé si c'est la première requête
		if count == 1 {
			database.RedisClient.Expire(ctx, key, time.Duration(window)*time.Second)
		}

		// vérification du nombre de requêtes
		if count > int64(limit) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Trop de requêtes. Veuillez réessayer plus tard.",
			})
			c.Abort()
			return
		}

		remaining := limit - int(count)
		c.Writer.Header().Set("X-RateLimit-Limit", strconv.Itoa(limit))
		c.Writer.Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))

		c.Next()
	}
}
