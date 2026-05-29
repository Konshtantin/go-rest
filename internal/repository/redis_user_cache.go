package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Konshtantin/go-rest/internal/model"
)

type cachedUserRepo struct {
	repo   UserRepository
	client *redis.Client
	ttl    time.Duration
}

func NewCachedUserRepo(repo UserRepository, client *redis.Client) UserRepository {
	return &cachedUserRepo{
		repo:   repo,
		client: client,
		ttl:    5 * time.Minute,
	}
}

func (r *cachedUserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	key := fmt.Sprintf("user:%d", id)

	data, err := r.client.Get(ctx, key).Bytes()
	if err == nil {
		var user model.User
		if err := json.Unmarshal(data, &user); err == nil {
			return &user, nil
		}
	}

	user, err := r.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(user); err == nil {
		r.client.Set(ctx, key, data, r.ttl)
	}

	return user, nil
}

func (r *cachedUserRepo) invalidate(ctx context.Context, id int) {
	r.client.Del(ctx, fmt.Sprintf("user:%d", id))
}

func (r *cachedUserRepo) Create(ctx context.Context, user *model.User) error {
	return r.repo.Create(ctx, user)
}

func (r *cachedUserRepo) List(ctx context.Context) ([]*model.User, error) {
	return r.repo.List(ctx)
}

func (r *cachedUserRepo) Delete(ctx context.Context, id int) error {
	err := r.repo.Delete(ctx, id)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return err
	}

	r.invalidate(ctx, id)

	// TODO: need to invalidate all cached users on delete (now only by ID)

	return err
}