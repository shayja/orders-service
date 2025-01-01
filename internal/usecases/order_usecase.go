// usecases/order_usecase.go
package usecases

import (
	"github.com/shayja/orders-service/internal/entities"
)

type OrderRepository interface {
	GetAllOrders(page int, userID string) ([]*entities.Order, error)
	GetByID(id string) (*entities.Order, error)
	Create(orderRequest *entities.OrderRequest) (string, error)
	UpdateStatus(id string, status int) (*entities.Order, error)
}

type OrderUsecase struct {
	OrderRepo OrderRepository
}

func (uc *OrderUsecase) GetOrders(page int, userID string) ([]*entities.Order, error) {
	return uc.OrderRepo.GetAllOrders(page, userID)
}

func (uc *OrderUsecase) GetByID(id string) (*entities.Order, error) {
	return uc.OrderRepo.GetByID(id)
}

func (uc *OrderUsecase) Create(orderRequest *entities.OrderRequest) (string, error) {
	return uc.OrderRepo.Create(orderRequest)
}

func (uc *OrderUsecase) UpdateStatus(id string, status int) (*entities.Order, error) {
	return uc.OrderRepo.UpdateStatus(id, status)
}