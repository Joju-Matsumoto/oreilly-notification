package updaterepository

import (
	"context"
	"errors"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/repository"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/webapi"
	"github.com/Joju-Matsumoto/oreilly-notification/pkg/oreillyapi"
)

const (
	FormatBook            = "book"
	PublisherOreillyJapan = "O'Reilly Japan, Inc."
	LanguageJa            = "ja"
	Limit                 = 100
)

//go:generate moq -out usecase_moq_test.go . bookRepository bookWebAPI

type bookRepository interface {
	repository.GetBookRepository
	repository.SaveBookRepository
}

type bookWebAPI interface {
	webapi.SearchBookWebAPI
}

func NewUsecase(bookWebAPI bookWebAPI, bookRepository bookRepository) *usecase {
	return &usecase{
		bookWebAPI: bookWebAPI,
		repository: bookRepository,
	}
}

type usecase struct {
	bookWebAPI bookWebAPI
	repository bookRepository
}

// UpdateRepository implements NotifyUsecase.
func (u *usecase) UpdateRepository(ctx context.Context) ([]*model.Book, error) {
	// 1. webapiで最近追加された本リストを取得
	// 2. repositoryに保存されていない本があれば保存し、リストに追加
	// 3. 2で保存されていない本がなくなるまで1,2を繰り返す
	// 4. リストを返す

	opt := oreillyapi.SearchOption{
		Formats: []string{
			FormatBook,
		},
		Languages: []string{
			LanguageJa,
		},
		Publishers: []string{
			PublisherOreillyJapan,
		},
		Sort:  oreillyapi.DateAdded,
		Order: oreillyapi.Desc,
		Page:  0,
		Limit: Limit,
	}

	var recentlyAddedBooks []*model.Book
	for {
		books, err := u.bookWebAPI.Search(ctx, opt)
		if err != nil {
			return nil, err
		}
		if len(books) == 0 {
			// api returned no books
			break
		}
		var newBookFound bool
		for _, book := range books {
			if _, err := u.repository.Get(ctx, book.ID()); err == nil {
				// found in repository -> do noting
				continue
			} else if !errors.Is(err, repository.ErrBookNotFound) {
				// unexpected error
				return nil, err
			}

			// not found -> save repository and add list
			newBookFound = true
			if err := u.repository.Save(ctx, book); err != nil {
				return nil, err
			}
			recentlyAddedBooks = append(recentlyAddedBooks, book)
		}
		if !newBookFound {
			// no new books found
			break
		}
		// search next page
		opt.Page += 1
	}

	return recentlyAddedBooks, nil
}

var _ Usecase = (*usecase)(nil)
