package orders

import (
	"context"
	"database/sql"
	"fmt"
	"restful/exception"
	"restful/model/domain"
	"strings"
)

type Repository interface {
	CreateOrders(ctx context.Context, tx *sql.Tx, ord *domain.Orders) (int, error)
	CreateOrderItems(ctx context.Context, tx *sql.Tx, ord *domain.OrdersItems) error
	GetOrders(ctx context.Context, tx *sql.Tx, userid int, status string) ([]domain.Orders, error)
	GetOrderItems(ctx context.Context, tx *sql.Tx, userId int, OrderId ...int) ([]domain.OrdersItems, error)
	GetOrdersById(ctx context.Context, tx *sql.Tx, userid int, status string, id int) ([]domain.Orders, error)
	UpdateOrders(ctx context.Context, tx *sql.Tx, order *domain.Orders) error
}

type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}

func (r *RepositoryImpl) CreateOrders(ctx context.Context, tx *sql.Tx, ord *domain.Orders) (int, error) {
	result, err := tx.ExecContext(ctx, "insert into orders(total,status,userid) values(?,?,?)", ord.Total, ord.Status, ord.Userid)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return 0, err
	}
	return int(id), nil

}

func (r *RepositoryImpl) CreateOrderItems(ctx context.Context, tx *sql.Tx, ord *domain.OrdersItems) error {
	result, err := tx.ExecContext(ctx, "insert into orders_items (total,total_price,order_id,product_id,users_id) values(?,?,?,?,?)", ord.Total, ord.TotalPrice, ord.OrderId, ord.ProductId, ord.UserId)
	if err != nil {
		return err
	}
	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return fmt.Errorf("rows affected equal to %v , error equal to %v", affected, err)
	}
	return nil
}

func (r *RepositoryImpl) GetOrders(ctx context.Context, tx *sql.Tx, userid int, status string) ([]domain.Orders, error) {

	rows, err := tx.QueryContext(ctx, "select * from orders where userid = ? or status = ? ", userid, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []domain.Orders{}
	for rows.Next() {
		ord := &domain.Orders{}
		err := rows.Scan(&ord.Id, &ord.Total, &ord.Status, &ord.Userid)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *ord)
	}

	if len(orders) <= 0 {
		return nil, exception.ErrNoOrderRow
	}

	return orders, nil
}

func (r *RepositoryImpl) GetOrderItems(ctx context.Context, tx *sql.Tx, userId int, OrderId ...int) ([]domain.OrdersItems, error) {
	placeholders := strings.Repeat("?,", len(OrderId))

	//turn orderId to slice of interface{}
	args := make([]any, len(OrderId))
	for i, v := range OrderId {
		args[i] = v
	}

	query := fmt.Sprintf("select * from orders_items where users_id = %d or order_id in (%s%d)", userId, placeholders, 0)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	orders := []domain.OrdersItems{}
	defer rows.Close()
	for rows.Next() {
		o := &domain.OrdersItems{}
		err := rows.Scan(&o.Id, &o.Total, &o.TotalPrice, &o.OrderId, &o.ProductId, &o.UserId)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *o)
	}
	return orders, nil

}

func (r *RepositoryImpl) GetOrdersById(ctx context.Context, tx *sql.Tx, userid int, status string, id int) ([]domain.Orders, error) {
	rows, err := tx.QueryContext(ctx, "select * from orders where userid = ? and id = ? or status = ? ", userid, id, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := []domain.Orders{}
	for rows.Next() {
		ord := &domain.Orders{}
		err := rows.Scan(&ord.Id, &ord.Total, &ord.Status, &ord.Userid)
		if err != nil {
			return nil, err
		}
		orders = append(orders, *ord)
	}

	if len(orders) <= 0 {
		return nil, exception.ErrNoOrderRow
	}

	return orders, nil
}

func (r *RepositoryImpl) UpdateOrders(ctx context.Context, tx *sql.Tx, order *domain.Orders) error {
	result, err := tx.ExecContext(ctx, "update orders set total = ? ,status = ? where id = ? and userid = ?", order.Total, order.Status, order.Id,order.Userid)
	if err != nil {
		return err
	}

	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return fmt.Errorf("update product failed ,rows affected = %d and err = %v", affected, err)
	}

	return nil
}
