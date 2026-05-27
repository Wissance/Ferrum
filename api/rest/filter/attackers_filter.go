package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/wissance/Ferrum/logging"
	"github.com/wissance/Ferrum/security/bruteforce"
	sf "github.com/wissance/stringFormatter"
	"net/http"
)

func AttackersFilterMiddleware(service *bruteforce.ProtectionService, logger *logging.AppLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check attackers, if it was blocked, return 418 (fun) or 429 (too many requests)
		// TODO(UMV): add check deviceId when we get first method with browser and redirection (auth)
		ipAddress := c.ClientIP()
		isBlocked := (*service).IsIpAddressBlocked(ipAddress)
		if !isBlocked {
			c.Next()
		} else {
			msg := sf.Format("Registered attempt to access from blocked IP Address: {0}", ipAddress)
			logger.Debug(msg)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": msg,
			})
		}
	}
}
