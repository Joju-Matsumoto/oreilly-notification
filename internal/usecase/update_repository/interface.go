package updaterepository

import (
	"context"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
)

type Usecase interface {
	UpdateRepository(ctx context.Context) ([]*model.Book, error)
}
