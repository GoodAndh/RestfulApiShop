package orders

import (
	"context"
	"database/sql"
	"restful/model/domain"
	"restful/model/web"
	"restful/utils"
)

type Service interface {
	CreateOrders(ctx context.Context, ord *web.OrdersCreatePayload) (int, error)
	CreateOrderItems(ctx context.Context, ord *web.OrderItemsCreatePayload) error
	GetOrders(ctx context.Context, userid int, status string) ([]web.Orders, error)
	GetOrderItems(ctx context.Context, userId int, OrderId ...int) ([]web.OrdersItems, error)
	GetOrdersById(ctx context.Context, userid int, status string, id int) ([]web.Orders, error)
	UpdateOrders(ctx context.Context, order *web.OrdersUpdatePayload) error
}

type ServiceImpl struct {
	db   *sql.DB
	repo Repository
}

func NewService(db *sql.DB, repo Repository) Service {
	return &ServiceImpl{
		db:   db,
		repo: repo,
	}
}

func (s *ServiceImpl) CreateOrders(ctx context.Context, ord *web.OrdersCreatePayload) (int, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer utils.CommitOrRollback(tx)

	orderId, err := s.repo.CreateOrders(ctx, tx, &domain.Orders{
		Total:  ord.Total,
		Status: ord.Status,
		Userid: ord.Userid,
	})
	if err != nil {
		return 0, err
	}

	return orderId, nil
}

func (s *ServiceImpl) CreateOrderItems(ctx context.Context, ord *web.OrderItemsCreatePayload) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	err = s.repo.CreateOrderItems(ctx, tx, &domain.OrdersItems{
		Total:      ord.Total,
		TotalPrice: ord.TotalPrice,
		OrderId:    ord.OrderId,
		ProductId:  ord.ProductId,
		UserId:     ord.Userid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) GetOrders(ctx context.Context, userid int, status string) ([]web.Orders, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	ord, err := s.repo.GetOrders(ctx, tx, userid, status)
	if err != nil {
		return nil, err
	}
	return utils.ConvertOrdersToSlice(ord), nil
}

func (s *ServiceImpl) GetOrderItems(ctx context.Context, userId int, OrderId ...int) ([]web.OrdersItems, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	ord, err := s.repo.GetOrderItems(ctx, tx, userId, OrderId...)
	if err != nil {
		return nil, err
	}
	return utils.ConvertOrderItemsToSlice(ord), nil
}

func (s *ServiceImpl) GetOrdersById(ctx context.Context, userid int, status string, id int) ([]web.Orders, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	ord, err := s.repo.GetOrdersById(ctx, tx, userid, status, id)
	if err != nil {
		return nil, err
	}
	return utils.ConvertOrdersToSlice(ord), nil
}

func (s *ServiceImpl) UpdateOrders(ctx context.Context, order *web.OrdersUpdatePayload) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	err = s.repo.UpdateOrders(ctx, tx, &domain.Orders{
		Id:     order.Id,
		Total:  order.Total,
		Status: order.Status,
		Userid: order.Userid,
	})
	if err != nil {
		return err
	}

	return nil
}
