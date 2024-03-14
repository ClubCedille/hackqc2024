package comment

import (
	"log"
	"time"

	"github.com/ClubCedille/hackqc2024/pkg/account"
	"github.com/ClubCedille/hackqc2024/pkg/database"
	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
	uuid "github.com/satori/go.uuid"
)

type Comment struct {
	Id          string `json:"_id" clover:"_id"`
	Comment     string `json:"comment" clover:"comment"`
	OwnerId     string `json:"owner_id" clover:"owner_id"`
	TargetId    string `json:"target_id" clover:"target_id"`
	CommentDate string `json:"comment_date" clover:"comment_date"`
}

type CommentFormData struct {
	Comment     string `json:"comment"`
	OwnerName   string `json:"owner_name"`
	TargetId    string `json:"target_id"`
	CommentDate string `json:"comment_date"`
}

func GetComments(conn *clover.DB, eventId string) ([]Comment, error) {
	docs, err := conn.FindAll(query.NewQuery(database.CommentCollection).Where(query.Field("target_id").Eq(eventId)))
	if err != nil {
		return nil, err
	}

	comments := []Comment{}
	for _, d := range docs {
		var comment Comment
		d.Unmarshal(&comment)
		comments = append(comments, comment)
	}

	return comments, nil
}

func CreateComment(conn *clover.DB, comment Comment) error {
	comment.Id = uuid.NewV4().String()
	comment.CommentDate = time.Now().Format(time.RFC3339)
	commentDoc := document.NewDocumentOf(comment)
	err := conn.Insert(database.CommentCollection, commentDoc)
	if err != nil {
		return err
	}

	return nil
}

func GetCommentsByTargetId(conn *clover.DB, targetId string) ([]Comment, error) {
	docs, err := conn.FindAll(query.NewQuery(database.CommentCollection).Where(query.Field("target_id").Eq(targetId)))
	if err != nil {
		return nil, err
	}

	comments := []Comment{}
	for _, d := range docs {
		var comment Comment
		d.Unmarshal(&comment)
		comments = append(comments, comment)
	}

	return comments, nil
}

func GetCommentsFormData(db *clover.DB, targetId string) ([]CommentFormData, error) {
	comments, err := GetCommentsByTargetId(db, targetId)
	if err != nil {
		log.Println("Error fetching comments:", err)
		return nil, err
	}

	commentsFormData := []CommentFormData{}
	for _, comment := range comments {
		owner, err := account.GetAccountById(db, comment.OwnerId)
		if err != nil {
			log.Println("Error getting owner account:", err)
			return nil, err
		}
		commentsFormData = append(commentsFormData, CommentFormData{
			Comment:     comment.Comment,
			OwnerName:   owner.FirstName + " " + owner.LastName,
			TargetId:    comment.TargetId,
			CommentDate: comment.CommentDate,
		})
	}

	return commentsFormData, nil
}
