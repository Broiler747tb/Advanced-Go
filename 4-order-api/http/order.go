package http

import (
	"errors"
	"net/http"
	"strconv"

	"order-api/Order"
	"order-api/Product"
	"order-api/User"
	"order-api/middleware"
	"order-api/req"
	"order-api/res"
)

type OrderHandlerDeps struct {
	OrderRepository   *Order.OrderRepository
	ProductRepository *Product.ProductRepository
	UserRepository    *User.UserRepository
	Secret            string
}

type OrderHandler struct {
	OrderRepository   *Order.OrderRepository
	ProductRepository *Product.ProductRepository
	UserRepository    *User.UserRepository
}

func NewOrderHandler(router *http.ServeMux, deps OrderHandlerDeps) {
	handler := &OrderHandler{
		OrderRepository:   deps.OrderRepository,
		ProductRepository: deps.ProductRepository,
		UserRepository:    deps.UserRepository,
	}
	router.Handle("POST /order", middleware.IsAuthed(handler.Create(), deps.Secret))
	router.Handle("GET /order/{id}", middleware.IsAuthed(handler.GetById(), deps.Secret))
	router.Handle("GET /my-orders", middleware.IsAuthed(handler.GetMyOrders(), deps.Secret))
}

func (handler *OrderHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[Order.OrderCreateRequest](w, r)
		if err != nil {
			return
		}
		user, err := handler.currentUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		products, err := handler.ProductRepository.GetByIds(body.ProductIds)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(products) != len(body.ProductIds) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		createdOrder, err := handler.OrderRepository.Create(&Order.Order{
			UserID:   user.ID,
			Products: products,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, createdOrder, http.StatusCreated)
	}
}

func (handler *OrderHandler) GetById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseUint(r.PathValue("id"), 10, 32)
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		user, err := handler.currentUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		order, err := handler.OrderRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if order.UserID != user.ID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		res.Json(w, order, http.StatusOK)
	}
}

func (handler *OrderHandler) GetMyOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := handler.currentUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		orders, err := handler.OrderRepository.GetByUserId(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(w, orders, http.StatusOK)
	}
}

func (handler *OrderHandler) currentUser(r *http.Request) (*User.User, error) {
	phone, ok := middleware.PhoneFromContext(r.Context())
	if !ok {
		return nil, errors.New("unauthorized")
	}
	return handler.UserRepository.FindByPhone(phone)
}
