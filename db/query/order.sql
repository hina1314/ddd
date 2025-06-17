-- name: CreateOrder :one
INSERT INTO "order" (order_no, user_id, status, total_amount, paid_at)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM "order"
WHERE id = $1;

-- name: GetOrderByOrderNo :one
SELECT * FROM "order"
WHERE order_no = $1;

-- name: ListOrdersByUserID :many
SELECT * FROM "order"
WHERE user_id = $1
ORDER BY created_at DESC
    LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatusByOrderNo :exec
UPDATE "order"
SET status = $1
WHERE order_no = $2;

-- name: UpdatePaidAtAndStatusByOrderNo :exec
UPDATE "order"
SET paid_at = $1,
    status = $2
WHERE order_no = $3;

-- name: DeleteOrderByID :exec
DELETE FROM "order"
WHERE id = $1;
