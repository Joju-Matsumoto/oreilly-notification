package notifyrecentlyaddedbooks

import (
	"context"
	"fmt"
	"testing"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
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

func Test_usecase_NotifyRecentlyAddedBooks(t *testing.T) {
	ErrSomething := fmt.Errorf("something wrong")
	type fields struct {
		updateRepositoryUsecase updateRepositoryUsecase
		newBookNotifier         newBookNotifier
	}
	type wants struct {
		err error
	}
	tests := []struct {
		name   string
		fields fields
		wants  wants
	}{
		{
			name: "success",
			fields: fields{
				updateRepositoryUsecase: &updateRepositoryUsecaseMock{
					UpdateRepositoryFunc: func(ctx context.Context) ([]*model.Book, error) {
						return books(t, 3), nil
					},
				},
				newBookNotifier: &newBookNotifierMock{
					NewBookFunc: func(ctx context.Context, books ...*model.Book) error {
						if len(books) != 3 {
							return fmt.Errorf("len(books) expected 3, but got %d", len(books))
						}
						return nil
					},
				},
			},
			wants: wants{
				err: nil,
			},
		},
		{
			name: "fail when updateRepositoryUsecase returned error",
			fields: fields{
				updateRepositoryUsecase: &updateRepositoryUsecaseMock{
					UpdateRepositoryFunc: func(ctx context.Context) ([]*model.Book, error) {
						return nil, ErrSomething
					},
				},
				newBookNotifier: nil,
			},
			wants: wants{
				err: ErrSomething,
			},
		},
		{
			name: "success when no added books found",
			fields: fields{
				updateRepositoryUsecase: &updateRepositoryUsecaseMock{
					UpdateRepositoryFunc: func(ctx context.Context) ([]*model.Book, error) {
						return nil, nil
					},
				},
				newBookNotifier: nil,
			},
			wants: wants{
				err: nil,
			},
		},
		{
			name: "fail when newBookNotifier.NewBook returned error",
			fields: fields{
				updateRepositoryUsecase: &updateRepositoryUsecaseMock{
					UpdateRepositoryFunc: func(ctx context.Context) ([]*model.Book, error) {
						return books(t, 3), nil
					},
				},
				newBookNotifier: &newBookNotifierMock{
					NewBookFunc: func(ctx context.Context, books ...*model.Book) error {
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
			u := NewUsecase(tt.fields.updateRepositoryUsecase, tt.fields.newBookNotifier)
			err := u.NotifyRecentlyAddedBooks(context.Background())
			require.Equal(t, tt.wants.err, err)
		})
	}
}
