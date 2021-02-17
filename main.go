package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gitlab.com/telegram/wishlist/common"
	"gitlab.com/telegram/wishlist/service/telegram"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//
//func main() {
//	// подключаемся к боту с помощью токена
//	bot, err := tgbotapi.NewBotAPI("1478590262:AAF_lrTTUwz2Xr-b07Kg-Z4RZwdfw33opWw")
//	if err != nil {
//		log.Panic(err)
//	}
//	bot.Debug = true
//	log.Printf("Authorized on account %s", bot.Self.UserName)
//
//	// инициализируем канал, куда будут прилетать обновления от API
//	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
//	ucfg.Timeout = 60
//	updates, err := bot.GetUpdatesChan(ucfg)
//	// читаем обновления из канала
//
//	for update := range updates {
//		if update.Message == nil { // ignore any non-Message Updates
//			continue
//		}
//
//		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
//
//		UserName := update.Message.From.UserName
//
//		// ID чата/диалога.
//		// Может быть идентификатором как чата с пользователем
//		// (тогда он равен UserID) так и публичного чата/канала
//		ChatID := update.Message.Chat.ID
//
//		// Текст сообщения
//		Text := update.Message.Text
//
//		log.Printf("[%s] %d %s", UserName, ChatID, Text)
//		switch update.Message.Command() {
//		case "help":
//			Text = "I understand /sayhi and /status."
//		case "sayhi":
//			Text = "Hi :)"
//		case "status":
//			Text = "I'm ok."
//		default:
//			Text = "I don't know that command"
//		}
//		// Ответим пользователю его же сообщением
//		reply := Text
//		// Созадаем сообщение
//		msg := tgbotapi.NewMessage(ChatID, reply)
//
//		bot.Send(msg)
//	}
//}

func main() {
	if err := common.Init(); err != nil {
		logrus.Fatal(err)
	}
	botApi, err := tgbotapi.NewBotAPI(common.Instance.Bot.Token)
	if err != nil {
		log.Panic(err)
	}
	mongoClient := initMongoClient()
	redisClient := initRedisClient()
	logrus.Info("authorized bot =", botApi.Self.UserName)
	bot := telegram.NewBot(botApi)

	if err := bot.Start(mongoClient, redisClient); err != nil {
		logrus.Error("main bot start:", err)
	}
}
func initMongoClient() *mongo.Client {
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s", common.Instance.Server.MongoUsername, common.Instance.Server.MongoPassword, common.Instance.Server.MongoURL)
	if common.Instance.Server.MongoUsername == "" {
		mongoURI = fmt.Sprintf("mongodb://%s", common.Instance.Server.MongoURL)
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	logrus.Info("Connecting to Mongo at:", common.Instance.Server.MongoURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Connected to Mongo:", common.Instance.Server.MongoURL)
	return client
}

func initRedisClient() *redis.Client {
	address := common.Instance.Server.RedisURL
	logrus.Info("Connecting to Redis at:", address)
	client := redis.NewClient(&redis.Options{Addr: address})
	_, err := client.Ping().Result() // check connection
	if err != nil {
		logrus.Fatalf("Error connecting to redis with address: %v :%v", address, err)
	}
	return client
}
