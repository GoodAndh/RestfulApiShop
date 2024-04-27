package users

import (
	"context"
	"database/sql"
	"restful/model/domain"
)

type Repository interface {
	GetByEmail(ctx context.Context, tx *sql.Tx, Email string) (*domain.Users, error)
	CreateUsers(ctx context.Context, tx *sql.Tx, user domain.Users) error
}
type RepositoryImpl struct{}

func NewRepository() Repository {
	return &RepositoryImpl{}
}
func (r *RepositoryImpl) GetByEmail(ctx context.Context, tx *sql.Tx, Email string) (*domain.Users, error) {
	u := &domain.Users{}
	err := tx.QueryRowContext(ctx, "select * from users where email = ?", Email).Scan(&u.Id, &u.Name, &u.Password, &u.Email)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *RepositoryImpl) CreateUsers(ctx context.Context, tx *sql.Tx, user domain.Users) error {
	_, err := tx.ExecContext(ctx, "insert into users(name,password,email) values(?,?,?)", user.Name, user.Password, user.Email)
	if err != nil {
		return err
	}
	return nil
}

//only use when need return slice of domain.users

// func rangeusers(rows *sql.Rows) (*domain.Users, error) {
// 	u := &domain.Users{}
// 	for rows.Next() {
// 		err := rows.Scan(&u.Id, &u.Name, &u.Password, &u.Email)
// 		if err != nil {
// 			if err == sql.ErrNoRows {
// 				return nil, err
// 			}
// 			return nil, err
// 		}
// 	}
// 	return u, nil
// }
