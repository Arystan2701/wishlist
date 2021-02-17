package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"gitlab.com/telegram/wishlist/common"
)

var mainKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(common.BotMessageAddWishlistItem),
		tgbotapi.NewKeyboardButton(common.BotMessageGetWishlistList)),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(common.BotMessageSearchUserWishlist),
	),
)

var itemWishlist0 = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(common.BotMessageWishlistItemName),
		tgbotapi.NewKeyboardButton(common.BotMessageWishlistItemPhoto),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(common.BotMessageBackToMain)),
)

var accessPhoneNumberKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButtonContact(common.BotMessageAccessPhoneNumber),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(common.BotMessageBackToMain)),
)

var backMainKeyboards = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(common.BotMessageBackToMain)),
)
