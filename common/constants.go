package common

const (
	BotCommandStart = "start"
)

//keyboards
const (
	BotMessageGetWishlistList    = "📝 Мой wishlist"
	BotMessageAddWishlistItem    = "❤️ Добавить в wishlist"
	BotMessageSearchUserWishlist = "🔍 Найти wishlist друга"
	BotMessageBackToMain         = "🔙 На главную"
	BotMessageAccessPhoneNumber  = "⚫️ Дать разрешение"
	BotMessageCopyWishlistLink   = "📌 Скопировать ссылку"
	BotMessageWishlistItemName   = "🖊Ввести название"
	BotMessageWishlistItemPhoto  = "Загрузить фотографию"
	BotMessageWishlistItemUrl    = "Указать ссылку где можно купить"
)

//answer keyboards
const (
	BotAnswerMessageSearchUserWishlist        = "Введите номер или username вашего друга"
	BotAnswerMessageAddWishlist               = "Введите название или добавьте фотографию(также можно отправить ссылку)."
	BotAnswerMessageAccessPhoneNumberQuestion = "Я умею показывать твой вишлист твоим друзьям по номеру телефона. Для этого необходимо дать разрешение"
)

const (
	SearchUserTypeSender    = "sender"
	SearchUserTypeRecipient = "recipient"
)

const (
	ActionTypeSearchUser = "search_user"
)
