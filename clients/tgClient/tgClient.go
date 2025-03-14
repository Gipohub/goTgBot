package tgClient

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/Gipohub/goTgBot/lib/e"
)

type Client struct {
	host     string
	basePath string
	owner    string
	client   http.Client
}

const (
	errReqMsg         = "cant do request"
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string, owner string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		owner:    owner,
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

// офсет число первого апдейта с которого мы бы хотели получить обновление
// (видимо из общего числа апдейтов?)
// лимит количество получаемых апдейтов
func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.Wrap("can't get updates", err) }()

	q := url.Values{}
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, e.Wrap(errReqMsg, err)
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

// берем первые два поля в сенд месседж методе телеги chat_id и text
func (c *Client) SendMesages(chatID int, text string) error {
	q := url.Values{}
	//fmt.Println(chatID)
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("cant send message", err)
	}

	return nil
}

// construct and send the message to chatID user, with topMessage text on the top of message,
// buttonsTextAndCallback button by button,
// and rangeLines is a sequence of numbers for the number of buttons in a row
func (c *Client) SendButtons(chatID int, topMessage string, buttonsTextAndCallback []InlineButton, rangeLines []int) error {
	count := len(buttonsTextAndCallback)
	if count == 0 {
		return e.WrapNew("no buttons provided")
	}

	var result InlineKeyboard
	i := 0
	for _, line := range rangeLines {
		end := i + line
		if end > count { // Проверяем границы
			end = count
		}
		result.RowsKeyboard = append(result.RowsKeyboard, buttonsTextAndCallback[i:end])
		i = end
	}

	replyMarkup, err := json.Marshal(result)
	if err != nil {
		return e.Wrap("failed to marshal buttons: %w", err)
	}

	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", topMessage)
	q.Add("reply_markup", string(replyMarkup))

	_, err = c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("cant send message", err)
	}

	return nil
}

func (c *Client) GetOwner() string {
	return c.owner
}

// код для отправки запроса аналогичен, поэтому отдельно
func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errReqMsg, err)
	}

	//метод енкод (query.Encode()) приведет параметры к такому виду,
	//которые мы сможем отправлять на сервер
	req.URL.RawQuery = query.Encode()

	//собственно запрос
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(errReqMsg, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(errReqMsg, err)
	}

	return body, nil
}
