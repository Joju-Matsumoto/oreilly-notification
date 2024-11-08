package updaterepository

import (
	"context"
	"fmt"
	"testing"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/repository"
	"github.com/Joju-Matsumoto/oreilly-notification/pkg/oreillyapi"
	"github.com/stretchr/testify/require"
)

func books(t *testing.T, size int) []*model.Book {
	t.Helper()

	books := make([]*model.Book, size)
	for i := range size {
		books[i] = model.NewBookFromOreillyResult(oreillyapi.Result{
			Id:    fmt.Sprintf("book_id_%d", i),
			Title: fmt.Sprintf("book_title_%d", i),
		})
	}
	return books
}

func Test_usecase_UpdateRepository(t *testing.T) {
	ErrSomething := fmt.Errorf("something wrong")

	type fields struct {
		bookWebAPI bookWebAPI
		repository bookRepository
	}
	type wants struct {
		books []*model.Book
		err   error
	}
	tests := []struct {
		name   string
		fields fields
		wants  wants
	}{
		{
			name: "success when webapi returned books only page 0, all books not found in repository",
			fields: fields{
				bookWebAPI: &bookWebAPIMock{
					SearchFunc: func(ctx context.Context, opt oreillyapi.SearchOption) ([]*model.Book, error) {
						if opt.Page > 0 {
							return nil, nil
						}
						return books(t, opt.Limit), nil
					},
				},
				repository: &bookRepositoryMock{
					GetFunc: func(ctx context.Context, id string) (*model.Book, error) {
						return nil, fmt.Errorf("%w: id='%s'", repository.ErrBookNotFound, id)
					},
					SaveFunc: func(ctx context.Context, book *model.Book) error {
						return nil
					},
				},
			},
			wants: wants{
				books: books(t, Limit),
				err:   nil,
			},
		},
		{
			name: "success when bookWebAPI returns books only page 0, all books found in repository",
			fields: fields{
				bookWebAPI: &bookWebAPIMock{
					SearchFunc: func(ctx context.Context, opt oreillyapi.SearchOption) ([]*model.Book, error) {
						return books(t, opt.Limit), nil
					},
				},
				repository: &bookRepositoryMock{
					GetFunc: func(ctx context.Context, id string) (*model.Book, error) {
						return &model.Book{}, nil
					},
					SaveFunc: func(ctx context.Context, book *model.Book) error {
						return fmt.Errorf("expect not to be called")
					},
				},
			},
			wants: wants{
				books: nil,
				err:   nil,
			},
		},
		{
			name: "fail when bookWebAPI.Search returns error",
			fields: fields{
				bookWebAPI: &bookWebAPIMock{
					SearchFunc: func(ctx context.Context, opt oreillyapi.SearchOption) ([]*model.Book, error) {
						return nil, ErrSomething
					},
				},
				repository: nil,
			},
			wants: wants{
				err: ErrSomething,
			},
		},
		{
			name: "fail when bookRepository.Get returns unexpected error",
			fields: fields{
				bookWebAPI: &bookWebAPIMock{
					SearchFunc: func(ctx context.Context, opt oreillyapi.SearchOption) ([]*model.Book, error) {
						return books(t, opt.Limit), nil
					},
				},
				repository: &bookRepositoryMock{
					GetFunc: func(ctx context.Context, id string) (*model.Book, error) {
						return nil, ErrSomething
					},
					SaveFunc: func(ctx context.Context, book *model.Book) error {
						return fmt.Errorf("expect not to be called")
					},
				},
			},
			wants: wants{
				err: ErrSomething,
			},
		},
		{
			name: "fail when bookRepository.Save returns error",
			fields: fields{
				bookWebAPI: &bookWebAPIMock{
					SearchFunc: func(ctx context.Context, opt oreillyapi.SearchOption) ([]*model.Book, error) {
						return books(t, opt.Limit), nil
					},
				},
				repository: &bookRepositoryMock{
					GetFunc: func(ctx context.Context, id string) (*model.Book, error) {
						return nil, fmt.Errorf("%w: id='%s'", repository.ErrBookNotFound, id)
					},
					SaveFunc: func(ctx context.Context, book *model.Book) error {
						return ErrSomething
					},
				},
			},
			wants: wants{
				err: ErrSomething,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUsecase(tt.fields.bookWebAPI, tt.fields.repository)
			got, err := u.UpdateRepository(context.Background())
			require.ErrorIs(t, err, tt.wants.err)
			if tt.wants.err != nil {
				return
			}

			require.Equal(t, len(got), len(tt.wants.books))
			for i, want := range tt.wants.books {
				require.Equal(t, want, got[i])
			}
		})
	}
}
