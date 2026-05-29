package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Konshtantin/go-rest/internal/model"
)

type postgresUserRepo struct {
    pool *pgxpool.Pool
}

func NewPostgresUserRepo(pool *pgxpool.Pool) UserRepository {
    return &postgresUserRepo{pool: pool}
}

func (r *postgresUserRepo) Create(ctx context.Context, user *model.User) error {
    query := `
        INSERT INTO users (name, email)
        VALUES ($1, $2)
        RETURNING id, created_at`

    return r.pool.QueryRow(ctx, query, user.Name, user.Email).
        Scan(&user.ID, &user.CreatedAt)
}

func (r *postgresUserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
    query := `SELECT id, name, email, created_at FROM users WHERE id = $1`

    user := &model.User{}
    err := r.pool.QueryRow(ctx, query, id).
        Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

    if errors.Is(err, pgx.ErrNoRows) {
        return nil, ErrNotFound
    }
    return user, err
}

func (r *postgresUserRepo) List(ctx context.Context) ([]*model.User, error) {
    query := `SELECT id, name, email, created_at FROM users ORDER BY id`

    rows, err := r.pool.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*model.User
    
    for rows.Next() {
        u := &model.User{}
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, rows.Err()
}

func (r *postgresUserRepo) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM users WHERE id = $1`

    result, err := r.pool.Exec(ctx, query, id)
    if err != nil {
        return err
    }
    
    if result.RowsAffected() == 0 {
        return ErrNotFound
    }

    return nil
}