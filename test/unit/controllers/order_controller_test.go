package controllers

import (
	"testing"

	"github.com/shayja/orders-service/internal/entities"
	"github.com/shayja/orders-service/internal/usecases"
	"github.com/shayja/orders-service/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetOrders_ValidRequest(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	mockUsecase := &usecases.OrderUsecase{OrderRepo: mockRepo}

	// Mock data
	mockOrders := []*entities.Order{
		{ID: "1", UserID: "user123", Status: 2},
	}
	mockRepo.On("GetAllOrders", 1, "user123").Return(mockOrders, nil)

	// Call the usecase
	orders, err := mockUsecase.GetOrders(1, "user123")

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.Equal(t, 1, len(orders))
	assert.Equal(t, "1", orders[0].ID)

	// Ensure expectations
	mockRepo.AssertExpectations(t)
}