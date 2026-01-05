package orders

import (
	"context"

	repo "github.com/KirillZharkov/Ecommerce-API/internal/adapters/postgresql/sqlc"
)

type orderItem struct {
	ProductId int64 `json:"productId"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	CustomerID int64       `json:"customer_id"`
	Items      []orderItem `json:"items"`
}


type Service interface {
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}
