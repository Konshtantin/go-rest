package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Konshtantin/go-rest/internal/model"
	"github.com/Konshtantin/go-rest/internal/repository"
	"github.com/Konshtantin/go-rest/internal/service"
)

type PostHandler struct {
	svc *service.PostService
	log *slog.Logger
}

func NewPostHandler(svc *service.PostService, log *slog.Logger) *PostHandler {
	return &PostHandler{svc: svc, log: log}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.Create(r.Context(), &post); err != nil {
		h.log.Error("create post", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	post, err := h.svc.GetByID(r.Context(), id)
	
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		h.log.Error("get post", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) List(w http.ResponseWriter, r *http.Request) {
	posts, err := h.svc.List(r.Context())
	if err != nil {
		h.log.Error("list posts", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if posts == nil {
		posts = []*model.Post{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "post not found", http.StatusNotFound)
			return
		}
		h.log.Error("delete post", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}