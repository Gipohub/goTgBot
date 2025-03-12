package tgClient

//	телега отравляет ответ буль ок и
//  если он тру то полезную нагрузку полями резалт
type UpdatesResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

//пара полей из запроса апдейт к айпиай
type Update struct {
	ID int `json:"update_id"`
	//поле составное, содержит лишь некоторые,
	//необходимые сейчас поля обьекта
	//фром от кого, чат для обратной отправки, текст команды и ссылки
	Message  *IncomingMessage `json:"message"`
	Callback *CallbackQuery   `json:"callback_query"`
}

// входящее сообщение а не отправленное нами
type IncomingMessage struct {
	Text string `json:"text"`
	From From   `json:"from"`
	Chat Chat   `json:"chat"`
}

//только необходимое поле структуры
type From struct {
	Username string `json:"username"`
}

//только необходимое поле структуры
type Chat struct {
	ID int `json:"id"`
}

type CallbackQuery struct {
	From    From    `json:"from"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}

type Message struct {
	Chat Chat `json:"chat"`
}

// type buttonsTextAndCallback struct {
// 	Text string
// }

// Структура для inline-кнопок
type InlineKeyboard struct {
	RowsKeyboard [][]InlineButton `json:"inline_keyboard"`
}

type InlineButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}
