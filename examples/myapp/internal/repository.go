package internal

import (
	"context"

	"github.com/peractio/gdk/examples/myapp/internal/entity"
)

// UserRepository represents a repository for managing users.
type UserRepository interface {
	// Returns a user by ID.
	FindUserByID(ctx context.Context, id int) (*entity.User, error)

	// Returns a user by username.
	FindUserByUsername(ctx context.Context, username string) (*entity.User, error)

	// Creates a new user.
	CreateUser(ctx context.Context, user *entity.User) error
}
