package order

import (
	"context"
	"database/sql"
	"study/internal/domain/order"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 插入订单
	query := `
		INSERT INTO "order" (order_no, user_id, status, total_amount, expire_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, query,
		o.OrderNo, o.UserID, o.Status, o.TotalAmount, o.ExpireAt, o.CreatedAt,
	).Scan(&o.ID)
	if err != nil {
		return err
	}

	// 插入订单项
	for i := range o.Items {
		itemQuery := `
			INSERT INTO order_item (order_id, product_id, quantity, unit_price)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		err = tx.QueryRowContext(ctx, itemQuery,
			o.ID, o.Items[i].ProductID, o.Items[i].Quantity, o.Items[i].UnitPrice,
		).Scan(&o.Items[i].ID)
		if err != nil {
			return err
		}
		o.Items[i].OrderID = o.ID
	}

	return tx.Commit()
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	query := `
		SELECT id, order_no, user_id, status, total_amount, paid_at, expire_at, created_at
		FROM "order"
		WHERE id = $1
	`

	var o order.Order
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&o.ID, &o.OrderNo, &o.UserID, &o.Status, &o.TotalAmount,
		&o.PaidAt, &o.ExpireAt, &o.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}

	// 获取订单项
	items, err := r.getOrderItems(ctx, id)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return &o, nil
}

func (r *OrderRepository) getOrderItems(ctx context.Context, orderID int64) ([]order.OrderItem, error) {
	query := `
		SELECT id, order_id, product_id, quantity, unit_price
		FROM order_item
		WHERE order_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []order.OrderItem
	for rows.Next() {
		var item order.OrderItem
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.UnitPrice,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, nil
}

func (r *OrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*order.Order, error) {
	query := `
		SELECT id, order_no, user_id, status, total_amount, paid_at, expire_at, created_at
		FROM "order"
		WHERE order_no = $1
	`

	var o order.Order
	err := r.db.QueryRowContext(ctx, query, orderNo).Scan(
		&o.ID, &o.OrderNo, &o.UserID, &o.Status, &o.TotalAmount,
		&o.PaidAt, &o.ExpireAt, &o.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, order.ErrOrderNotFound
		}
		return nil, err
	}

	// 获取订单项
	items, err := r.getOrderItems(ctx, o.ID)
	if err != nil {
		return nil, err
	}
	o.Items = items

	return &o, nil
}

func (r *OrderRepository) Update(ctx context.Context, o *order.Order) error {
	query := `
		UPDATE "order" 
		SET status = $1, paid_at = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query, o.Status, o.PaidAt, o.ID)
	return err
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID int64, offset, limit int) ([]*order.Order, error) {
	query := `
		SELECT id, order_no, user_id, status, total_amount, paid_at, expire_at, created_at
		FROM "order"
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*order.Order
	for rows.Next() {
		var o order.Order
		err := rows.Scan(
			&o.ID, &o.OrderNo, &o.UserID, &o.Status, &o.TotalAmount,
			&o.PaidAt, &o.ExpireAt, &o.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	return orders, nil
}
