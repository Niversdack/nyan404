package sqlstore

import (
	"github.com/jmoiron/sqlx"
	"github.com/nyan404/internal/app/store"

	// pq ...
	_ "github.com/lib/pq"
)

// Store ...
type Store struct {
	db             *sqlx.DB
	userRepository *UserRepository
}

// New ...
func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository

}
