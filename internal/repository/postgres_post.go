package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Konshtantin/go-rest/internal/model"
)

type postgresPostRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresPostRepo(pool *pgxpool.Pool) PostRepository {
	return &postgresPostRepo{pool: pool}
}

func (r *postgresPostRepo) Create(ctx context.Context, post *model.Post) error {
	query := `
		INSERT INTO posts (title, body, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query, post.Title, post.Body, post.UserID).
		Scan(&post.ID, &post.CreatedAt)
}

func (r *postgresPostRepo) GetByID(ctx context.Context, id int) (*model.Post, error) {
	query := `SELECT id, title, body, user_id, created_at FROM posts WHERE id = $1`

	post := &model.Post{}
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&post.ID, &post.Title, &post.Body, &post.UserID, &post.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return post, err
}

func (r *postgresPostRepo) List(ctx context.Context) ([]*model.Post, error) {
	query := `SELECT id, title, body, user_id, created_at FROM posts ORDER BY id`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		p := &model.Post{}
		if err := rows.Scan(&p.ID, &p.Title, &p.Body, &p.UserID, &p.CreatedAt); err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

func (r *postgresPostRepo) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM posts WHERE id = $1`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}