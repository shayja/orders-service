package entities

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// Order represents an order entity
type Order struct {
	Id         string    `json:"id"`
	UserId     string    `json:"user_id"`
	TotalPrice float64   `json:"total_price"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// OrderDetail represents details of an order
type OrderDetail struct {
	Id         string  `json:"id"`
	OrderId    string  `json:"order_id"`
	ProductId  string  `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderRequest struct {
	UserId       string            `json:"user_id"`
	TotalPrice   float64           `json:"total_price"`
	Status       int               `json:"status"`
	OrderDetails []OrderDetail `json:"order_details"`
}

// Convert order details to database-compatible array
func (v OrderDetail) Value() (driver.Value, error) {
    return []byte(fmt.Sprintf("(%s,%d,%f)", v.ProductId, v.Quantity, v.UnitPrice)), nil
}