package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/shayja/orders-service/internal/adapters/controllers"
	"github.com/shayja/orders-service/internal/entities"
	"github.com/shayja/orders-service/internal/usecases"
	"github.com/shayja/orders-service/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func setupRouter(orderController *controllers.OrderController) *gin.Engine {
	router := gin.Default()

	// Middleware for user ID (Mock user extraction)
	router.Use(func(c *gin.Context) {
		// Mock userId extraction from the token
		c.Set("userId", "123e4567-e89b-12d3-a456-426614174000")
		c.Next()
	})

	api := router.Group("/api/v1")
	{
		api.GET("/orders", orderController.GetOrders)
		api.GET("/order/:id", orderController.GetById)
		api.POST("/order", orderController.Create)
		api.PUT("/order/:id/status", orderController.UpdateStatus)
	}
	return router
}

func TestGetOrdersIntegration(t *testing.T) {

	// Mock Repository
	mockRepo := &MockOrderRepository{}

	// Usecase
	orderUsecase := &usecases.OrderUsecase{
		OrderRepo: mockRepo, // Fixed field name
	}

	// Controller
	orderController := &controllers.OrderController{
		OrderUsecase: orderUsecase,
	}

	// Setup Router
	router := setupRouter(orderController)

	// Mock Data
	mockOrder := entities.Order{
		Id:	"1",
		UserId:	"123e4567-e89b-12d3-a456-426614174000",
		TotalPrice:	100,
		Status:	1,// "Pending"
	}

	// Setup Mock Repository Data
	mockRepo.orders = []*entities.Order{&mockOrder}

	// Create Test Request
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/orders?page=1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check Response
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Status string	`json:"status"`
		Data   []*entities.Order `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.Equal(t, 1, len(response.Data))
	assert.Equal(t, mockOrder.Id, response.Data[0].Id)
}

func TestCreateOrderIntegration(t *testing.T) {
	// Mock Repository
	mockRepo := &MockOrderRepository{}

	// Usecase
	orderUsecase := &usecases.OrderUsecase{
		OrderRepo: mockRepo, // Fixed field name
	}

	// Controller
	orderController := &controllers.OrderController{
		OrderUsecase: orderUsecase,
	}

	// Setup Router
	router := setupRouter(orderController)

	// Mock Request Data
	orderRequest := &entities.OrderRequest{
		UserId:	"123e4567-e89b-12d3-a456-426614174000",
		TotalPrice: 150,
		Status: 1,//"Pending",
		OrderDetails: []entities.OrderDetail{
			{ProductId: "1", Quantity: 2},
		},
	}
	body, _ := json.Marshal(orderRequest)

	// Create Test Request
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/order", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check Response
	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Status string `json:"status"`
		Id     string `json:"id"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "success", response.Status)
	assert.NotEmpty(t, response.Id)

	entity, _ := mockRepo.GetById(response.Id)
	assert.NotNil(t, entity)
}

// Mock Repository
type MockOrderRepository struct {
	orders []*entities.Order
}

func (m *MockOrderRepository) GetAllOrders(page int, userId string) ([]*entities.Order, error) {
	var result []*entities.Order
	for _, order := range m.orders {
		if order.UserId == userId {
			result = append(result, order)
		}
	}
	return result, nil
}

func (m *MockOrderRepository) GetById(id string) (*entities.Order, error) {
	for _, order := range m.orders {
		if order.Id == id {
			return order, nil
		}
	}
	return nil, nil
}

func (m *MockOrderRepository) Create(orderRequest *entities.OrderRequest) (string, error) {
	newId := utils.CreateNewUUID().String()
	newOrder := &entities.Order{
		Id:         newId,
		UserId:     orderRequest.UserId,
		TotalPrice: orderRequest.TotalPrice,
		Status:     orderRequest.Status,
	}
	m.orders = append(m.orders, newOrder)
	return newId, nil
}

func (m *MockOrderRepository) UpdateStatus(id string, status int) (*entities.Order, error) {
	for _, order := range m.orders {
		if order.Id == id {
			order.Status = status
			return order, nil
		}
	}
	return nil, nil
}
