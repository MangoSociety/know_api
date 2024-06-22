package telegram

import (
	"context"
	"github.com/MangoSociety/know_api/internal/categories/domain"
	notesDomain "github.com/MangoSociety/know_api/internal/notes/domain"
	e "github.com/MangoSociety/know_api/pkg/error"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

const (
	HelpCmd           = "/help"
	ChooseSphereCmd   = "choose_sphere:root"
	ChooseCategoryCmd = "choose_category:root"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	switch text {
	case HelpCmd:
		return p.sendHelp(ctx, chatID, "Привет! Я бот, который поможет тебе учиться!")
	case ChooseSphereCmd:
		return p.showRootCategories(ctx, chatID)
	case ChooseCategoryCmd:
		return p.showRootCategories(ctx, chatID)
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

func (p *Processor) sendHelp(ctx context.Context, chatID int, previewText string) error {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Кидай вопрос", "questions:get"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Выберите категорию", "choose_category:root"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Выберите тему", "choose_sphere:root"),
		},
	}
	return p.tg.SendButton(int64(chatID), previewText+"Ты можешь...", buttons)
}

func (p *Processor) showRootCategories(ctx context.Context, chatID int) error {
	rootCategories, err := p.categoryService.GetCategoriesTree("")
	//categories, err := p.questionsUC.GetRootCategories(ctx)
	if err != nil {
		return err
	}
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, item := range rootCategories {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(item.Name, "choose_category:"+item.Name),
		})
	}
	return p.tg.SendButton(int64(chatID), "Выбери 1 категорию", buttons)
}

// ! Знаю ! Надо повторить ! Впервые вижу ! Что-то слышал !
func (p *Processor) sendQuestion(chatID int, note notesDomain.Note) {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Впервые вижу", "questions:send:"+note.ID.Hex()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Что-то слышал", "questions:send:"+note.ID.Hex()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Надо повторить", "questions:send:"+note.ID.Hex()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Знаю", "questions:send:"+note.ID.Hex()),
		},
	}
	p.tg.SendButton(int64(chatID), note.Title, buttons)
}

// ! Понял ! Есть вопрос к ответу ! Сложно для понимания ! Совсем ничего не понятно
func (p *Processor) sendAnswer(chatID int, note notesDomain.Note) {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Понял", "questions:1:answer:"+note.ID.Hex()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Есть вопрос к ответу", "questions:2:answer:"+note.ID.Hex()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Сложно для понимания", "questions:3:answer:"+note.ID.Hex()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Совсем ничего не понятно", "questions:4:answer:"+note.ID.Hex()),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Задать вопрос к ответу", "questions:5:answer:"+note.ID.Hex()),
		},
	}
	p.tg.SendButton(int64(chatID), note.Content, buttons)
}

func (p *Processor) sendQuestionByAnswer(chatID int) {
	p.tg.SendMessage(context.Background(), chatID, "Окей, у тебя есть вопросы к ответу, это нормально, расскажи, что именно тебе непонятно?")
}

func (p *Processor) sendText(chatID int, text string) {
	p.tg.SendMessage(context.Background(), chatID, text)
}

func (p *Processor) showCategories(ctx context.Context, chatID int, categories []*domain.Category) error {
	buttons := [][]tgbotapi.InlineKeyboardButton{}
	for _, item := range categories {
		buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(item.Name, "choose_category:"+item.Name),
		})
		log.Println(item)
	}
	buttons = append(buttons, []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Дальше не пойдем", "choose_category:all"),
	})
	return p.tg.SendButton(int64(chatID), "Выбери 1 категорию", buttons)
}

func (p *Processor) chooserCategory(chatID int) {
	err := p.tg.SendMessage(context.Background(), chatID, "Поздравляю, ты выбрал категории! Что дальше?")
	if err != nil {
		return
	}
	err = p.showStandartMenu(context.Background(), chatID)
	if err != nil {
		return
	}
}

func (p *Processor) showStandartMenu(ctx context.Context, chatID int) error {
	buttons := [][]tgbotapi.InlineKeyboardButton{
		{
			tgbotapi.NewInlineKeyboardButtonData("Выбрать категорию", "choose_category:root"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Выбрать тему", "choose_sphere:root"),
		},
		{
			tgbotapi.NewInlineKeyboardButtonData("Получить рандомный вопрос", "questions:get"),
		},
	}
	return p.tg.SendButton(int64(chatID), "Ты можешь...", buttons)
}

func (p *Processor) sendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}
