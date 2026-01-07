package orders

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

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	//проверка полезной нагрузки
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one item is required")
	}

	tx, err := s.db.Begin(ctx) //запускает транзакцию
	if err != nil {
		return repo.Order{}, err
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)
	//создание заказа
	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, err
	}
	//найдем заказ, если он существует
	for _, item := range tempOrder.Items {
		product, err := qtx.FindPoductsByID(ctx, item.ProductId)
		if err != nil {
			return repo.Order{}, ErrProductNotFound
		}
		// if product.Quantity < item.Quantity {
		// 	return repo.Order{}, ErrProductNoStock
		// }
		//Обновите количество товара на складе
		rows, err := qtx.UpdateProductQuantity(ctx, repo.UpdateProductQuantityParams{
			ID:       item.ProductId,
			Quantity: item.Quantity, // количество для вычитания
		})
		if err != nil {
			return repo.Order{}, fmt.Errorf("failed to update product quantity: %w", err)
		}
		if rows == 0 {
			return repo.Order{}, ErrProductNoStock
		}
		//пишем заказ в бд
		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:    order.ID,
			ProductID:  item.ProductId,
			Quantity:   item.Quantity,
			PriceCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, err
	}
	return order, nil
}

func (s *svc) FindOrderByID(ctx context.Context, id int64) (repo.Order, error) {
	return s.repo.FindOrdersByID(ctx, id)
}
