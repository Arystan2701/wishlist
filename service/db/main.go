package db

import (
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	databaseWishlist = "wishlist"
)

const (
	collectionUsers     = "users"
	collectionFavorites = "favorites"
	collectionActions   = "actions"
)

var UserRepo *userRepository
var FavoriteRepo *FavoriteRepository
var ActionRepo *actionRepo

func Init(client *mongo.Client, redisClient *redis.Client) {
	UserRepo = NewUserRepository(client, redisClient)
	FavoriteRepo = NewFavoriteRepo(client, redisClient)
	ActionRepo = NewActionRepo(client, redisClient)

}
