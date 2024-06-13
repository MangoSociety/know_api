package telegram

import (
	"context"
	e "github.com/MangoSociety/know_api/pkg/error"
	"log"
	"strings"
)

const (
	HelpCmd = "/help"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	switch text {
	case HelpCmd:
		return p.sendHelp(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) getRandomQuestionAndroid(ctx context.Context, chatId int) error {
	//data, err := p.storage.GetNote(ctx, "category", "Object")
	//if err != nil {
	//	fmt.Println("error" + err.Error())
	//	return err
	//}
	//err = p.tg.SendMessageWithSpoilerMd(ctx, int64(chatId), data.Title, data.Content)
	//if err != nil {
	//	fmt.Println("error" + err.Error())
	//	return err
	//}
	return nil
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	const op = "commands.savePage"
	log.Println(op, "got new page to save", pageURL, "from", username)

	if err := p.tg.SendMessage(ctx, chatID, "Нужный функционал удален"); err != nil {
		return err
	}

	return nil
}

func (p *Processor) saveNote(ctx context.Context, chatId int) error {
	//text := p.quest_usecase.z()
	//log.Println("get random note = ", text)
	//
	//note := &storage.Note{
	//	Title:    text.Title,
	//	Sphere:   text.Sphere,
	//	Category: text.Category,
	//	Content:  text.Content,
	//}
	//
	//return p.storage.SaveNote(ctx, note)
	return nil
}

func (p *Processor) sendHelp(ctx context.Context, chatID int) error {
	//err := p.questionsUC.AutoMigration(ctx)
	//if err != nil {
	//	return err
	//}
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}
