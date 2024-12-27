package mocks

import (
	"github.com/shayja/orders-service/internal/entities"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

// Mock implementation for GetAllOrders
func (m *MockOrderRepository) GetAllOrders(page int, useID string) ([]*entities.Order, error) {
	args := m.Called(page, useID)
	return args.Get(0).([]*entities.Order), args.Error(1)
}

// Mock implementation for GetByID
func (m *MockOrderRepository) GetByID(id string) (*entities.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Order), args.Error(1)
}

// Mock implementation for Create
func (m *MockOrderRepository) Create(orderRequest *entities.OrderRequest) (string, error) {
	args := m.Called(orderRequest)
	return args.String(0), args.Error(1)
}

// Mock implementation for UpdateStatus
func (m *MockOrderRepository) UpdateStatus(id string, status int) (*entities.Order, error) {
	args := m.Called(id, status)
	return args.Get(0).(*entities.Order), args.Error(1)
}
