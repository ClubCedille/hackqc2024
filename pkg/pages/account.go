package pages

import (
	"log"
	"net/http"

	"github.com/ClubCedille/hackqc2024/pkg/account"
	"github.com/ClubCedille/hackqc2024/pkg/geometry"
	"github.com/ClubCedille/hackqc2024/pkg/session"
	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func GetCreateAccount(c *gin.Context) {
	c.HTML(http.StatusOK, "forms/accountForm.html", gin.H{
		"SigningUp": "true",
	})
}

func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "forms/loginForm.html", gin.H{
		"LoggingIn": "true",
	})
}

func CreateAccount(c *gin.Context, db *clover.DB) {
	var data account.Account
	var err error
	data.UserName = c.PostForm("user_name")
	data.FirstName = c.PostForm("first_name")
	data.LastName = c.PostForm("last_name")
	data.Email = c.PostForm("email")
	data.PhoneNumber = c.PostForm("phone_number")
	data.Coordinates, err = geometry.ParseCoordinatesString(c.PostForm("coordinates"))
	if err != nil {
		data.Coordinates = []float64{}
		err = nil
	}

	exist, err := account.AccountExistByEmailAndUsername(db, c.PostForm("email"), c.PostForm("user_name"))
	if err != nil {
		log.Println("Error looking up account:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if !exist {
		err := account.CreateAccount(db, data)
		if err != nil {
			log.Println("Error creating account:", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		setActiveSession(c, db, data.UserName)

	} else {
		log.Println("Account already exists")
		c.HTML(http.StatusOK, "forms/accountForm.html", gin.H{
			"Error": "Le nom d'utilisateur existe déjà ou l'adresse courriel est déjà utilisée.",
		})
		return
	}

	log.Println("Account created successfully")
	c.Redirect(http.StatusSeeOther, "/map")
}

func Login(c *gin.Context, db *clover.DB) {
	exist, err := account.AccountExistByEmailAndUsername(db, c.PostForm("email"), c.PostForm("user_name"))
	if err != nil {
		log.Println("Error looking up account:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if !exist {
		c.HTML(http.StatusUnauthorized, "forms/loginForm.html", gin.H{
			"Error": "Le compte n'existe pas.",
		})
		return
	} else {
		setActiveSession(c, db, c.PostForm("user_name"))
		redirectURL, err := c.Cookie("redirect_url")
		if err != nil || redirectURL == "" {
			redirectURL = "/map"
		}

		c.SetCookie("redirect_url", "", -1, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, redirectURL)
	}
}

func Logout(c *gin.Context) {
	session.ClearActiveSession(c)
	c.Redirect(http.StatusSeeOther, "/map")
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
	c.Redirect(http.StatusSeeOther, "/map")
}

func setActiveSession(c *gin.Context, db *clover.DB, userName string) {
	acc, err := account.GetAccountByUsername(db, userName)
	if err != nil {
		log.Println("Error getting account:", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	session.SetActiveSession(c, session.SessionInformation{
		AccountId: acc.Id,
		Email:     acc.Email,
		UserName:  acc.UserName,
	})
}
