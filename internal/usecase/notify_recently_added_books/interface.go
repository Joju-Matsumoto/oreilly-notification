package notifyrecentlyaddedbooks

import "context"

type Usecase interface {
	NotifyRecentlyAddedBooks(ctx context.Context) error
}
