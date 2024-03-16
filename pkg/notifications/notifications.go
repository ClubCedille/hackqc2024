package notifications

import (
	"fmt"
	"log"

	"github.com/ClubCedille/hackqc2024/pkg/account"
	"github.com/ClubCedille/hackqc2024/pkg/geometry"
	mapobject "github.com/ClubCedille/hackqc2024/pkg/map_object"
	"github.com/ostafen/clover/v2"
)

func SendNotification(message string, recipients []string) {
	// TODO: Implement notification sending
	fmt.Println("Sending Notification, message: ", message, " recipients: ", recipients)
}

func NotifyNearby(db *clover.DB, message string, geom mapobject.Geometry) error {
	accounts, err := account.GetAllAccounts(db)
	if err != nil {
		log.Println("Error getting accounts:", err)
		return err
	}

	accountIds := []string{}
	for _, account := range accounts {
		if geometry.IsInGeom(account.Coordinates, geom) {
			accountIds = append(accountIds, account.Id)
		}
	}

	SendNotification(
		message,
		accountIds,
	)

	return nil
}
