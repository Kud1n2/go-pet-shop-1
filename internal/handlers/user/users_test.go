package user

import (
	"context"
	"errors"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi"
)

func TestCreateUser_Success(t *testing.T) {
	mock := &UsersMock{
		CreateUserFunc: func(ctx context.Context, user models.Customer) error {
			return nil
		},
	}

	body := `{
		"name":"bob",
		"email":"bob@gmail.net"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateUser(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateUser_BadRequest(t *testing.T) {
	mock := &UsersMock{
		CreateUserFunc: func(ctx context.Context, user models.Customer) error {
			return errors.New("Bad request")
		},
	}

	body := `{
		"name":bob,
		"email": ""
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateUser(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateUser_Error(t *testing.T) {
	mock := &UsersMock{
		CreateUserFunc: func(ctx context.Context, user models.Customer) error {
			return errors.New("DB error")
		},
	}

	body := `{
		"name":"bob",
		"email":"bob@gmail.net"
	}`

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateUser(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	mock := &UsersMock{
		GetUserByEmailFunc: func(ctx context.Context, email string) (models.Customer, error) {
			return models.Customer{
				ID: 1, Name: "Bob", Email: "example@email.com",
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users/example@email.com", nil)
	w := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "example@email.com")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	handler := New(slog.Default(), mock)
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetUserByEmail_BadRequest(t *testing.T) {
	mock := &UsersMock{
		GetUserByEmailFunc: func(ctx context.Context, email string) (models.Customer, error) {
			return models.Customer{}, errors.New("Bad request")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users/example@email.com", nil)
	w := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	handler := New(slog.Default(), mock)
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetUserByEmail_Error(t *testing.T) {
	mock := &UsersMock{
		GetUserByEmailFunc: func(ctx context.Context, email string) (models.Customer, error) {
			return models.Customer{}, errors.New("DB error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users/example@email.com", nil)
	w := httptest.NewRecorder()

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("email", "example@email.com")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	handler := New(slog.Default(), mock)
	handler.GetUserByEmail(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetAllUsers_Success(t *testing.T) {
	mock := &UsersMock{
		GetAllUsersFunc: func(ctx context.Context) ([]models.Customer, error) {
			return []models.Customer{
				{ID: 1, Name: "Bob", Email: "example@email.com"},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllUsers(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetAllUsers_Error(t *testing.T) {
	mock := &UsersMock{
		GetAllUsersFunc: func(ctx context.Context) ([]models.Customer, error) {
			return nil, errors.New("DB error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllUsers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
