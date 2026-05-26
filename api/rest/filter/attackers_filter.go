package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/wissance/Ferrum/security/bruteforce"
)

func AttackersFilterMiddleware(service bruteforce.ProtectionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// check attackers, if it was blocked, return
		c.Next()
	}
}
