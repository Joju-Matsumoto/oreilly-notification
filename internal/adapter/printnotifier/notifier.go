package printnotifier

import (
	"context"
	"fmt"
	"io"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/notifier"
)

func New(w io.Writer) *printBookNotifier {
	return &printBookNotifier{
		w: w,
	}
}

type printBookNotifier struct {
	w io.Writer
}

// NewBook implements notifier.BookNotifier.
func (p *printBookNotifier) NewBook(ctx context.Context, books ...*model.Book) error {
	for _, book := range books {
		fmt.Fprintf(p.w, fmt.Sprintf("『%s』 %s\n", book.Title(), book.URL()))
	}
	return nil
}

var _ notifier.BookNotifier = (*printBookNotifier)(nil)
