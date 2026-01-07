package products

import (
	"context"
	"time"

	repo "github.com/KirillZharkov/Ecommerce-API/internal/adapters/postgresql/sqlc"
)

type createProductsParams struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	PriceInCents uint32    `json:"price_in_cents"`
	Quantity     int32     `json:"quantity"`
	CreatedAt    time.Time `json:"created_at"`
}
type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindPoductsByID(ctx context.Context, id int64) (repo.Product, error)
	PlaceProduct(ctx context.Context, tempOrder repo.CreateProductParams) (repo.Product, error)
}
