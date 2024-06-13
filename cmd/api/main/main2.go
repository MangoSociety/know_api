package main

//
//import (
//	"fmt"
//	"github.com/gomarkdown/markdown/ast"
//	"github.com/russross/blackfriday/v2"
//	"golang.org/x/net/html"
//	"io"
//	"regexp"
//	"strings"
//	"time"
//)
//
////import (
////	"fmt"
////	"github.com/gomarkdown/markdown"
////	"github.com/gomarkdown/markdown/html"
////	"io"
////	"regexp"
////	"strings"
////	"time"
////
////	"github.com/yuin/goldmark"
////	"github.com/yuin/goldmark/ast"
////	_ "github.com/yuin/goldmark/parser"
////	"github.com/yuin/goldmark/text"
////	_ "github.com/yuin/goldmark/util"
////)
//
//const testString = `#готово
//
//Theme : #common
//Title: Расскажи что такое SharedPreferences Какие данные можно хранить Какие плюсы и минусы
//Sphere: #android
//
//### Content
//
//### Какие данные можно хранить
//
//С SharedPreferences, вы можете хранить следующие типы данных:
//
//- boolean
//- float
//- int
//- long
//- String
//- Set<String> (используется для хранения коллекций строк, например, списка значений)
//
//Эти типы данных покрывают большинство нужд приложения в сохранении простых конфигурационных параметров и пользовательских предпочтений.
//
//### Плюсы SharedPreferences
//
//1. **Простота использования**: Интерфейс SharedPreferences прост и интуитивно понятен, что делает его легким в использовании для хранения и извлечения простых данных.
//2. **Легкость интеграции**: Он хорошо интегрирован в Android SDK, предоставляя прямой и эффективный способ сохранения легковесных данных.
//3. **Быстрый доступ**: Доступ к данным, хранящимся в SharedPreferences, осуществляется быстро, что делает его подходящим для хранения данных, необходимых при старте приложения.
//4. **Поддержка асинхронного сохранения**: С API 9 (Android 2.3, Gingerbread) и выше, SharedPreferences предлагает метод apply(), который асинхронно сохраняет изменения, минимизируя задержки UI.
//
//### Минусы SharedPreferences
//
//1. **Ограниченный объем и типы данных**: SharedPreferences подходит только для примитивных типов данных и не предназначен для хранения сложных объектов или больших объемов данных.
//2. **Безопасность**: Данные, сохраненные в SharedPreferences, хранятся в виде обычных файлов XML без шифрования, что делает их уязвимыми для атак, если устройство скомпрометировано.
//3. **Отсутствие поддержки структурированных данных**: SharedPreferences не подходит для хранения структурированных данных или реализации сложных иерархий настроек.
//4. **Проблемы с многопоточностью**: При неправильном использовании может возникнуть состояние гонки или другие проблемы с многопоточностью, особенно если одновременно производится чтение и запись из разных потоков.
//
//SharedPreferences является удобным инструментом для хранения небольших объемов данных, таких как настройки пользователя или простые флаги состояния. Однако для более сложных или объемных данных следует рассмотреть другие варианты хранения данных, такие как базы данных SQLite или хранилище на основе файлов с использованием внутренней или внешней памяти.
//
//### External Link
//
//-
//
//### Internal Link
//
//- ....
//`
//
////package main
//
////import (
////"fmt"
////"io"
////"regexp"
////"strings"
////
////"github.com/gomarkdown/markdown"
////"github.com/gomarkdown/markdown/ast"
////"github.com/gomarkdown/markdown/html"
////)
//
//// levels tracks how deep we are in a heading "structure"
//var levels []int
//
//func hasLevels() bool {
//	return len(levels) > 0
//}
//
//func lastLevel() int {
//	if hasLevels() {
//		return levels[len(levels)-1]
//	}
//	return 0
//}
//
//func popLevel() int {
//	level := lastLevel()
//	levels = levels[:len(levels)-1]
//	return level
//}
//
//func pushLevel(x int) {
//	levels = append(levels, x)
//}
//
//var reID = regexp.MustCompile(`\s+`)
//
//// renderSections catches an ast.Heading node, and wraps the node
//// and its "children" nodes in <section>...</section> tags; there's no
//// real hierarchy in Markdown, so we make one up by saying things like:
//// - H2 is a child of H1, and so forth from 1 → 2 → 3 ... → N
//// - an H1 is a sibling of another H1
//func renderSections(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
//	openSection := func(level int, id string) {
//		w.Write([]byte(fmt.Sprintf("<section id=\"%s\">\n", id)))
//		pushLevel(level)
//	}
//	closeSection := func() {
//		w.Write([]byte("</section>\n"))
//		popLevel()
//	}
//
//	if _, ok := node.(*ast.Heading); ok {
//		level := node.(*ast.Heading).Level
//		if entering {
//			// close heading-sections deeper than this level; we've "come up" some number of levels
//			for lastLevel() > level {
//				closeSection()
//			}
//
//			txtNode := node.GetChildren()[0]
//			if _, ok := txtNode.(*ast.Text); !ok {
//				panic(fmt.Errorf("expected txtNode to be *ast.Text; got %T", txtNode))
//			}
//			headTxt := string(txtNode.AsLeaf().Literal)
//			id := strings.ToLower(reID.ReplaceAllString(headTxt, "-"))
//
//			openSection(level, id)
//		}
//	}
//
//	// at end of document
//	if _, ok := node.(*ast.Document); ok {
//		if !entering {
//			for hasLevels() {
//				closeSection()
//			}
//		}
//	}
//
//	// continue as normal
//	return ast.GoToNext, false
//}
//
//// Section представляет раздел внутри контента Markdown
//type Section struct {
//	Title   string `bson:"title"`
//	Content string `bson:"content"`
//}
//
//// MarkdownData структура для хранения данных из Markdown
//type MarkdownData struct {
//	Theme        string    `bson:"theme"`
//	Title        string    `bson:"title"`
//	Sphere       string    `bson:"sphere"`
//	Sections     []Section `bson:"sections"`
//	ExternalLink string    `bson:"external_link"`
//	InternalLink string    `bson:"internal_link"`
//	Date         time.Time `bson:"date"`
//}
//
//func main() {
//	data, _ := parseString(testString)
//	fmt.Printf("%+v\n", data)
//}
//
//func parseString(content string) (*MarkdownData, error) {
//	output := blackfriday.Run([]byte(content))
//	doc, err := html.Parse(strings.NewReader(string(output)))
//	if err != nil {
//		return nil, err
//	}
//
//	sections, theme, title, sphere, externalLink, internalLink := parseHTMLContent(doc)
//
//	return &MarkdownData{
//		Theme:        theme,
//		Title:        title,
//		Sphere:       sphere,
//		Sections:     sections,
//		ExternalLink: externalLink,
//		InternalLink: internalLink,
//		Date:         time.Now(),
//	}, nil
//}
//
//func parseHTMLContent(n *html.Node) ([]Section, string, string, string, string, string) {
//	var sections []Section
//	var currentSection *Section
//	var theme, title, sphere, externalLink, internalLink string
//
//	var f func(*html.Node)
//	f = func(n *html.Node) {
//		if n.Type == html.ElementNode {
//			switch n.Data {
//			case "h1":
//				theme = getTextContent(n)
//			case "h2":
//				title = getTextContent(n)
//			case "h3":
//				sphere = getTextContent(n)
//			case "h4":
//				if currentSection != nil {
//					sections = append(sections, *currentSection)
//				}
//				currentSection = &Section{
//					Title:   getTextContent(n),
//					Content: "",
//				}
//			case "ul", "ol":
//				if currentSection != nil {
//					currentSection.Content += getListContent(n) + "\n"
//				}
//			case "strong":
//				if currentSection != nil {
//					currentSection.Content += formatBoldText(n) + "\n"
//				}
//			case "em":
//				if currentSection != nil {
//					currentSection.Content += formatItalicText(n) + "\n"
//				}
//			case "a":
//				if externalLink == "" && internalLink == "" {
//					href := getHref(n)
//					if strings.Contains(getTextContent(n), "External Link") {
//						externalLink = href
//					} else if strings.Contains(getTextContent(n), "Internal Link") {
//						internalLink = href
//					}
//				}
//			default:
//				if currentSection != nil {
//					currentSection.Content += getTextContent(n) + "\n"
//				}
//			}
//		}
//		for c := n.FirstChild; c != nil; c = c.NextSibling {
//			f(c)
//		}
//	}
//
//	f(n)
//
//	if currentSection != nil {
//		sections = append(sections, *currentSection)
//	}
//
//	return sections, theme, title, sphere, externalLink, internalLink
//}
//
//func getTextContent(n *html.Node) string {
//	if n.Type == html.TextNode {
//		return n.Data
//	}
//	var result string
//	for c := n.FirstChild; c != nil; c = c.NextSibling {
//		result += getTextContent(c)
//	}
//	return result
//}
//
//func getListContent(n *html.Node) string {
//	var result string
//	var f func(*html.Node)
//	f = func(n *html.Node) {
//		if n.Type == html.ElementNode && n.Data == "li" {
//			result += "- " + getTextContent(n) + "\n"
//		}
//		for c := n.FirstChild; c != nil; c = c.NextSibling {
//			f(c)
//		}
//	}
//	f(n)
//	return result
//}
//
//func formatBoldText(n *html.Node) string {
//	return "**" + getTextContent(n) + "**"
//}
//
//func formatItalicText(n *html.Node) string {
//	return "*" + getTextContent(n) + "*"
//}
//
//func getHref(n *html.Node) string {
//	for _, a := range n.Attr {
//		if a.Key == "href" {
//			return a.Val
//		}
//	}
//	return ""
//}
