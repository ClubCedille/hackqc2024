package account

import (
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
)

type Account struct {
	UserName  string `clover:"user_name"`
	FirstName string `clover:"first_name"`
	LastName  string `clover:"last_name"`
	Email     string `clover:"email"`
}

func CreateAccount(conn *clover.DB, account Account) error {
	accountDoc := document.NewDocumentOf(account)
	err := conn.Insert(database.AccountCollection, accountDoc)
	if err != nil {
		return err
	}

	return nil
}
