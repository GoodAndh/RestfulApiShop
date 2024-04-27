package users

import (
	"context"
	"database/sql"
	"fmt"
	"restful/layer/auth"
	"restful/model/domain"
	"restful/model/web"
	"restful/utils"
)

type Service interface {
	GetEmail(ctx context.Context, Email string) (*web.Users, error)
	CreateUsers(ctx context.Context, user web.UsersRegisterPayload) error
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

func (s *ServiceImpl) GetEmail(ctx context.Context, Email string) (*web.Users, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer utils.CommitOrRollback(tx)

	user, err := s.repo.GetByEmail(ctx, tx, Email)
	if err != nil {
		return nil, err
	}
	return utils.ConvertUserToSlice(user), nil
}

func (s *ServiceImpl) CreateUsers(ctx context.Context, user web.UsersRegisterPayload) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer utils.CommitOrRollback(tx)

	//check if email is already exist
	_, err = s.repo.GetByEmail(ctx, tx, user.Email)
	if err == nil {
		return fmt.Errorf("email already exist")
	}

	//hash the password
	password, err := auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	usr := &domain.Users{
		Name:     user.Name,
		Password: password,
		Email:    user.Email,
	}

	err = s.repo.CreateUsers(ctx, tx, *usr)
	if err != nil {
		return err
	}

	return nil

}
