package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"gitlab.com/telegram/wishlist/common"
	"gitlab.com/telegram/wishlist/service/db"
	"strconv"
)

func (b *Bot) getUserWishlist(message *tgbotapi.Message) error {
	wishlists, err := db.FavoriteRepo.GetUserWishlist(message.From.ID)
	if err != nil {
		logrus.Errorf("getUserWishlist:GetUserWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
		return err
	}
	if len(wishlists) > 0 {
		for i, w := range wishlists {
			if w.Image != "" {
				msg := tgbotapi.NewPhotoShare(message.Chat.ID, w.Image)
				msg.Caption = fmt.Sprintf("%v. %v", i+1, w.Title)
				_, err := b.bot.Send(msg)
				if err != nil {
					logrus.Errorf("getUserWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
					return err
				}
			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%v. %v", i+1, w.Title))
				//msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				//	tgbotapi.NewInlineKeyboardRow(
				//		tgbotapi.NewInlineKeyboardButtonData("Удалить", w.ID.Hex()),
				//	),
				//)

				_, err := b.bot.Send(msg)
				if err != nil {
					logrus.Errorf("getUserWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
					return err
				}
			}
		}
		user, err := db.UserRepo.GetUser(message.From.ID)
		if err != nil {
			logrus.Errorf("getUserWishlist:GetUser: chat_id = %v,  err = %v", message.Chat.ID, err)
			return b.handleErrorMessage(message, common.BotErrorInvalidHandler)
		}

		if user.Phone == "" {
			msg := tgbotapi.NewMessage(message.Chat.ID, common.BotAnswerMessageAccessPhoneNumberQuestion)
			msg.ReplyMarkup = accessPhoneNumberKeyboard
			newMsg, err := b.bot.Send(msg)
			if err != nil {
				logrus.Errorf("getUserWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
				return err
			}
			go func() {
				if err := db.FavoriteRepo.SaveSessionByCreateWishlist(message.From.ID, &newMsg); err != nil {
					logrus.Errorf("getUserWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
				}
			}()

		}
		return nil
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Твой список желаний пуст(( Давай исправим это!")
		msg.ReplyMarkup = mainKeyboard
		_, err := b.bot.Send(msg)
		if err != nil {
			logrus.Errorf("getUserWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
			return err
		}
	}
	return nil
}

func (b *Bot) addItemInWishlist(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, common.BotAnswerMessageAddWishlist)
	msg.ReplyMarkup = itemWishlist0
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) saveSessionInRedisByCreateWishlist(message *tgbotapi.Message, answer string) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, answer)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	newMessage, err := b.bot.Send(msg)
	if err := db.FavoriteRepo.SaveSessionByCreateWishlist(message.From.ID, &newMessage); err != nil {
		logrus.Errorf("saveSessionByCreateWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
	}
	return err
}

func (b *Bot) insertItemForWishlist(message *tgbotapi.Message) error {
	wishlist, err := db.FavoriteRepo.InsertItemForWishlist(message)
	if err != nil {
		logrus.Errorf("insertItemForWishlist:InsertItemForWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
		return err
	}
	if wishlist.Image != "" {
		msg := tgbotapi.NewPhotoShare(message.Chat.ID, wishlist.Image)
		msg.Caption = "Добавлен:" + wishlist.Title
		msg.ReplyMarkup = mainKeyboard
		_, err = b.bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Добавлен: "+wishlist.Title)
		msg.ReplyMarkup = mainKeyboard
		_, err = b.bot.Send(msg)
	}
	go b.removeOldMessageIDOnRedis(message)
	return err
}

func (b *Bot) removeOldMessageIDOnRedis(message *tgbotapi.Message) {
	if err := db.FavoriteRepo.RemoveOldMessageID(message.From.ID); err != nil {
		logrus.Error("insertItemForWishlist:RemoveOldMessageID:", err)
	}
}

//func (b *Bot) insertItemForWishlist(message *tgbotapi.Message) error {
//	if err := db.FavoriteRepo.InsertItemForWishlist(message); err != nil {
//		logrus.Errorf("insertItemForWishlist:InsertItemForWishlist: chat_id = %v,  err = %v", message.Chat.ID, err)
//		return err
//	}
//	photos := *message.Photo
//	url, err := b.bot.GetFileDirectURL(photos[0].FileID)
//	logrus.Info("url =            ", url, "               err =     ", err)
//	logrus.Info(b.bot.GetFile(tgbotapi.FileConfig{FileID: photos[0].FileID}))
//	logrus.Info(message.Text)
//	logrus.Info(message.Photo)
//	msg := tgbotapi.NewPhotoShare(message.Chat.ID, photos[0].FileID)
//	msg.ReplyMarkup = mainKeyboard
//	_, err = b.bot.Send(msg)
//	go func() {
//		if err := db.FavoriteRepo.RemoveOldMessageID(message.From.ID); err != nil {
//			logrus.Error("insertItemForWishlist:RemoveOldMessageID:", err)
//		}
//	}()
//	return err
//}

func (b *Bot) shareWishlistLink(message *tgbotapi.Message) error {
	if err := db.UserRepo.UpdateUserPhone(message.Contact); err != nil {
		return err
	}
	//msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Вот твоя ссылка на вишлист. скопиру ее и поделись с друзьями  https://9f8b51b92dec.ngrok.io?username=%v", message.From.UserName))
	msg := tgbotapi.NewMessage(message.Chat.ID, "Теперь твой вишлист будет виден твоим друзьям")
	msg.ReplyMarkup = mainKeyboard
	_, err := b.bot.Send(msg)
	go func() {
		b.removeOldMessageIDOnRedis(message)
	}()
	return err
}

func (b *Bot) searchWishlistByPhoneNumber(message *tgbotapi.Message) error {
	query := common.VerifySearchQuery(message.Text)
	if message.Contact != nil {
		query = message.Contact.PhoneNumber
	}
	go b.updateUserSearchCount(message.From.ID, common.SearchUserTypeSender)
	user, err := db.UserRepo.SearchUser(query)
	if err != nil {
		logrus.Errorf("searchWishlistByPhoneNumber:SearchUser: chat_id = %v,  err = %v", message.Chat.ID, err)
		go b.insertUserAction(common.ActionTypeSearchUser, message, query, nil)
		return b.handleErrorMessage(message, common.BotErrorNotFoundSearchWishlist)
	}
	go func() {
		b.updateUserSearchCount(user.TelegramUserID, common.SearchUserTypeRecipient)
		b.insertUserAction(common.ActionTypeSearchUser, message, query, user)
	}()
	wishlists, err := db.FavoriteRepo.GetUserWishlistByPhone(user.TelegramUserID)
	if len(wishlists) > 0 {
		text := fmt.Sprint("*Пользователь*:*" + user.UserName + "*\n_Количество хотелок_:*" + strconv.Itoa(len(wishlists)) + "*\n")
		msg := tgbotapi.NewMessage(message.Chat.ID, text)
		msg.ParseMode = "markdown"
		_, err := b.bot.Send(msg)
		if err != nil {
			logrus.Errorf("searchWishlistByPhoneNumber: chat_id = %v,  err = %v", message.Chat.ID, err)
			return err
		}
		for i, w := range wishlists {
			if w.Image != "" {
				msg := tgbotapi.NewPhotoShare(message.Chat.ID, w.Image)
				msg.Caption = fmt.Sprintf("%v. %v", i+1, w.Title)
				msg.ReplyMarkup = mainKeyboard
				_, err := b.bot.Send(msg)
				if err != nil {
					logrus.Errorf("searchWishlistByPhoneNumber: chat_id = %v,  err = %v", message.Chat.ID, err)
					return err
				}
			} else {
				msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%v. %v", i+1, w.Title))
				msg.ReplyMarkup = mainKeyboard
				_, err := b.bot.Send(msg)
				if err != nil {
					logrus.Errorf("searchWishlistByPhoneNumber: chat_id = %v,  err = %v", message.Chat.ID, err)
					return err
				}
			}
		}
		go b.removeOldMessageIDOnRedis(message)
	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Его список желаний пуст((")
		msg.ReplyMarkup = mainKeyboard
		_, err := b.bot.Send(msg)
		if err != nil {
			logrus.Errorf("searchWishlistByPhoneNumber: chat_id = %v,  err = %v", message.Chat.ID, err)
			return err
		}
	}
	return nil
}

func (b *Bot) updateUserSearchCount(senderID int, requestType string) {
	if err := db.UserRepo.IncrementUserSearchCountByType(senderID, requestType); err != nil {
		logrus.Errorf("updateUserActionSearchUsers:sender_id:%v, err: %v", senderID, err)
	}
}

func (b *Bot) insertUserAction(actionType string, message *tgbotapi.Message, query string, recipient *db.User) {
	if err := db.ActionRepo.InsertAction(actionType, message, query, recipient); err != nil {
		logrus.Errorf("insertUserAction:chat_id:%v, err: %v", message.Chat.ID, err)
	}
}
