package product

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

// Get Product - Ready
func TestGetAllProducts_Success(t *testing.T) {
	// Мокаем storage — он вернёт один продукт.
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{
				{ID: 1, Name: "Dog Food"},
			}, nil
		},
	}

	// Создаем HTTP-запрос GET /products
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	// Создаем хендлер с мок-хранилищем
	handler := New(slog.Default(), mock)

	// Вызываем метод GetAllProducts, который является http.HandlerFunc
	handler.GetAllProducts(w, req)

	// Проверяем HTTP-код
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}
func TestGetAllProducts_Error(t *testing.T) {
	// Мокаем storage — он будет возвращать ошибку
	mock := &ProductsMock{
		GetAllProductsFunc: func(ctx context.Context) ([]models.Product, error) {
			return nil, errors.New("DB error")
		},
	}

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetAllProducts(w, req)

	// Ожидаем HTTP 500
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Create Product
// =======================

func TestCreateProduct_Success(t *testing.T) {
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 1, nil
		},
	}
	prod := `{
		"name": "banana",
		"price": 12.5,
		"stock": 4
	}`
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(prod))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	// TODO: Написать Unit-тест для создания продукта с невалидным JSON (400 Bad Request)
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 0, errors.New("Empty product name")
		},
	}

	prod := `{
		"name": name,
		"price": 12.5,
		"stock": 4
	}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(prod))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для создания продукта при ошибке сервиса (500 Internal Server Error)
	mock := &ProductsMock{
		CreateProductFunc: func(ctx context.Context, product models.Product) (int, error) {
			return 0, errors.New("Bad request")
		},
	}

	prod := `{
		"name": "boss",
		"price": 12.5,
		"stock": 4
	}`

	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(prod))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Update Product
// =======================

func TestUpdateProduct_Success(t *testing.T) {
	// TODO: Написать Unit-тест для обновления продукта (200 OK)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return nil
		},
	}
	prod := `{
		"name": "boss",
		"price": 12.5,
		"stock": 4
	}`

	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(prod))
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdateProduct_BadRequest(t *testing.T) {
	// TODO: Написать Unit-тест для обновления продукта с невалидным JSON (400 Bad Request)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return errors.New("Bad Request")
		},
	}
	prod := `{
		"name": boss,
		"price": 12.5,
		"stock": 4
	}`

	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(prod))
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для обновления продукта при ошибке сервиса (500 Internal Server Error)
	mock := &ProductsMock{
		UpdateProductFunc: func(ctx context.Context, product models.Product) error {
			return errors.New("DB error")
		},
	}
	prod := `{
		"name": "boss",
		"price": 12.5,
		"stock": 4
	}`

	req := httptest.NewRequest(http.MethodPut, "/products/1", strings.NewReader(prod))
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.UpdateProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// =======================
// Delete Product
// =======================

func TestDeleteProduct_Success(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта (200 OK)
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return nil
		},
	}

	req := httptest.NewRequest(http.MethodPut, "/products/1", nil)
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта с пустым id (400 Bad Request)
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return errors.New("Bad request")
		},
	}

	req := httptest.NewRequest(http.MethodDelete, "/products/abc", nil)
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "abc")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteProduct_Fail(t *testing.T) {
	// TODO: Написать Unit-тест для удаления продукта при ошибке сервиса (500 Internal Server Error)
	mock := &ProductsMock{
		DeleteProductFunc: func(ctx context.Context, id int) error {
			return errors.New("DB error")
		},
	}

	req := httptest.NewRequest(http.MethodPut, "/products/1", nil)
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.DeleteProduct(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetProductByID_Success(t *testing.T) {
	mock := &ProductsMock{
		GetProductsByIDFunc: func(ctx context.Context, id int) (models.Product, error) {
			return models.Product{}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.GetProductsByID(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetProductByID_BadRequest(t *testing.T) {
	mock := &ProductsMock{
		GetProductsByIDFunc: func(ctx context.Context, id int) (models.Product, error) {
			return models.Product{}, errors.New("Bad Request")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "abc")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.GetProductsByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetProductByID_Fail(t *testing.T) {
	mock := &ProductsMock{
		GetProductsByIDFunc: func(ctx context.Context, id int) (models.Product, error) {
			return models.Product{}, errors.New("DB error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.GetProductsByID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
