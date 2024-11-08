package jsonrepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/model"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/repository"
)

func New(path string) *jsonBookRepository {
	return &jsonBookRepository{
		path: path,
		data: map[string]model.Book{},
	}
}

type jsonBookRepository struct {
	path string
	data map[string]model.Book
}

func (j *jsonBookRepository) Open() error {
	if _, err := os.Stat(j.path); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(j.path)
		if err != nil {
			return err
		}
		defer f.Close()
		return nil
	}
	f, err := os.Open(j.path)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)

	if err := dec.Decode(&j.data); err != nil {
		return err
	}
	return nil
}

func (j *jsonBookRepository) Close() error {
	f, err := os.Create(j.path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(&j.data); err != nil {
		return err
	}
	return err
}

// Get implements repository.BookRepository.
func (j *jsonBookRepository) Get(ctx context.Context, id string) (*model.Book, error) {
	if book, ok := j.data[id]; ok {
		return &book, nil
	}
	return nil, fmt.Errorf("%w: id='%s'", repository.ErrBookNotFound, id)
}

// Save implements repository.BookRepository.
func (j *jsonBookRepository) Save(ctx context.Context, book *model.Book) error {
	j.data[book.ID()] = *book
	return nil
}

// List implements repository.BookRepository.
func (j *jsonBookRepository) List(ctx context.Context, opt repository.ListBookOption) ([]*model.Book, error) {
	books := make([]*model.Book, 0, len(j.data))
	for _, book := range j.data {
		books = append(books, &book)
	}
	return books, nil
}

var _ repository.BookRepository = (*jsonBookRepository)(nil)
