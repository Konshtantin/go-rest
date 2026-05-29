package repository

import (
    "context"

    "github.com/Konshtantin/go-rest/internal/model"
)

type PostRepository interface {
    Create(ctx context.Context, post *model.Post) error
    GetByID(ctx context.Context, id int) (*model.Post, error)
    List(ctx context.Context) ([]*model.Post, error)
    Delete(ctx context.Context, id int) error
}