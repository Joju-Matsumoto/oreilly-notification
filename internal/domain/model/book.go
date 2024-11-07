package model

import "github.com/Joju-Matsumoto/oreilly-notification/pkg/oreillyapi"

type Book struct {
	Result oreillyapi.Result
}

func NewBookFromOreillyResult(result oreillyapi.Result) *Book {
	return &Book{
		Result: result,
	}
}

func (b *Book) ID() string {
	return b.Result.Id
}

func (b *Book) Title() string {
	return b.Result.Title
}
