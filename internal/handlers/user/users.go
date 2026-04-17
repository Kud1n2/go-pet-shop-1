package user

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Customer interface {
	CreateUser(ctx context.Context, user models.Customer) error
	GetUserByEmail(ctx context.Context, email string) (models.Customer, error)
	GetAllUsers(ctx context.Context) ([]models.Customer, error)
}

type Handler struct {
	log     *slog.Logger
	storage Customer
}

func New(log *slog.Logger, storage Customer) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.user.GetAllUsers"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	users, err := h.storage.GetAllUsers(r.Context())

	if err != nil {
		log.Error("failed to get users", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to get users",
		})
		return
	}

	log.Info("Retrieved users successfully",
		slog.String("url", r.URL.String()),
		slog.Int("count", len(users)),
	)

	render.JSON(w, r, users)
}

func (h *Handler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.user.GetUserByEmail"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	email := r.URL.String()[7:]
	if email == "" {
		log.Error("Empty url email")
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Request must have email",
		})
		return
	}

	var user models.Customer
	var err error

	if user, err = h.storage.GetUserByEmail(r.Context(), email); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") ||
			strings.Contains(strings.ToLower(err.Error()), "no rows") ||
			strings.Contains(strings.ToLower(err.Error()), "rows affected: 0") {
			log.Warn("product not found", slog.String("email", email))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{
				"error":   "Not found",
				"message": fmt.Sprintf("User with email %s does not exist", email),
				"email":   email,
			})
			return
		}
		log.Error("Failed to get user by email", slog.Any("error", err))
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to get user by email",
		})
		return
	}

	render.JSON(w, r, user)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	const fn = "handlers.user.CreateUser"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Creating new user", slog.String("url", r.URL.String()))

	var user models.Customer

	if err := render.DecodeJSON(r.Body, &user); err != nil {
		log.Error("failed to decode request body", slog.Any("error", err))
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Failed to decode request body",
		})
		return
	}

	if user.Name == "" {
		log.Error("user name is empty")
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "user name is empty",
		})
		return
	}

	if user.Email == "" {
		log.Error("user email is empty")
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "user email is empty",
		})
		return
	}

	err := h.storage.CreateUser(r.Context(), user)
	if err != nil {
		log.Error("failed to create user")
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "failed to create user",
		})
		return
	}

	log.Info("User created successfully",
		slog.String("name", user.Name),
		slog.String("email", user.Email),
	)

	render.JSON(w, r, map[string]interface{}{
		"status": "User created successfully",
		"name":   user.Name,
		"email":  user.Email,
	})
}
