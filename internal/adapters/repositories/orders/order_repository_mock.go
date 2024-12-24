package repositories

import (
	"github.com/shayja/orders-service/internal/entities"
	"github.com/stretchr/testify/mock"
)

var MockDatabase = make(map[string][]entities.Order)



type MockOrderRepository struct {
	mock.Mock
}

func NewMockOrderRepository() *MockOrderRepository {
	return &MockOrderRepository{}
}


func (m *MockOrderRepository) GetAllOrders(page int, userId string) ([]*entities.Order, error) {
	args := m.Called(page, userId)
	if orders, ok := args.Get(0).([]*entities.Order); ok {
		return orders, args.Error(1)
	}
	return nil, args.Error(1)
}
