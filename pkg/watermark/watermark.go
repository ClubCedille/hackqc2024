package watermark

import (
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

type Watermark struct {
	Id        string    `json:"_id" clover:"_id"`
	Name      string    `json:"name" clover:"name"`
	Watermark time.Time `json:"watermark" clover:"watermark"`
}

func WatermarkExistsByName(db *clover.DB, name string) (bool, error) {
	return db.Exists(query.NewQuery(database.WatermarkCollection).Where(query.Field("name").Eq(name)))
}

func GetWatermark(db *clover.DB, name string) (Watermark, error) {
	docs, err := db.FindFirst(query.NewQuery(database.WatermarkCollection).Where(query.Field("name").Eq(name)))
	if err != nil {
		return Watermark{}, err
	}

	watermark := Watermark{}
	docs.Unmarshal(&watermark)

	return watermark, nil
}

func CreateWatermark(db *clover.DB, watermark Watermark) error {
	doc := document.NewDocumentOf(watermark)
	err := db.Insert(database.WatermarkCollection, doc)
	if err != nil {
		return err
	}

	return nil
}

func UpdateWatermark(db *clover.DB, watermark Watermark) error {
	err := db.UpdateById(database.WatermarkCollection, watermark.Id, func(doc *document.Document) *document.Document {
		doc.Set("watermark", watermark.Watermark)
		return doc
	})
	if err != nil {
		return err
	}

	return nil
}
