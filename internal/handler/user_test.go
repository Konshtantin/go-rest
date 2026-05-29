package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"

	"github.com/Konshtantin/go-rest/internal/handler"
	"github.com/Konshtantin/go-rest/internal/model"
	"github.com/Konshtantin/go-rest/internal/repository"
	"github.com/Konshtantin/go-rest/internal/service"
)

type mockUserRepo struct {
	users map[int]*model.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[int]*model.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	user.ID = len(m.users) + 1
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int) (*model.User, error) {
	user, ok := m.users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return user, nil
}

func (m *mockUserRepo) List(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id int) error {
	if _, ok := m.users[id]; !ok {
		return repository.ErrNotFound
	}

	delete(m.users, id)
	return nil
}

func newTestHandler() (*handler.UserHandler, *mockUserRepo) {
	repo := newMockUserRepo()
	svc := service.NewUserService(repo)
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	return handler.NewUserHandler(svc, log), repo
}

func TestUserHandler_GetByID_Found(t *testing.T) {
	h, repo := newTestHandler()
	repo.users[1] = &model.User{ID: 1, Name: "Konstantin", Email: "test@example.com"}


	r := chi.NewRouter()
	r.Get("/users/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var user model.User
	err := json.NewDecoder(rr.Body).Decode(&user)
	
	require.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Konstantin", user.Name)
}

func TestUserHandler_GetByID_NotFound(t *testing.T) {
	h, _ := newTestHandler()

	r := chi.NewRouter()
	r.Get("/users/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUserHandler_GetByID_InvalidID(t *testing.T) {
	h, _ := newTestHandler()

	r := chi.NewRouter()
	r.Get("/users/{id}", h.GetByID)

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}