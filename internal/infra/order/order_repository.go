package order

import (
	"context"
	"study/db/model"
	"study/internal/domain/order"
)

type OrderRepository struct {
	db model.TxManager
}

func NewOrderRepository(db model.TxManager) order.Repository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, o *order.Order) error {
	panic("implement me")
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*order.Order, error) {
	panic("implement me")
}

func (r *OrderRepository) getOrderItems(ctx context.Context, orderID int64) ([]order.OrderItem, error) {
	panic("implement me")
}

func (r *OrderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*order.Order, error) {
	panic("implement me")
}

func (r *OrderRepository) Update(ctx context.Context, o *order.Order) error {
	panic("implement me")
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID int64, offset, limit int) ([]*order.Order, error) {
	panic("implement me")
}

func (r *OrderRepository) AddCart(ctx context.Context, o *order.Cart) error {
	q := r.db.Querier(ctx)

	arg := model.AddCartParams{
		UserID:   o.UserID,
		SkuID:    o.SkuID,
		Quantity: o.Quantity,
		Price:    o.Price,
	}

	cart, err := q.AddCart(ctx, arg)
	if err != nil {
		// Todo: if duplicate, increase quantity
		//if infra.IsDuplicateKeyError(err) {
		//	return errors.New(errors.ErrUserAlreadyExists, "User already exists")
		//}
		return err
	}

	o.ID = cart.ID

	return nil
}
