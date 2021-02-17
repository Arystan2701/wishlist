package db

import (
	"context"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	UserCol *mongo.Collection
	Redis   *redis.Client
}

func NewUserRepository(client *mongo.Client, redisClient *redis.Client) *userRepository {
	userCol := client.Database(databaseWishlist).Collection(collectionUsers)
	return &userRepository{
		UserCol: userCol,
		Redis:   redisClient,
	}
}

type ShortUser struct {
	TelegramUserID int    `json:"telegram_user_id" bson:"telegram_user_id"`
	FirstName      string `json:"first_name" bson:"first_name"`
	LastName       string `json:"last_name" bson:"last_name"`
	UserName       string `json:"username" bson:"username"`
	Phone          string `json:"phone" bson:"phone"`
}

func (u *userRepository) CreateUser(message *tgbotapi.Message) (*ShortUser, error) {
	user := &ShortUser{
		TelegramUserID: message.From.ID,
		FirstName:      message.From.FirstName,
		LastName:       message.From.LastName,
		UserName:       message.From.UserName,
	}
	_, err := u.UserCol.InsertOne(context.TODO(), user)
	return user, err
}

func (u *userRepository) GetUser(userID int) (*ShortUser, error) {
	res := u.UserCol.FindOne(context.TODO(), bson.M{"telegram_user_id": userID})
	if err := res.Err(); err != nil {
		return nil, err
	}
	var user ShortUser
	err := res.Decode(&user)
	return &user, err

}

func (u *userRepository) UpdateUserPhone(contact *tgbotapi.Contact) error {
	filter := bson.M{"telegram_user_id": contact.UserID}
	update := bson.M{
		"$set": bson.M{
			"phone": contact.PhoneNumber,
		},
	}
	_, err := u.UserCol.UpdateOne(context.TODO(), filter, update)
	return err

}

func (u *userRepository) SearchUser(query string) (*ShortUser, error) {
	filter := bson.M{"$or": bson.A{bson.M{"username": query}, bson.M{"phone": query}}}
	res := u.UserCol.FindOne(context.TODO(), filter)
	if err := res.Err(); err != nil {
		return nil, err
	}
	var user ShortUser
	err := res.Decode(&user)
	return &user, err
}
