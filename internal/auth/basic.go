package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func BasicAuth(c *gin.Context) {
// 	u, p, ok := c.Request.BasicAuth()
// 	log.Println("---basic auth---")
// 	log.Println(u, p, ok)
// 	if !ok || u != "admin" || p != "1234" {
// 		c.Writer.Header().Set("WWW-Authenticate", "Basic")
// 		c.AbortWithStatus(http.StatusUnauthorized)
// 		return
// 	}
// }

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// func BasicAuth() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		u, p, ok := c.Request.BasicAuth()
// 		log.Println("---basic auth---")
// 		log.Println(u, p, ok)
// 		if !ok {
// 			c.Writer.Header().Set("WWW-Authenticate", "Basic")
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 			return
// 		}
// 		wrongCredentials := u != "admin" || p != "1234"
// 		if wrongCredentials {
// 			c.Writer.Header().Set("WWW-Authenticate", "Basic") /// in case of open in web browser it appear login popup
// 			c.AbortWithStatus(http.StatusUnauthorized)
// 			return
// 		}
// 	}
// }

func BasicAuth(credentials []Credential) gin.HandlerFunc {
	return func(c *gin.Context) {
		username, password, ok := c.Request.BasicAuth()
		log.Println("---basic auth---")
		log.Println(username, password, ok)
		if !ok {
			c.Writer.Header().Set("WWW-Authenticate", "Basic")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		for _, v := range credentials {
			if v.Username == username && v.Password == password {
				c.Next()
				return
			}
		}
		c.Writer.Header().Set("WWW-Authenticate", "Basic")
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
