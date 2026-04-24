package user

import (
	"context"
	"go-pet-shop/internal/models"
)

type UsersMock struct {
	CreateUserFunc          func(ctx context.Context, user models.Customer) error
	GetUserByEmailFunc      func(ctx context.Context, email string) (models.Customer, error)
	GetAllUsersFunc         func(ctx context.Context) ([]models.Customer, error)
	GetUserOrderHistoryFunc func(ctx context.Context, email string) ([]models.OrderDetail, error)
}

func (m *UsersMock) CreateUser(ctx context.Context, user models.Customer) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	return nil
}

func (m *UsersMock) GetUserByEmail(ctx context.Context, email string) (models.Customer, error) {
	if m.GetUserByEmailFunc != nil {
		return m.GetUserByEmailFunc(ctx, email)
	}
	return models.Customer{}, nil
}

func (m *UsersMock) GetAllUsers(ctx context.Context) ([]models.Customer, error) {
	if m.GetAllUsersFunc != nil {
		return m.GetAllUsersFunc(ctx)
	}
	return nil, nil
}

func (m *UsersMock) GetUserOrderHistory(ctx context.Context, email string) ([]models.OrderDetail, error) {
	if m.GetUserOrderHistoryFunc != nil {
		return m.GetUserOrderHistoryFunc(ctx, email)
	}
	return nil, nil
}
