package db

import (
	"context"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gitlab.com/telegram/wishlist/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"strconv"
	"time"
)

type FavoriteRepository struct {
	Col   *mongo.Collection
	Redis *redis.Client
}

func NewFavoriteRepo(client *mongo.Client, redis *redis.Client) *FavoriteRepository {
	col := client.Database(databaseWishlist).Collection(collectionFavorites)
	return &FavoriteRepository{
		Col:   col,
		Redis: redis,
	}
}

const (
	KeyOtpCustomer = "otp:customer:"
)

type Wishlist struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title" bson:"title"`
	AuthorID  int                `json:"author_id" bson:"author_id"`
	ChatID    int64              `json:"chat_id" bson:"chat_id"`
	CreatedAt int64              `json:"created_at" bson:"created_at"`
	Image     string             `json:"image" bson:"image"`
}

func (f *FavoriteRepository) GetUserWishlist(userID int) ([]Wishlist, error) {
	filter := bson.M{"author_id": userID}
	cursor, err := f.Col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	wishlists := make([]Wishlist, 0)
	for cursor.Next(context.TODO()) {
		var wishlist Wishlist
		if err := cursor.Decode(&wishlist); err != nil {
			logrus.Error("FavoriteRepository:GetUserWishlist:", err)
			continue
		}

		wishlists = append(wishlists, wishlist)
	}
	return wishlists, nil
}

func (f *FavoriteRepository) InsertItemForWishlist(message *tgbotapi.Message) (*Wishlist, error) {
	item := &Wishlist{
		ID:        primitive.NewObjectID(),
		Title:     message.Text,
		AuthorID:  message.From.ID,
		ChatID:    message.Chat.ID,
		CreatedAt: common.Timestamp(),
	}
	if message.Photo != nil {
		photos := *message.Photo
		item.Title = message.Caption
		item.Image = photos[0].FileID
	}
	_, err := f.Col.InsertOne(context.TODO(), item)
	return item, err
}

func (f *FavoriteRepository) SaveSessionByCreateWishlist(userID int, message *tgbotapi.Message) error {
	_, err := f.Redis.Set(KeyOtpCustomer+strconv.Itoa(userID), message.Text, 5*time.Minute).Result()
	return err
}

func (f *FavoriteRepository) GetSessionByCreateWishlist(message *tgbotapi.Message) (string, error) {
	messageIDStr, err := f.Redis.Get(KeyOtpCustomer + strconv.Itoa(message.From.ID)).Result()
	if err != nil {
		return "", err
	}
	return messageIDStr, err
}

func (f *FavoriteRepository) RemoveOldMessageID(userID int) interface{} {
	return f.Redis.Del(KeyOtpCustomer + strconv.Itoa(userID)).Err()
}

//func (f *FavoriteRepository) UpdateUserPhone(contact *tgbotapi.Contact) error {
//	filter := bson.M{"author.id": contact.UserID}
//	update := bson.M{
//		"$set": bson.M{
//			"author.phone": contact.PhoneNumber,
//		},
//	}
//	_, err := f.Col.UpdateMany(context.TODO(), filter, update)
//	return err
//}

func (f *FavoriteRepository) GetUserWishlistByPhone(userID int) ([]Wishlist, error) {
	filter := bson.M{"author_id": userID}
	cursor, err := f.Col.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	wishlists := make([]Wishlist, 0)
	for cursor.Next(context.TODO()) {
		var wishlist Wishlist
		if err := cursor.Decode(&wishlist); err != nil {
			logrus.Error("FavoriteRepository:GetUserWishlistByPhone:", err)
			continue
		}

		wishlists = append(wishlists, wishlist)
	}
	return wishlists, nil
}
