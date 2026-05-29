package service

import (
	"context"

	"github.com/Konshtantin/go-rest/internal/model"
	"github.com/Konshtantin/go-rest/internal/repository"
)

type PostService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(ctx context.Context, post *model.Post) error {
	return s.repo.Create(ctx, post)
}

func (s *PostService) GetByID(ctx context.Context, id int) (*model.Post, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PostService) List(ctx context.Context) ([]*model.Post, error) {
	return s.repo.List(ctx)
}

func (s *PostService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}