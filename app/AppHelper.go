package app

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/services_utils/api_error"
	"github.com/johannes-kuhfuss/services_utils/logger"
	"golang.org/x/crypto/bcrypt"
)

func basicAuth(users map[string]string) gin.HandlerFunc {
	realm := "Basic realm=Authorization Required"
	return func(c *gin.Context) {
		authHeaders := c.Request.Header["Authorization"]
		found, user := checkAuth(users, authHeaders)
		_ = user
		if !found {
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		} else {
			session := sessions.Default(c)
			session.Set("uploadUser", user)
			err := session.Save()
			if err != nil {
				msg := "could not save session"
				logger.Error(msg, err)
				apiErr := api_error.NewInternalServerError(msg, err)
				c.JSON(apiErr.StatusCode(), apiErr)
			}
		}
	}
}

func checkAuth(users map[string]string, authHeaders []string) (found bool, user string) {
	for _, header := range authHeaders {
		if strings.HasPrefix(header, "Basic ") {
			b64Creds := header[6:]
			creds, err := base64.StdEncoding.DecodeString(b64Creds)
			if err == nil {
				parts := strings.Split(string(creds), ":")
				if len(parts) == 2 {
					givenUser := parts[0]
					givenPass := parts[1]
					pwhash, exists := users[givenUser]
					if !exists {
						logger.Warn(fmt.Sprintf("User %v not found", givenUser))
						return false, ""
					} else {
						testHash, _ := base64.StdEncoding.DecodeString(pwhash)
						err := bcrypt.CompareHashAndPassword(testHash, []byte(givenPass))
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
