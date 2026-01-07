package products

import (
	"context"
	"errors"
	"fmt"

	repo "github.com/KirillZharkov/Ecommerce-API/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product has not enough stock")
)

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}
func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindPoductsByID(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindPoductsByID(ctx, id)
}

func (s *svc) PlaceProduct(ctx context.Context, tempOrder repo.CreateProductParams) (repo.Product, error) {
	//проверка полезной нагрузки
	if tempOrder.ID == 0 {
		return repo.Product{}, fmt.Errorf("ID is required")
	}
	if tempOrder.Name == "" {
		return repo.Product{}, fmt.Errorf("Name is required")
	}
	// if len(tempOrder.Items) == 0 {
	// 	return repo.Order{}, fmt.Errorf("at least one item is required")
	// }

	tx, err := s.db.Begin(ctx) //запускает транзакцию
	if err != nil {
		return repo.Product{}, err
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)
	//создание продукта
	products, err := qtx.CreateProduct(ctx, repo.CreateProductParams{
		ID:           tempOrder.ID,
		Name:         tempOrder.Name,
		PriceInCents: int32(tempOrder.PriceInCents),
		Quantity:     tempOrder.Quantity,
		//CreatedAt:    tempOrder.CreatedAt,
	})
	if err != nil {
		return repo.Product{}, err
	}
	// //найдем заказ, если он существует
	// for _, item := range tempOrder {
	// 	product, err := qtx.FindPoductsByID(ctx, item.ProductId)
	// 	if err != nil {
	// 		return repo.Product{}, ErrProductNotFound
	// 	}
	// 	if product.Quantity < item.Quantity {
	// 		return repo.Product{}, ErrProductNoStock
	// 	}
	// 	//Обновите количество товара на складе
	// 	rows, err := qtx.UpdateProductQuantity(ctx, repo.UpdateProductQuantityParams{
	// 		ID:       item.ProductId,
	// 		Quantity: item.Quantity, // количество для вычитания
	// 	})
	// 	if err != nil {
	// 		return repo.Product{}, fmt.Errorf("failed to update product quantity: %w", err)
	// 	}
	// 	if rows == 0 {
	// 		return repo.Product{}, ErrProductNoStock
	// 	}
	// 	//пишем заказ в бд
	// 	_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{})
	// 	if err != nil {
	// 		return repo.Product{}, err
	// 	}
	// }
	if err := tx.Commit(ctx); err != nil {
		return repo.Product{}, err
	}
	return products, nil
}
