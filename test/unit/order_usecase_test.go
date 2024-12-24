package usecases_test

import (
	"errors"
	"testing"

	"github.com/shayja/orders-service/internal/entities"
	"github.com/shayja/orders-service/internal/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func (m *OrderRepositoryMock) GetAllOrders(page int, userId string) ([]*entities.Order, error) {
	args := m.Called(page, userId)
	return args.Get(0).([]*entities.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetById(id string) (*entities.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Order), args.Error(1)
}

func (m *OrderRepositoryMock) Create(orderRequest *entities.OrderRequest) (string, error) {
	args := m.Called(orderRequest)
	return args.String(0), args.Error(1)
}

func (m *OrderRepositoryMock) UpdateStatus(id string, status int) (*entities.Order, error) {
	args := m.Called(id, status)
	return args.Get(0).(*entities.Order), args.Error(1)
}

func TestOrderUsecase_GetOrders(t *testing.T) {
	orderRepositoryMock := new(OrderRepositoryMock)
	orderUsecase := &usecases.OrderUsecase{OrderRepo: orderRepositoryMock}

	expectedOrders := []*entities.Order{}
	orderRepositoryMock.On("GetAllOrders", 1, "test-user-id").Return(expectedOrders, nil)

	orders, err := orderUsecase.GetOrders(1, "test-user-id")
	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
	orderRepositoryMock.AssertCalled(t, "GetAllOrders", 1, "test-user-id")
}

func TestOrderUsecase_GetOrders_Error(t *testing.T) {
	orderRepositoryMock := new(OrderRepositoryMock)
	orderUsecase := &usecases.OrderUsecase{OrderRepo: orderRepositoryMock}

	// Mock the behavior: return nil orders and an error
	orderRepositoryMock.On("GetAllOrders", 1, "test-user-id").Return(([]*entities.Order)(nil), errors.New("db error"))

	orders, err := orderUsecase.GetOrders(1, "test-user-id")
	assert.Error(t, err)           // Expecting an error
	assert.Nil(t, orders)          // Expecting orders to be nil
	orderRepositoryMock.AssertCalled(t, "GetAllOrders", 1, "test-user-id")
}
