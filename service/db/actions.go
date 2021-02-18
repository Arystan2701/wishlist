package db

import (
	"context"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/telegram/wishlist/common"
	"go.mongodb.org/mongo-driver/mongo"
)

type actionRepo struct {
	Col   *mongo.Collection
	Redis *redis.Client
}

func NewActionRepo(client *mongo.Client, redis *redis.Client) *actionRepo {
	col := client.Database(databaseWishlist).Collection(collectionActions)
	return &actionRepo{
		Col:   col,
		Redis: redis,
	}
}

type Action struct {
	Sender      ShortUser  `json:"sender" bson:"sender"`
	Type        string     `json:"type" bson:"type"`
	Recipient   *ShortUser `json:"recipient" bson:"recipient`
	SearchQuery string     `json:"search_query" bson:"search_query"`
}

func (a *actionRepo) InsertAction(actionType string, message *tgbotapi.Message, query string, recipient *User) error {
	action := &Action{
		Sender: ShortUser{
			TelegramUserID: message.From.ID,
			UserName:       message.From.UserName,
		},
		Type: actionType,
	}
	if actionType == common.ActionTypeSearchUser {
		action.SearchQuery = query
		if recipient != nil {
			action.Recipient = &ShortUser{
				TelegramUserID: recipient.TelegramUserID,
				UserName:       recipient.UserName,
			}
		}
	}
	_, err := a.Col.InsertOne(context.TODO(), action)
	return err
}
