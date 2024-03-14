package main

import (
	"log"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/session"
	"github.com/gin-gonic/gin"
)

// Check if user is logged in to enable a group of features
func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session.GetActiveSession(c)
		log.Printf("(middleware) accountid value: %s", session.ActiveSession.AccountId)
		if session.ActiveSession.AccountId == "" {
			redirectURL := c.Request.URL.Path
			c.SetCookie("redirect_url", redirectURL, 3600, "/", "", false, true)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}
	}
}

func LoginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("(middleware) getting active session")
		session.GetActiveSession(c)
	}
}
