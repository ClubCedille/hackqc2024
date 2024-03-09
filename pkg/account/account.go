package account

import (
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

type Account struct {
	Id        string `clover:"_id"`
	UserName  string `clover:"user_name"`
	FirstName string `clover:"first_name"`
	LastName  string `clover:"last_name"`
	Email     string `clover:"email"`
}

func AccountExistsById(conn *clover.DB, id string) (bool, error) {
	return conn.Exists(query.NewQuery(database.AccountCollection).Where(query.Field("_id").Eq(id)))
}

func CreateAccount(conn *clover.DB, account Account) error {
	accountDoc := document.NewDocumentOf(account)
	err := conn.Insert(database.AccountCollection, accountDoc)
	if err != nil {
		return err
	}

	return nil
}
