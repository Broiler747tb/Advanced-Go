package http

import (
	"fmt"
	"net/http"
	"order-api/Product"
	"order-api/req"
	"order-api/res"
	"strconv"

	"gorm.io/gorm"
)

type ProductHandlerDeps struct {
}

type ProductHandler struct {
	ProductRepository Product.ProductRepository
}

func NewLinkHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{}
	router.HandleFunc("POST /link", handler.Create())
	router.HandleFunc("PATCH /link/{id}", handler.Update())
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.GoTo())
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ProductRepository = ProductRepo.ProductRepository
		body, err := req.HandleBody[Product.LinkCreateRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		product := Product.NewProduct(body.Name)
		for {
			existingProduct, err := handler.ProductRepository.GetByHash(product.Hash)
			if err != nil {

			}
			if existingProduct == nil {
				break
			}
			product.GenerateHash()
			fmt.Println(product.Hash)
		}
		createdProduct, err := handler.ProductRepository.Create(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res.Json(w, createdProduct, 201)
	}
}

func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ProductRepository = ProductRepo.ProductRepository
		body, err := req.HandleBody[Product.LinkUpdateRequest](w, r)
		if err != nil {
			return
		}
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		product, err := handler.ProductRepository.Update(&Product.Product{
			Name:        body.Name,
			Hash:        body.Hash,
			Model:       gorm.Model{ID: uint(id)},
			Description: body.Description,
			ImageLinks:  body.ImageLinks,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotModified)
		}
		res.Json(w, product, 201)
	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ProductRepository = ProductRepo.ProductRepository
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		existingProduct, err := handler.ProductRepository.GetById(uint(id))
		if existingProduct == nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = handler.ProductRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		res.Json(w, nil, 200)
	}
}

func (handler *ProductHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ProductRepository = ProductRepo.ProductRepository
		hash := r.PathValue("hash")
		foundProduct, err := handler.ProductRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		res.Json(w, foundProduct, 200)
	}
}
