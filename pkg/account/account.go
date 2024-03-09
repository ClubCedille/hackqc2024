package account

import (
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Id        string `json:"_id" clover:"_id"`
	UserName  string `json:"user_name" clover:"user_name"`
	FirstName string `json:"first_name" clover:"first_name"`
	LastName  string `json:"last_name" clover:"last_name"`
	Email     string `json:"email" clover:"email"`
}

func CreateAccount(conn *clover.DB, account Account) error {
	account.Id = uuid.NewV4().String()
	accountDoc := document.NewDocumentOf(account)
	err := conn.Insert(database.AccountCollection, accountDoc)
	if err != nil {
		return err
	}

	return nil
}

func UpdateAccount(conn *clover.DB, account Account) error {
	err := conn.UpdateById(database.EventCollection, account.Id, func(doc *document.Document) *document.Document {
		doc.Set("user_name", account.UserName)
		doc.Set("first_name", account.FirstName)
		doc.Set("last_name", account.LastName)
		doc.Set("email", account.Email)
		return doc
	})
	if err != nil {
		return err
	}

	return nil
}
