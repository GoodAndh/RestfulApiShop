package produk

import (
	"context"
	"database/sql"
	"fmt"
	"restful/exception"
	"restful/model/domain"
)

type Repository interface {
	GetAllProduk(ctx context.Context, tx *sql.Tx) ([]domain.Product, error)
	GetById(ctx context.Context, tx *sql.Tx, id int) (*domain.Product, error)
	GetByName(ctx context.Context, tx *sql.Tx, name string) ([]domain.Product, error)
	CreateProduct(ctx context.Context, tx *sql.Tx, product *domain.Product) error
	UpdateProduct(ctx context.Context, tx *sql.Tx, product *domain.Product) error
}
type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}

func (r *RepositoryImpl) GetAllProduk(ctx context.Context, tx *sql.Tx) ([]domain.Product, error) {
	rows, err := tx.QueryContext(ctx, "select * from product")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pr, err := rangeProduct(rows)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (r *RepositoryImpl) GetById(ctx context.Context, tx *sql.Tx, id int) (*domain.Product, error) {
	pr := &domain.Product{}
	err := tx.QueryRowContext(ctx, "select * from product where id = ?", id).Scan(&pr.Id, &pr.Name, &pr.Deskripsi, &pr.Category, &pr.Quantity, &pr.Price, &pr.Userid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, exception.ErrNotFound
		}
		return nil, err
	}

	return pr, nil

}

func (r *RepositoryImpl) GetByName(ctx context.Context, tx *sql.Tx, name string) ([]domain.Product, error) {
	query := "select * from product where name in( select name from product where name like ? ) or category in (select category from product where category like ?) "
	rows, err := tx.QueryContext(ctx, query, name+"%", name+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	pr, err := rangeProduct(rows)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (r *RepositoryImpl) CreateProduct(ctx context.Context, tx *sql.Tx, product *domain.Product) error {
	_, err := tx.ExecContext(ctx, "insert into product (name,deskripsi,category,quantity,price,userid) values(?,?,?,?,?,?)", product.Name, product.Deskripsi, product.Category, product.Quantity, product.Price, product.Userid)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepositoryImpl) UpdateProduct(ctx context.Context, tx *sql.Tx, product *domain.Product) error {
	result, err := tx.ExecContext(ctx, "update product set name = ? ,deskripsi = ? ,category = ? ,quantity = ? ,price = ? where id = ? ", product.Name, product.Deskripsi, product.Category, product.Quantity, product.Price, product.Id)
	if err != nil {
		return err
	}

	if affected, err := result.RowsAffected(); affected == 0 || err != nil {
		return fmt.Errorf("update product failed ,rows affected = %d and err = %v", affected, err)
	}

	return nil
}

func rangeProduct(rows *sql.Rows) ([]domain.Product, error) {
	p := []domain.Product{}

	for rows.Next() {
		pr := &domain.Product{}
		err := rows.Scan(&pr.Id, &pr.Name, &pr.Deskripsi, &pr.Category, &pr.Quantity, &pr.Price, &pr.Userid)
		if err != nil {
			return nil, err
		}
		p = append(p, *pr)
	}
	return p, nil

}
