package session

import (
	"log"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionInformation struct {
	AccountId string
	Email     string
	UserName  string
}

var ActiveSession SessionInformation

func SetActiveSession(c *gin.Context, sessionInformation SessionInformation) {
	session := sessions.Default(c)
	session.Set("account_id", sessionInformation.AccountId)
	session.Set("user_name", sessionInformation.UserName)
	session.Set("email", sessionInformation.Email)
	session.Save()

	ActiveSession = sessionInformation

	log.Println("Session set:", sessionInformation)
}

func GetActiveSession(c *gin.Context) {
	session := sessions.Default(c)
	accountId := session.Get("account_id")
	userName := session.Get("user_name")
	email := session.Get("email")

	if accountId != nil && userName != nil && email != nil {
		ActiveSession.AccountId = accountId.(string)
		ActiveSession.UserName = userName.(string)
		ActiveSession.Email = email.(string)
	}
}

func ClearActiveSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()

	ActiveSession = SessionInformation{}
}
