package order

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

func TestCreateOrder(t *testing.T) {
	mock := &OrdersMock{
		CreateOrderFunc: func(ctx context.Context, order models.Order) (int, error) {
			return 1, nil
		},
	}

	order := `{
		"userID": 1,
		"total_price":0
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(order))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateOrder(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateOrder_BadRequest(t *testing.T) {
	mock := &OrdersMock{
		CreateOrderFunc: func(ctx context.Context, order models.Order) (int, error) {
			return 0, errors.New("Empty userID")
		},
	}

	order := `{
		"userID":,
		"total_price":0
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(order))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateOrder_Error(t *testing.T) {
	mock := &OrdersMock{
		CreateOrderFunc: func(ctx context.Context, order models.Order) (int, error) {
			return 0, errors.New("DB error")
		},
	}

	order := `{
		"userID":1,
		"total_price":0
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(order))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.CreateOrder(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestAddOrderItem_Success(t *testing.T) {
	mock := &OrdersMock{
		AddOrderItemFunc: func(ctx context.Context, orderItem models.OrderItem) error {
			return nil
		},
	}

	orderItem := `{
		"product_id": 1,
		"quantity": 1
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders/1/items", strings.NewReader(orderItem))
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.AddOrderItem(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestAddOrderItem_BadRequest(t *testing.T) {
	mock := &OrdersMock{
		AddOrderItemFunc: func(ctx context.Context, orderItem models.OrderItem) error {
			return errors.New("Bad Request")
		},
	}

	orderItem := `{
		"product_id": 1,
		"quantity": 1
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders/1/items", strings.NewReader(orderItem))
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "abc")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.AddOrderItem(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAddOrderItem_Error(t *testing.T) {
	mock := &OrdersMock{
		AddOrderItemFunc: func(ctx context.Context, orderItem models.OrderItem) error {
			return errors.New("DB error")
		},
	}

	orderItem := `{
		"product_id": 1,
		"quantity": 1
	}`

	req := httptest.NewRequest(http.MethodPost, "/orders/1/items", strings.NewReader(orderItem))
	w := httptest.NewRecorder()

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	handler := New(slog.Default(), mock)
	handler.AddOrderItem(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetOrderByID_Success(t *testing.T) {
	mock := &OrdersMock{
		GetOrderByIDFunc: func(ctx context.Context, id int) (models.Order, error) {
			return models.Order{}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	w := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routerCtx))

	handler := New(slog.Default(), mock)
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetOrderByID_BadRequest(t *testing.T) {
	mock := &OrdersMock{
		GetOrderByIDFunc: func(ctx context.Context, id int) (models.Order, error) {
			return models.Order{}, errors.New("Bad Request")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/abc", nil)
	w := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("id", "abc")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routerCtx))

	handler := New(slog.Default(), mock)
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetOrderByID_Error(t *testing.T) {
	mock := &OrdersMock{
		GetOrderByIDFunc: func(ctx context.Context, id int) (models.Order, error) {
			return models.Order{}, errors.New("DB error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/orders/1", nil)
	w := httptest.NewRecorder()

	routerCtx := chi.NewRouteContext()
	routerCtx.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routerCtx))

	handler := New(slog.Default(), mock)
	handler.GetOrderByID(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_Success(t *testing.T) {
	mock := &OrdersMock{
		GetOrdersByUserEmailFunc: func(ctx context.Context, email string) ([]models.Order, error) {
			return []models.Order{
				{CustomerID: 1, TotalPrice: 100},
			}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users/orders?email=example@mail.net", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_BadRequest(t *testing.T) {
	mock := &OrdersMock{
		GetOrdersByUserEmailFunc: func(ctx context.Context, email string) ([]models.Order, error) {
			return nil, errors.New("Bad request")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users/orders?gmail=example@mail.net", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetOrdersByUserEmail_Error(t *testing.T) {
	mock := &OrdersMock{
		GetOrdersByUserEmailFunc: func(ctx context.Context, email string) ([]models.Order, error) {
			return nil, errors.New("DB error")
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/users/orders?email=example@mail.net", nil)
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.GetOrdersByUserEmail(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestPlaceOrder_Success(t *testing.T) {
	mock := &OrdersMock{
		PlaceOrderFunc: func(ctx context.Context, userEmail string, items []models.OrderItem) (int, error) {
			return 1, nil
		},
	}

	body := `{
		"email":"example@email.net",
		"orderItems":[{"productID":1,"quantity":3}]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPlaceOrder_BadRequest(t *testing.T) {
	mock := &OrdersMock{
		PlaceOrderFunc: func(ctx context.Context, userEmail string, items []models.OrderItem) (int, error) {
			return 0, errors.New("Bad request")
		},
	}

	body := `{
		"email":"",
		"orderItems":[{"productID":1,"quantity":3}]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPlaceOrder_Error(t *testing.T) {
	mock := &OrdersMock{
		PlaceOrderFunc: func(ctx context.Context, userEmail string, items []models.OrderItem) (int, error) {
			return 0, errors.New("DB error")
		},
	}

	body := `{
		"email":"example@email.com",
		"orderItems":[{"productID":1,"quantity":3}]
	}`

	req := httptest.NewRequest(http.MethodPost, "/checkout", strings.NewReader(body))
	w := httptest.NewRecorder()

	handler := New(slog.Default(), mock)
	handler.PlaceOrder(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
