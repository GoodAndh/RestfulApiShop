package produk

import (
	"context"
	"database/sql"
	"restful/model/domain"
	"restful/model/web"
	"restful/utils"
)

type Service interface {
	GetAllProduk(ctx context.Context) ([]web.Product, error)
	GetById(ctx context.Context, id int) (*web.Product, error)
	GetByName(ctx context.Context, name string) ([]web.Product, error)
	CreateProduct(ctx context.Context, product *web.ProductCreatePayload) error
	UpdateProduct(ctx context.Context, product *web.ProductUpdatePayload) error
}
type ServiceImpl struct {
	repo Repository
	db   *sql.DB
}

func NewService(repo Repository, db *sql.DB) Service {
	return &ServiceImpl{
		repo: repo,
		db:   db,
	}
}

func (s *ServiceImpl) GetAllProduk(ctx context.Context) ([]web.Product, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	pr, err := s.repo.GetAllProduk(ctx, tx)
	if err != nil {
		return nil, err
	}

	return utils.ConvertProductToSlice(pr), nil
}

func (s *ServiceImpl) GetById(ctx context.Context, id int) (*web.Product, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	pr, err := s.repo.GetById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return utils.ConvertProductToWeb(pr), nil
}

func (s *ServiceImpl) GetByName(ctx context.Context, name string) ([]web.Product, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	pr, err := s.repo.GetByName(ctx, tx, name)
	if err != nil {
		return nil, err
	}
	return utils.ConvertProductToSlice(pr), nil

}

func (s *ServiceImpl) CreateProduct(ctx context.Context, product *web.ProductCreatePayload) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)
	err = s.repo.CreateProduct(ctx, tx, &domain.Product{
		Name:      product.Name,
		Deskripsi: product.Deskripsi,
		Category:  product.Category,
		Quantity:  product.Quantity,
		Price:     product.Price,
		Userid:    product.Userid,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) UpdateProduct(ctx context.Context, product *web.ProductUpdatePayload) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	err = s.repo.UpdateProduct(ctx, tx, &domain.Product{
		Id:        product.Id,
		Name:      product.Name,
		Deskripsi: product.Deskripsi,
		Category:  product.Category,
		Quantity:  product.Quantity,
		Price:     product.Price,
	})
	if err != nil {
		return err
	}
	return nil
}
