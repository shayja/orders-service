// adapters/repositories/orders/order_repository.go
package repositories

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/shayja/orders-service/internal/entities"
	"github.com/shayja/orders-service/pkg/utils"
)

type OrderRepository struct {
	Db *sql.DB
}

const PAGE_SIZE = 20
// Get all user orders
func (r *OrderRepository) GetAllOrders(page int, userID string) ([]*entities.Order, error) {
	offset := PAGE_SIZE * (page - 1)
	query := `SELECT * FROM get_user_orders($1, $2, $3)`
	rows, err := r.Db.Query(query, userID, offset, PAGE_SIZE)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	var orders []*entities.Order
	for rows.Next() {
		order := &entities.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// Get order by ID
func (r *OrderRepository) GetByID(id string) (*entities.Order, error) {
	query := `SELECT * FROM get_order($1)`
	rows, err := r.Db.Query(query, id)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer rows.Close()

	order := &entities.Order{}
	if rows.Next() {
		err := rows.Scan(&order.ID, &order.UserID, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
	}
	return order, nil
}

// Create a new order
func (r *OrderRepository) Create(orderRequest *entities.OrderRequest) (string, error) {
	newID := utils.CreateNewUUID().String()
	_, err := r.Db.Exec(
		`CALL orders_insert($1, $2, $3, $4::order_detail_type[], $5)`,
		orderRequest.UserID,
		orderRequest.TotalPrice,
		orderRequest.Status,
		pq.Array(orderRequest.OrderDetails),
		&newID)
	if err != nil {
		fmt.Print(err)
		return "", err
	}

	fmt.Printf("Order %s created successfully\n", newID)
	return newID, nil
}

// Update order status
func (r *OrderRepository) UpdateStatus(id string, status int) (*entities.Order, error) {
	_, err := r.Db.Exec("CALL orders_update_status($1, $2)", id, status)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	return r.GetByID(id)
}