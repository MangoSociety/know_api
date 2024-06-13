package telegram

import (
	"context"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	e "know_api/pkg_1/error"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
	bot      *tgbotapi.BotAPI
}

const (
	getUpdatesMethod  = "getUpdates"
	sendMessageMethod = "sendMessage"

	formatMessageMd = "MarkdownV2"
)

func NewTelegramClient(host string, token string, tgBot *tgbotapi.BotAPI) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
		bot:      tgBot,
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(ctx context.Context, offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(ctx, getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendMessageMd(ctx context.Context, chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)
	q.Add("parse_mode", "markdown")

	_, err := c.doRequest(ctx, sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) SendMessageWithSpoilerMd(ctx context.Context, chatID int64, title string, text string) error {
	result := replacingString(title) + "\n" + "||" + replacingString(text) + "||"

	msg := tgbotapi.NewMessage(chatID, result)
	msg.ParseMode = formatMessageMd

	// TODO("мб что-то надо вытаскивать из ответа отправки и сохранять для альнейшей модерации ответа")
	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}

	return nil
}

func replacingString(data string) string {
	symbols := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
	var result = data
	for _, value := range symbols {
		result = strings.ReplaceAll(result, value, "\\"+value)
	}
	return result
}

func (c *Client) doRequest(ctx context.Context, method string, query url.Values) (data []byte, err error) {
	defer func() { err = e.WrapIfErr("can't do request", err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
