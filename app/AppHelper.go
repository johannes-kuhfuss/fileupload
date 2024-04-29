package app

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/fileupload/domain"
	"golang.org/x/crypto/bcrypt"
)

func basicAuth(users domain.Users) gin.HandlerFunc {
	realm := "Basic realm=Authorization Required"
	return func(c *gin.Context) {
		authHeaders := c.Request.Header["Authorization"]
		found, user := checkAuth(users, authHeaders)
		if !found {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(gin.AuthUserKey, user)
	}
}

func checkAuth(users domain.Users, authHeaders []string) (found bool, user string) {
	for _, header := range authHeaders {
		if strings.HasPrefix(header, "Basic ") {
			b64Creds := header[6:]
			creds, err := base64.StdEncoding.DecodeString(b64Creds)
			if err == nil {
				parts := strings.Split(string(creds), ":")
				if len(parts) == 2 {
					givenUser := parts[0]
					givenPass := parts[1]
					exists, pwhash := users.CheckUser(givenUser)
					if !exists {
						return false, ""
					} else {
						err := bcrypt.CompareHashAndPassword([]byte(givenPass), []byte(pwhash))
						if err == nil {
							return true, givenUser
						}
					}
				}
			}
		}
	}
	return false, ""
}
