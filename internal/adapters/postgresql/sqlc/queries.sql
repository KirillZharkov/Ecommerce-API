-- name: ListProducts :many
SELECT * FROM products;

-- name: FindPoductsByID :one
SELECT * FROM products WHERE id = $1;
-- Входные данные: customer_id = 42
-- Выходные данные: {id: 1, customer_id: 42, created_at: '2024-01-15 10:30:00'}

-- name: CreateOrder :one
INSERT INTO orders (customer_id) VALUES ($1) RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price_cents) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateProductQuantity :execrows
UPDATE products
SET quantity = quantity - $1
WHERE id = $2 AND quantity >= $1;

-- name: FindOrdersByID :one
SELECT * FROM orders WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (id, name, price_in_cents, quantity) VALUES ($1, $2, $3, $4) RETURNING *;

