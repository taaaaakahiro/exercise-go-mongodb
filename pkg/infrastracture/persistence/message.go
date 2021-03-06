package persistence

import (
	"context"
	"fmt"

	"exercise-go-api/pkg/domain/entity"
	"exercise-go-api/pkg/domain/repository"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	messageCollection = "message"
)

type MessageRepo struct {
	col *mongo.Collection
}

var _ repository.IMessageRepository = (*MessageRepo)(nil)

func NewMessageRepository(db *mongo.Database) *MessageRepo {
	return &MessageRepo{
		col: db.Collection(messageCollection),
	}
}

func (r MessageRepo) ListMessages(ctx context.Context, userId int) ([]entity.Message, error) {
	messages := make([]entity.Message, 0)
	srt := bson.D{
		primitive.E{Key: "id", Value: -1},
	}
	opt := options.Find().SetSort(srt)
	flt := bson.D{
		primitive.E{Key: "userId", Value: userId},
	}

	cur, err := r.col.Find(ctx, flt, opt)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var message entity.Message
		err := cur.Decode(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	// data := `[
	// 	{"id": 1, "userId": 1, "message": "test message id 1"},
	// 	{"id": 2, "userId": 1, "message": "test message id 2"},
	// 	{"id": 3, "userId": 2, "message": "test message id 3"},
	// 	{"id": 3, "userId": 2, "message": "test message id 3"},
	// 	{"id": 4, "userId": 2, "message": "test message id 4"}
	// ]`
	// var str []entity.Message
	// err = json.Unmarshal([]byte(data), &str)
	// for _, res := range str {
	// 	fmt.Println(res)
	// 	fmt.Println(res.Id)
	// 	fmt.Println(res.Message)
	// }
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = cur.Close(ctx)
	if err != nil {
		return nil, err
	}

	return messages, nil
	// return str, nil
}
