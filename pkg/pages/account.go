package pages

import (
	"log"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/account"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func CreateAccount(c *gin.Context, db *clover.DB) {
	var data account.Account
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := account.CreateAccount(db, data)
	if err != nil {
		log.Println("Error creating account:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Account created successfully")
	c.Redirect(http.StatusSeeOther, "/")
}

func UpdateAccount(c *gin.Context, db *clover.DB) {
	var data account.Account
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := account.UpdateAccount(db, data)
	if err != nil {
		log.Println("Error updating account:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("Account updated successfully")
	c.Redirect(http.StatusSeeOther, "/")
}
