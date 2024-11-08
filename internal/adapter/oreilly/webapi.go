package oreilly

import (
	"context"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/webapi"
	"github.com/Joju-Matsumoto/oreilly-notification/pkg/oreillyapi"
)

func New() *oreillyAPI {
	return &oreillyAPI{
		apiClient: oreillyapi.New(),
	}
}

type oreillyAPI struct {
	apiClient *oreillyapi.Client
}

// Search implements webapi.BookWebAPI.
func (o *oreillyAPI) Search(ctx context.Context, opt oreillyapi.SearchOption) ([]*model.Book, error) {
	res, err := o.apiClient.Search(opt)
	if err != nil {
		return nil, err
	}

	books := make([]*model.Book, len(res.Results))
	for i, result := range res.Results {
		books[i] = model.NewBookFromOreillyResult(result)
	}

	return books, nil
}

var _ webapi.BookWebAPI = (*oreillyAPI)(nil)
