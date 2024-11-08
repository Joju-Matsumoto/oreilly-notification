package webapi

import (
	"context"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
	"github.com/Joju-Matsumoto/oreilly-notification/pkg/oreillyapi"
)

type (
	BookWebAPI interface {
		SearchBookWebAPI
	}

	SearchBookWebAPI interface {
		Search(ctx context.Context, opt oreillyapi.SearchOption) ([]*model.Book, error)
	}
)
