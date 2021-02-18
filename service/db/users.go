package db

import (
	"context"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/telegram/wishlist/common"
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

type User struct {
	TelegramUserID  int    `json:"telegram_user_id" bson:"telegram_user_id"`
	FirstName       string `json:"first_name" bson:"first_name"`
	LastName        string `json:"last_name" bson:"last_name"`
	UserName        string `json:"username" bson:"username"`
	Phone           string `json:"phone" bson:"phone"`
	SearchReqCount  int    `json:"search_req_count" bson:"search_req_count"`
	SearchRespCount int    `json:"search_resp_count" bson:"search_resp_count"`
}

type ShortUser struct {
	TelegramUserID int    `json:"telegram_user_id" bson:"telegram_user_id"`
	UserName       string `json:"username" bson:"username"`
}

func (u *userRepository) CreateUser(message *tgbotapi.Message) (*User, error) {
	user := &User{
		TelegramUserID: message.From.ID,
		FirstName:      message.From.FirstName,
		LastName:       message.From.LastName,
		UserName:       message.From.UserName,
	}
	_, err := u.UserCol.InsertOne(context.TODO(), user)
	return user, err
}

func (u *userRepository) GetUser(userID int) (*User, error) {
	res := u.UserCol.FindOne(context.TODO(), bson.M{"telegram_user_id": userID})
	if err := res.Err(); err != nil {
		return nil, err
	}
	var user User
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

func (u *userRepository) SearchUser(query string) (*User, error) {
	filter := bson.M{"$or": bson.A{bson.M{"username": query}, bson.M{"phone": query}}}
	res := u.UserCol.FindOne(context.TODO(), filter)
	if err := res.Err(); err != nil {
		return nil, err
	}
	var user User
	err := res.Decode(&user)
	return &user, err
}

func (u *userRepository) IncrementUserSearchCountByType(userUD int, requestType string) error {
	filter := bson.M{"telegram_user_id": userUD}
	update := bson.M{}
	if requestType == common.SearchUserTypeSender {
		update["search_req_count"] = 1
	} else {
		update["search_resp_count"] = 1
	}
	_, err := u.UserCol.UpdateOne(context.TODO(), filter, bson.M{"$inc": update})
	return err
}
