package account

import (
	"regexp"
	"slices"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	Id          string    `json:"_id" clover:"_id"`
	UserName    string    `json:"user_name" clover:"user_name"`
	FirstName   string    `json:"first_name" clover:"first_name"`
	LastName    string    `json:"last_name" clover:"last_name"`
	Email       string    `json:"email" clover:"email"`
	PhoneNumber string    `json:"phone_number" clover:"phone_number"`
	Coordinates []float64 `json:"coordinates" clover:"coordinates"`
}

func AccountExistsById(conn *clover.DB, id string) (bool, error) {
	return conn.Exists(query.NewQuery(database.AccountCollection).Where(query.Field("_id").Eq(id)))
}

func AccountExistByEmailAndUsername(conn *clover.DB, email string, userName string) (bool, error) {
	return conn.Exists(query.NewQuery(database.AccountCollection).Where(query.Field("user_name").Eq(userName).And(query.Field("email").Eq(email))))
}

func CreateAccount(conn *clover.DB, account Account) error {
	account.Id = uuid.NewV4().String()
	accountDoc := document.NewDocumentOf(account)
	r := regexp.MustCompile("[^0-9.]")
	account.PhoneNumber = r.ReplaceAllString(account.PhoneNumber, "")
	err := conn.Insert(database.AccountCollection, accountDoc)
	if err != nil {
		return err
	}

	return nil
}

func GetAccountByUsername(conn *clover.DB, username string) (Account, error) {

	docs, err := conn.FindFirst(query.NewQuery(database.AccountCollection).Where(query.Field("user_name").Eq(username)))
	if err != nil {
		return Account{}, err
	}

	account := Account{}
	docs.Unmarshal(&account)

	return account, nil
}

func GetAccountById(conn *clover.DB, id string) (Account, error) {
	docs, err := conn.FindFirst(query.NewQuery(database.AccountCollection).Where(query.Field("_id").Eq(id)))
	if err != nil {
		return Account{}, err
	}

	account := Account{}
	docs.Unmarshal(&account)

	return account, nil
}

func GetAccountsByIds(conn *clover.DB, ids []string) ([]Account, error) {

	filterQuery := query.NewQuery(database.AccountCollection)
	// filterQuery = filterQuery.MatchFunc(func(doc *document.Document) bool {
	// 	return slices.Contains(ids, doc.Get("_id").(string))
	// })

	docs, err := conn.FindAll(filterQuery)
	if err != nil {
		return []Account{}, err
	}

	accounts := []Account{}
	for _, d := range docs {
		if slices.Contains(ids, d.Get("_id").(string)) {
			var account Account
			d.Unmarshal(&account)
			accounts = append(accounts, account)
		}
	}

	return accounts, nil
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

func GetAllAccountsWithCoords(conn *clover.DB) ([]Account, error) {
	docs, err := conn.FindAll(query.NewQuery(database.AccountCollection).Where(query.Field("coordinates").IsNilOrNotExists().Not()))
	if err != nil {
		return []Account{}, err
	}

	accounts := []Account{}
	for _, d := range docs {
		var account Account
		err = d.Unmarshal(&account)
		if err != nil {
			return []Account{}, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}
