package cache

import (
	"context"

	"github.com/llan0/go-social/internal/store"
	"github.com/stretchr/testify/mock"
)

func NewMockStorage() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	mock.Mock
}

func (s *MockUserStore) Get(context.Context, int64) (*store.User, error) {
	return &store.User{}, nil
}
func (s *MockUserStore) Set(context.Context, *store.User) error {
	return nil
}
