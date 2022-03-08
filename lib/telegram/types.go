package telegram

type User struct {
	Id int `json:"id"`
}

type Chat struct {
	Id   int    `json:"id"`
	Type string `json:"type"`
}

type MessageEntity struct {
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Type   string `json:"type"`
}

type Message struct {
	Text string `json:"text"`
}

type IncomingMessage struct {
	Message
	MessageId int             `json:"message_id"`
	Chat      Chat            `json:"chat"`
	From      User            `json:"user"`
	Entities  []MessageEntity `json:"entities"`
	Date      int             `json:"date"`
}

type KeyboardButton struct {
	RequestLocation bool `json:"request_location"`
}

type ReplyKeyboardMarkup struct {
	Keyboard []KeyboardButton `json:"keyboard"`
}

type OutgoingMessage struct {
	Message
	ChatId      int                 `json:"chat_id"`
	ReplyMarkup ReplyKeyboardMarkup `json:"reply_markup"`
}

type Update struct {
	UpdateId int             `json:"update_id"`
	Message  IncomingMessage `json:"message"`
}

type SendMessageResponse struct {
	Result IncomingMessage `json:"result"`
}
