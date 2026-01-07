package products

import (
	"log"
	"net/http"
	"strconv"

	repo "github.com/KirillZharkov/Ecommerce-API/internal/adapters/postgresql/sqlc"
	"github.com/KirillZharkov/Ecommerce-API/internal/json"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// обработчик списка твоаров
func (h *Handler) PlaceProduct(w http.ResponseWriter, r *http.Request) {
	var tempOrder createProductsParams
	if err := json.Read(r, &tempOrder); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createdOrder, err := h.service.PlaceProduct(r.Context(), repo.CreateProductParams{
		ID:           tempOrder.ID,
		Name:         tempOrder.Name,
		PriceInCents: int32(tempOrder.PriceInCents),
		Quantity:     tempOrder.Quantity})
	if err != nil {
		log.Println(err)
		if err == ErrProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.Write(w, http.StatusCreated, createdOrder)
}

// обработчик списка твоаров
func (h *Handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// сервис , который будет обращаться к ListProducts
	//1.service->ListProducts(repository)
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//2.возращаем json в Http-ответе
	// products := struct {
	// 	Products []string `json:"products"`
	// }{}
	json.Write(w, http.StatusOK, products)
}

func (h *Handler) FindPoductsByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}
	productsid, err := h.service.FindPoductsByID(r.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Write(w, http.StatusOK, productsid)
}
