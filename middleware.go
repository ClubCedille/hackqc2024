package main

import (
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/session"
	"github.com/gin-gonic/gin"
)

// Check if user is logged in
func AuthRequiredMiddleware(c *gin.Context) {
	if session.ActiveSession.AccountId == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}
}
