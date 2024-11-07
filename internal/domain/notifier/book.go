package notifier

import (
	"context"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
)

type (
	BookNotifier interface {
		NewBookNotifier
	}

	NewBookNotifier interface {
		NewBook(ctx context.Context, books ...*model.Book) error
	}
)
