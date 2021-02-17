package telegram

import (
	"fmt"
	"github.com/go-redis/redis"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gitlab.com/telegram/wishlist/common"
	"gitlab.com/telegram/wishlist/service/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{
		bot: bot,
	}
}
func (b *Bot) Start(client *mongo.Client, redisClient *redis.Client) error {
	db.Init(client, redisClient)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	b.bot.Debug = true
	updates, _ := b.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore non-Message updates
			continue
		}
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				logrus.Errorf("handleCommand: chat_id = %v,  err = %v", update.Message.Chat.ID, err)
				return b.handleErrorMessage(update.Message, common.BotErrorInvalidHandler)
			}
			//continue
		} else if update.Message.Text != "" {
			if err := b.handleMessage(update.Message); err != nil {
				logrus.Errorf("handleMessage: chat_id = %v,  err = %v", update.Message.Chat.ID, err)
				return b.handleErrorMessage(update.Message, common.BotErrorInvalidHandler)

			}
		} else {
			if err := b.handleUnknownMessage(update.Message); err != nil {
				logrus.Errorf("handleUnknownMessage: chat_id = %v,  err = %v", update.Message.Chat.ID, err)
				return b.handleErrorMessage(update.Message, common.BotErrorInvalidHandler)
			}
		}
	}

	return nil
}

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case common.BotCommandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.Text = fmt.Sprintf("Привет, %v. Данный бот поможет тебе составить твой список желаний или посмотреть список желаний твоих друзей!", message.From.UserName)
	msg.ReplyMarkup = mainKeyboard
	_, err := b.bot.Send(msg)
	go func() {
		user, _ := db.UserRepo.GetUser(message.From.ID)
		if user == nil {
			if _, err := db.UserRepo.CreateUser(message); err != nil {
				logrus.Errorf("handleStartCommand:CreateUser: userID = %v,  err = %v", message.From.ID, err)
			}
		}

	}()
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Я не знаю такой команды")
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	switch message.Text {
	case common.BotMessageBackToMain:
		return b.handleStartCommand(message)
	case common.BotMessageGetWishlistList:
		return b.getUserWishlist(message)
	case common.BotMessageAddWishlistItem:
		return b.saveSessionInRedisByCreateWishlist(message, common.BotAnswerMessageAddWishlist)
	case common.BotMessageSearchUserWishlist:
		return b.saveSessionInRedisByCreateWishlist(message, common.BotAnswerMessageSearchUserWishlist)

	default:
		return b.handleUnknownMessage(message)
	}
}

func (b *Bot) handleUnknownMessage(message *tgbotapi.Message) error {
	if message.Contact != nil {

	}
	prevMessage, err := db.FavoriteRepo.GetSessionByCreateWishlist(message)
	if err != nil {
		logrus.Errorf("insertItemForWishlist:GetSessionByCreateWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
		return b.handleErrorMessage(message, common.BotErrorInvalidHandler)
	} else {
		if prevMessage == common.BotAnswerMessageAddWishlist {
			return b.insertItemForWishlist(message)
		} else if prevMessage == common.BotAnswerMessageSearchUserWishlist {
			return b.searchWishlistByPhoneNumber(message)
		} else if prevMessage == common.BotAnswerMessageAccessPhoneNumberQuestion && message.Contact != nil {
			if err := b.shareWishlistLink(message); err != nil {
				logrus.Errorf("shareWishlistLink: chat_id = %v,  err = %v", message.Chat.ID, err)
				return b.handleErrorMessage(message, common.BotErrorInvalidHandler)
			}
			return nil
		} else {
			msg := tgbotapi.NewMessage(message.Chat.ID, "Я понимаю только команды. Попробуем пообщаться через них:)")
			msg.ReplyMarkup = mainKeyboard
			_, err = b.bot.Send(msg)
			return err
		}
	}
}

func (b *Bot) handleErrorMessage(message *tgbotapi.Message, text string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, text)
	msg.ReplyMarkup = mainKeyboard
	_, err := b.bot.Send(msg)
	return err
}
