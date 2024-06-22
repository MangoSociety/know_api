package telegram

import (
	"context"
	"errors"
	"fmt"
	"github.com/MangoSociety/know_api/clients/telegram"
	categoriesService "github.com/MangoSociety/know_api/internal/categories/service"
	notesService "github.com/MangoSociety/know_api/internal/notes/service"
	spheresService "github.com/MangoSociety/know_api/internal/spheres/service"
	domainStatistics "github.com/MangoSociety/know_api/internal/statistics/domain"
	statisticsService "github.com/MangoSociety/know_api/internal/statistics/service"
	"github.com/MangoSociety/know_api/internal/telegram/events"
	"github.com/MangoSociety/know_api/internal/user/service"
	e "github.com/MangoSociety/know_api/pkg/error"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Processor struct {
	tg                   *telegram.Client
	noteService          notesService.NoteService
	categoryService      categoriesService.CategoryService
	sphereService        spheresService.SphereService
	userSelectionService service.UserSelectionService
	statisticsService    statisticsService.StatisticsService
	offset               int
}

type Meta struct {
	ChatID   int
	Username string
	Data     string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func NewProcessor(
	client *telegram.Client,
	noteService notesService.NoteService,
	categoriesService categoriesService.CategoryService,
	sphereService spheresService.SphereService,
	userSelectionService service.UserSelectionService,
	statService statisticsService.StatisticsService,

) *Processor {
	return &Processor{
		tg:                   client,
		noteService:          noteService,
		categoryService:      categoriesService,
		sphereService:        sphereService,
		userSelectionService: userSelectionService,
		statisticsService:    statService,
	}
}

func (p *Processor) Fetch(ctx context.Context, limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(ctx, p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(ctx context.Context, event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(ctx, event)
	case events.CallbackQuery:
		return p.processChooserButton(ctx, event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(ctx, event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func (p *Processor) processChooserButton(ctx context.Context, event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return err
	}

	if meta.Data == "questions:get" {
		categories := p.userSelectionService.GetCategoriesByChatID(meta.ChatID)
		note, _ := p.noteService.GetRandomNoteByCategory(categories)
		p.sendQuestion(meta.ChatID, *note)
		return nil
	}

	questionId, isAnswerAfterQuestion := getSubstringAfterPrefix(meta.Data, "questions:send:")
	if isAnswerAfterQuestion {
		noteHex := questionId
		note, _ := p.noteService.GetByHex(noteHex)
		p.sendAnswer(meta.ChatID, *note)
		return nil
	}

	number, objectID, err := extractNumberAndID(meta.Data, "questions:")
	if err == nil {
		if number == 5 {

		} else {
			fmt.Printf("Error extracting number and ID: %v\n", err)
			p.statisticsService.Create(domainStatistics.Statistics{
				Status:    number,
				NoteID:    objectID,
				CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
			})
			p.sendHelp(context.Background(), meta.ChatID, "Красавчик, еще 1 вопрос в помойку, что дальше?😎\n")
			return nil
		}
	}

	categoryName, isRoot := extractCategoryID(meta.Data)
	result, _ := p.categoryService.GetCategoriesTree(categoryName)
	if isRoot {
		if len(result) == 0 {
			p.saveChooseCategory(meta.ChatID, categoryName)
			p.chooserCategory(meta.ChatID)
			return nil
		}
		return p.showCategories(ctx, meta.ChatID, result)
	} else {
		if meta.Data == "choose_category:all" {
			p.chooserCategory(meta.ChatID)
			return nil
		}
		if err := p.doCmd(ctx, meta.Data, meta.ChatID, meta.Username); err != nil {
			return err
		}
	}

	return nil
}

func startsWithQuestionsSend(input string) bool {
	const prefix = "questions:send:"
	return strings.HasPrefix(input, prefix)
}

func getSubstringAfterPrefix(input, prefix string) (string, bool) {
	if strings.HasPrefix(input, prefix) {
		return input[len(prefix):], true
	}
	return "", false
}

func (p *Processor) saveChooseCategory(chatId int, categoryName string) {
	category, _ := p.categoryService.GetCategoriesByName(categoryName)
	sphere, _ := p.sphereService.GetById(category.SphereID.String())
	log.Println("sphere = ", sphere)
	log.Println("category = ", category)
	p.userSelectionService.AddCategoryToUser(chatId, sphere.ID, sphere.Name, category.ID, categoryName)
}

func extractCategoryID(input string) (string, bool) {
	re := regexp.MustCompile(`^choose_category:([^:]+)$`)
	if !re.MatchString(input) || strings.Contains(input, "root") {
		return "", false
	}

	matches := re.FindStringSubmatch(input)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// Функция для извлечения цифры и ID из строки
func extractNumberAndID(input, prefix string) (int, primitive.ObjectID, error) {
	if !strings.HasPrefix(input, prefix) {
		return 0, primitive.NilObjectID, fmt.Errorf("string does not start with prefix: %s", prefix)
	}

	// Удаляем префикс
	suffix := input[len(prefix):]

	// Находим индекс разделителя ":answer:"
	parts := strings.Split(suffix, ":answer:")
	if len(parts) != 2 {
		return 0, primitive.NilObjectID, fmt.Errorf("invalid string format")
	}

	// Извлекаем цифру и ID
	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, primitive.NilObjectID, fmt.Errorf("invalid number format: %w", err)
	}

	objectID, err := primitive.ObjectIDFromHex(parts[1])
	if err != nil {
		return 0, primitive.NilObjectID, fmt.Errorf("invalid hex ID: %w", err)
	}

	return number, objectID, nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	switch updType {
	case events.Message:
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}

	case events.CallbackQuery:
		res.Meta = Meta{
			ChatID:   upd.CallbackQuery.Message.Chat.ID,
			Username: upd.CallbackQuery.From.Username,
			Data:     upd.CallbackQuery.Data,
		}
	default:
		panic("unhandled default case")
	}
	//if updType == events.Message {
	//	res.Meta = Meta{
	//		ChatID:   upd.Message.Chat.ID,
	//		Username: upd.Message.From.Username,
	//	}
	//}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.CallbackQuery != nil {
		return events.CallbackQuery
	}
	if upd.Message != nil {
		return events.Message
	}

	return events.Unknown
}
