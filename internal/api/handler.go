package api

import (
	"net/http"
	"x5_test/internal/domain"
	"x5_test/internal/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler struct {
	orderService *service.OrderService
}

func NewHandler(os *service.OrderService) *Handler {
	return &Handler{orderService: os}
}

func (h *Handler) Register(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", h.HealthCheck)
	e.POST("/orders", h.CreateOrder)
	e.GET("/orders/:id", h.GetOrder)
	e.GET("/orders", h.ListOrders)
}

type createOrderRequest struct {
	CustomerID string        `json:"customer_id"`
	Items      []domain.Item `json:"items"`
}

func (h *Handler) CreateOrder(c echo.Context) error {
	var req createOrderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.CustomerID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "customer_id is required"})
	}

	order, err := h.orderService.CreateOrder(c.Request().Context(), req.CustomerID, req.Items)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Ответ - 201 created с созданным заказом
	return c.JSON(http.StatusCreated, order)
}

func (h *Handler) GetOrder(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid order id"})
	}

	order, err := h.orderService.GetOrder(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if order == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "order not found"})
	}

	return c.JSON(http.StatusOK, order)
}

func (h *Handler) ListOrders(c echo.Context) error {
	customerID := c.QueryParam("customer_id")
	status := domain.OrderStatus(c.QueryParam("status"))

	orders, err := h.orderService.ListOrders(c.Request().Context(), customerID, status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if orders == nil {
		orders = []domain.Order{}
	}

	return c.JSON(http.StatusOK, orders)
}

func (h *Handler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
