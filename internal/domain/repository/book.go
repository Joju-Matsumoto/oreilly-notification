package repository

import (
	"context"
	"fmt"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
)

var (
	ErrBookNotFound = fmt.Errorf("book not found")
)

type (
	BookRepository interface {
		GetBookRepository
		SaveBookRepository
		ListBookRepository
	}

	GetBookRepository interface {
		Get(ctx context.Context, id string) (*model.Book, error)
	}

	SaveBookRepository interface {
		Save(ctx context.Context, book *model.Book) error
	}

	ListBookOption struct {
	}
	ListBookRepository interface {
		List(ctx context.Context, opt ListBookOption) ([]*model.Book, error)
	}
)
