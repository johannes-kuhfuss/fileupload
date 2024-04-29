package app

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/johannes-kuhfuss/services_utils/logger"
)

func basicAuth(users map[string]string) gin.HandlerFunc {
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
						logger.Info(fmt.Sprintf("User %v not found", givenUser))
						return false, ""
					} else {
						/*
							err := bcrypt.CompareHashAndPassword([]byte(givenPass), []byte(pwhash))
							if err == nil {
								return true, givenUser
							} else {
								h, _ := bcrypt.GenerateFromPassword([]byte(givenPass), bcrypt.DefaultCost)
								logger.Info(fmt.Sprintf("Password %v does not equal hash %v. Hash: %v", givenPass, pwhash, string(h)))
							}
						*/
						if givenPass == pwhash {
							return true, givenUser
						}
					}
				}
			}
		}
	}
	return false, ""
}
