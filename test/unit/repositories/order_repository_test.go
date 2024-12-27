package repositories_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	repositories "github.com/shayja/orders-service/internal/adapters/repositories/orders"
	"github.com/shayja/orders-service/internal/entities"
	"github.com/stretchr/testify/assert"
)

func TestGetAllOrders(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.OrderRepository{Db: db}

	userID := "451fa817-41f4-40cf-8dc2-c9f22aa98a4f"
	page := 1
	expectedOrders := []*entities.Order{
		{
			ID:         "6204037c-30e6-408b-8aaa-dd8219860b4b",
			UserID:     userID,
			TotalPrice: 100.0,
			Status:     1,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "total_price", "status", "created_at", "updated_at"}).
		AddRow(expectedOrders[0].ID, userID, expectedOrders[0].TotalPrice, expectedOrders[0].Status, expectedOrders[0].CreatedAt, expectedOrders[0].UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM get_user_orders\\(\\$1, \\$2, \\$3\\)").
		WithArgs(userID, 0, repositories.PAGE_SIZE).
		WillReturnRows(rows)

	orders, err := repo.GetAllOrders(page, userID)

	assert.NoError(t, err)
	assert.Len(t, orders, len(expectedOrders))
	assert.Equal(t, expectedOrders[0].ID, orders[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.OrderRepository{Db: db}

	orderID := "6204037c-30e6-408b-8aaa-dd8219860b4b"
	expectedOrder := &entities.Order{
		ID:         orderID,
		UserID:     "451fa817-41f4-40cf-8dc2-c9f22aa98a4f",
		TotalPrice: 150.0,
		Status:     2,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "total_price", "status", "created_at", "updated_at"}).
		AddRow(expectedOrder.ID, expectedOrder.UserID, expectedOrder.TotalPrice, expectedOrder.Status, expectedOrder.CreatedAt, expectedOrder.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM get_order\\(\\$1\\)").
		WithArgs(orderID).
		WillReturnRows(rows)

	order, err := repo.GetByID(orderID)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.ID, order.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.OrderRepository{Db: db}

	orderRequest := &entities.OrderRequest{
		UserID:      "451fa817-41f4-40cf-8dc2-c9f22aa98a4f",
		TotalPrice:  200.0,
		Status:      1,
		OrderDetails: []entities.OrderDetail{
			{ProductID: "063d0ff7-e17e-4957-8d92-a988caeda8a1", Quantity: 2, UnitPrice: 50.0, TotalPrice: 100.0},
		},
	}
	//newID := "6204037c-30e6-408b-8aaa-dd8219860b4b"

	mock.ExpectExec("CALL orders_insert\\(\\$1, \\$2, \\$3, \\$4::order_detail_type\\[\\], \\$5\\)").
		WithArgs(orderRequest.UserID, orderRequest.TotalPrice, orderRequest.Status, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := repo.Create(orderRequest)

	assert.NoError(t, err)
	assert.NotEmpty(t, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repositories.OrderRepository{Db: db}

	orderID := "6204037c-30e6-408b-8aaa-dd8219860b4b"
	newStatus := 3
	expectedOrder := &entities.Order{
		ID:         orderID,
		UserID:     "451fa817-41f4-40cf-8dc2-c9f22aa98a4f",
		TotalPrice: 150.0,
		Status:     newStatus,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mock.ExpectExec("CALL orders_update_status\\(\\$1, \\$2\\)").
		WithArgs(orderID, newStatus).
		WillReturnResult(sqlmock.NewResult(1, 1))

	rows := sqlmock.NewRows([]string{"id", "user_id", "total_price", "status", "created_at", "updated_at"}).
		AddRow(expectedOrder.ID, expectedOrder.UserID, expectedOrder.TotalPrice, expectedOrder.Status, expectedOrder.CreatedAt, expectedOrder.UpdatedAt)

	mock.ExpectQuery("SELECT \\* FROM get_order\\(\\$1\\)").
		WithArgs(orderID).
		WillReturnRows(rows)

	order, err := repo.UpdateStatus(orderID, newStatus)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrder.Status, order.Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}
