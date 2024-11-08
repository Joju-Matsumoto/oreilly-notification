package notifyrecentlyaddedbooks

import (
	"context"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/notifier"
	updaterepository "github.com/Joju-Matsumoto/oreilly-notification/internal/usecase/update_repository"
)

func NewUsecase(updateRepositoryUsecase updateRepositoryUsecase, newBookNotifier newBookNotifier) *usecase {
	return &usecase{
		updateRepositoryUsecase: updateRepositoryUsecase,
		newBookNotifier:         newBookNotifier,
	}
}

//go:generate moq -out usecase_moq_test.go . updateRepositoryUsecase newBookNotifier

type updateRepositoryUsecase interface {
	updaterepository.Usecase
}

type newBookNotifier interface {
	notifier.NewBookNotifier
}

type usecase struct {
	updateRepositoryUsecase updateRepositoryUsecase
	newBookNotifier         newBookNotifier
}

// NotifyRecentlyAddedBooks implements NotifyRecentlyAddedBooksUsecase.
func (u *usecase) NotifyRecentlyAddedBooks(ctx context.Context) error {
	recentlyAddedBooks, err := u.updateRepositoryUsecase.UpdateRepository(ctx)
	if err != nil {
		return err
	}

	if len(recentlyAddedBooks) == 0 {
		// no added book found
		return nil
	}

	if err := u.newBookNotifier.NewBook(ctx, recentlyAddedBooks...); err != nil {
		return err
	}
	return nil
}

var _ Usecase = (*usecase)(nil)
