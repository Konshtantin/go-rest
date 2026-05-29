package repository

import (
    "context"

    "github.com/Konshtantin/go-rest/internal/model"
)

type UserRepository interface {
    Create(ctx context.Context, user *model.User) error
    GetByID(ctx context.Context, id int) (*model.User, error)
    List(ctx context.Context) ([]*model.User, error)
    Delete(ctx context.Context, id int) error
}