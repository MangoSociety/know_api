package telegram

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"read-adviser-bot/storage"
	"strings"

	"read-adviser-bot/lib/e"
)

const (
	RndCmd        = "/rnd"
	HelpCmd       = "/help"
	StartCmd      = "/start"
	RndAndroidCmd = "/aa"
	RndGolangCmd  = "/gg"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, username)

	if isAddCmd(text) {
		return p.savePage(ctx, chatID, text, username)
	}

	switch text {
	case RndCmd:
		err := p.saveNote(ctx, chatID)
		if err != nil {
			return err
		}
		return p.sendRandom(ctx, chatID, username)
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	case StartCmd:
		return p.sendHello(ctx, chatID)
	case RndAndroidCmd:
		return p.getRandomQuestionAndroid(ctx, chatID)
	case RndGolangCmd:
		return p.getRandomQuestionGolang(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) getRandomQuestionAndroid(ctx context.Context, chatId int) error {
	data, err := p.storage.GetNote(ctx, "category", "Object")
	if err != nil {
		fmt.Println("error" + err.Error())
		return err
	}
	err = p.tg.SendMessageWithSpoilerMd(ctx, int64(chatId), data.Title, data.Content)
	if err != nil {
		fmt.Println("error" + err.Error())
		return err
	}
	//result := p.quest_usecase.GetRandomQuestionAndroid()
	//for _, value := range strings.Split(data, "NEXT") {
	//	err := p.tg.SendMessageMd(ctx, chatId, value)
	//if err != nil {
	//	fmt.Println("error" + err.Error())
	//	return err
	//}
	//}
	return nil
}

func (p *Processor) getRandomQuestionGolang(ctx context.Context, chatId int) error {
	//result := p.quest_usecase.GetRandomQuestionGolang()
	//for _, value := range strings.Split(result, "NEXT") {
	//	//err := p.tg.SendMessageWithSpoilerMd(ctx, chatId, value)
	//	if err != nil {
	//		fmt.Println("error" + err.Error())
	//		return err
	//	}
	//}
	//err := p.tg.SendMessage(ctx, chatId, "completed")
	//if err != nil {
	//	return err
	//}

	return nil
}

//func sendMessage(ctx context.Context, chatID int, text string) error {
//	client := http.Client{}
//	q := url.Values{}
//	q.Add("chat_id", strconv.Itoa(chatID))
//	q.Add("text", text)
//
//	_, err := client.Do.doRequest(ctx, sendMessageMethod, q)
//	if err != nil {
//		return e.Wrap("can't send message", err)
//	}
//
//	return nil
//}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	const op = "commands.savePage"
	log.Println(op, "got new page to save", pageURL, "from", username)

	//page := &storage.Page{
	//	URL:      pageURL,
	//	UserName: username,
	//}

	//isExists, err := p.storage.IsExists(ctx, page)
	//if err != nil {
	//	return err
	//}
	//if isExists {
	//	return p.tg.SendMessage(ctx, chatID, msgAlreadyExists)
	//}
	//
	//if err := p.storage.Save(ctx, page); err != nil {
	//	return err
	//}

	if err := p.tg.SendMessage(ctx, chatID, "Нужный функционал удален"); err != nil {
		return err
	}

	return nil
}

func (p *Processor) saveNote(ctx context.Context, chatId int) error {
	text := p.quest_usecase.GetRandomQuestionAndroidStruct()
	log.Println("get random note = ", text)

	note := &storage.Note{
		Title:    text.Title,
		Sphere:   text.Sphere,
		Category: text.Category,
		Content:  text.Content,
	}

	return p.storage.SaveNote(ctx, note)
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	//defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	//page, err := p.storage.PickRandom(ctx, username)
	//if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
	//	return err
	//}
	//if errors.Is(err, storage.ErrNoSavedPages) {
	//	return p.tg.SendMessage(ctx, chatID, msgNoSavedPages)
	//}
	//
	//if err := p.tg.SendMessage(ctx, chatID, "пока нужный функционал удален"); err != nil {
	//	return err
	//}

	//return p.storage.Remove(ctx, page)
	return nil
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
