package order

import (
	"context"
	"fmt"
	"go-pet-shop/internal/models"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Order interface {
	CreateOrder(ctx context.Context, order models.Order) (int, error)
	AddOrderItem(ctx context.Context, orderItem models.OrderItem) error
	GetOrderByID(ctx context.Context, id int) (models.Order, error)
	GetOrdersByUserEmail(ctx context.Context, email string) ([]models.Order, error)
	GetOrderItemsByOrderID(ctx context.Context, orderID int) ([]models.OrderItem, error)
	PlaceOrder(ctx context.Context, userEmail string, items []models.OrderItem) (int, error)
}

type Handler struct {
	log     *slog.Logger
	storage Order
}

func New(log *slog.Logger, storage Order) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
	}
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.hanlers.order.CreateOrder"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Creating new order", slog.String("url", r.URL.String()))

	var order models.Order

	if err := render.DecodeJSON(r.Body, &order); err != nil {
		log.Error("failed to decode JSON body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Failed to decode JSON body",
		})
		return
	}

	order.CreatedAt = time.Now()

	orderId, err := h.storage.CreateOrder(r.Context(), order)
	if err != nil {
		log.Error("failed to create order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create order",
		})
		return
	}

	log.Info("Order created successfully",
		slog.Int("id", orderId),
		slog.Int("customer_id", order.CustomerID),
		slog.Time("created_ad", order.CreatedAt),
	)

	order.ID = orderId
	render.JSON(w, r, map[string]interface{}{
		"status": "Order created successfully",
		"order":  order,
	})
}

func (h *Handler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.hanlers.order.AddOrderItem"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Creating new order item", slog.String("url", r.URL.String()))

	strOrderID := chi.URLParam(r, "id")
	if strOrderID == "" {
		log.Error("empty order id")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order ID is required",
		})
		return
	}

	orderID, err := strconv.Atoi(strOrderID)
	if err != nil {
		log.Error("invalid id format", slog.Any("error", err), slog.String("id", strOrderID))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Product ID must be a number",
		})
		return
	}

	var order_item models.OrderItem

	order_item.OrderID = orderID

	if err := render.DecodeJSON(r.Body, &order_item); err != nil {
		log.Error("failed to decode JSON body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Failed to decode JSON body",
		})
		return
	}

	err = h.storage.AddOrderItem(r.Context(), order_item)
	if err != nil {
		log.Error("failed to create order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to create order",
		})
		return
	}

	log.Info("Order item created successfully",
		slog.Int("product_id", order_item.ProductID),
		slog.Int("order_id", order_item.OrderID),
	)

	render.JSON(w, r, map[string]interface{}{
		"status":    "Order created successfully",
		"orderID":   order_item.OrderID,
		"productID": order_item.ProductID,
		"quantity":  order_item.Quantity,
	})
}

func (h *Handler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.handlers.orders.GetOrderByID"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	idStr := chi.URLParam(r, "id")

	if idStr == "" {
		log.Error("empty id")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Order id is required",
		})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("invalid id", slog.Any("error", err), slog.String("id", idStr))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Invalid ID format",
		})
		return
	}

	var order models.Order

	if order, err = h.storage.GetOrderByID(r.Context(), id); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") ||
			strings.Contains(strings.ToLower(err.Error()), "no rows") ||
			strings.Contains(strings.ToLower(err.Error()), "rows affected: 0") {
			log.Warn("order not found", slog.Int("id", id))
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, map[string]interface{}{
				"error":   "Not found",
				"message": fmt.Sprintf("Order with ID %d does not exist", id),
				"id":      id,
			})
			return
		}

		log.Error("failed to get order by id", slog.Any("error", err), slog.Int("id", id))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to get order by id",
		})
		return
	}

	log.Info("Retrieved order by ID successfully",
		slog.Int("id", id),
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, order)
}

func (h *Handler) GetOrdersByUserEmail(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.handlers.orders.GetOrdersByUserEmail"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	email := r.URL.Query().Get("email")

	if email == "" {
		log.Error("empty email")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Email query is empty",
		})
		return
	}

	orders, err := h.storage.GetOrdersByUserEmail(r.Context(), email)

	if err != nil {
		log.Error("failed to get orders by email", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve orders by user email",
		})
		return
	}

	log.Info("Retrieved orders successfully",
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, orders)
}

func (h *Handler) GetOrderItemsByOrderID(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.handlers.orders.GetOrderItemsByOrderID"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	email := r.URL.Query().Get("email")

	if email == "" {
		log.Error("empty email")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Email query is empty",
		})
		return
	}

	orders, err := h.storage.GetOrdersByUserEmail(r.Context(), email)

	if err != nil {
		log.Error("failed to get orders by email", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to retrieve orders by user email",
		})
		return
	}

	log.Info("Retrieved orders successfully",
		slog.String("url", r.URL.String()),
	)

	render.JSON(w, r, orders)
}

func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	const fn = "internal.handlers.order.PlaceOrder"

	log := h.log.With(
		slog.String("fn", fn),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	log.Info("Placing order", slog.String("url", r.URL.String()))

	var order_item models.OrderItemsList

	err := render.DecodeJSON(r.Body, &order_item)
	if err != nil {
		log.Error("failed to decode JSON body", slog.Any("error", err))
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Failed to decode JSON body",
		})
		return
	}

	if order_item.Email == "" {
		log.Error("email is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error":   "Bad request",
			"message": "Email can't be empty",
		})
		return
	}

	id, err := h.storage.PlaceOrder(r.Context(), order_item.Email, order_item.OrderItems)
	if err != nil {
		log.Error("Failed to place order", slog.Any("error", err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error":   "Internal server error",
			"message": "Failed to place order",
		})
		return
	}

	log.Info("Order placed successfully", slog.Int("id", id))
	render.JSON(w, r, map[string]interface{}{
		"status":   "Order placed successfully",
		"order_id": id,
	})
}
