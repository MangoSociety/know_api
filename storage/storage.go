package storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"read-adviser-bot/lib/e"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	SaveNote(ctx context.Context, p *Note) error
	GetNote(ctx context.Context, conditionField string, conditionValue string) (note *Note, err error)
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
}

type Note struct {
	Title    string `bson:"title"`
	Sphere   string `bson:"sphere"`
	Category string `bson:"category"`
	Content  string `bson:"content"`
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't calculate hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
