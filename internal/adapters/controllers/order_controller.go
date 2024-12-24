package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shayja/orders-service/internal/entities"
	"github.com/shayja/orders-service/internal/usecases"
	"github.com/shayja/orders-service/pkg/utils"
)

type OrderController struct {
	OrderUsecase *usecases.OrderUsecase
}

// GetAll godoc
// @Summary      Get all orders
// @Description  Retrieve a paginated list of all orders
// @Tags         Orders
// @Param        page  query     int  true  "Page number"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Router       /order [get]
// @Security apiKey
func (oc *OrderController) GetOrders(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "Invalid page number"})
		return
	}

	// Get the userId from the token and not from the request, to ensure the user can only see their own orders
	userId, exists := c.Get("userId")
	if !exists {
		// Bad token - no userId. Stop here.
		c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "msg": "User ID not found in token"})
		return
	}

	// Validate the userId is a valid UUID
	if !utils.IsValidUUID(userId.(string)) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": "Invalid user id"})
		return
	}

	// Fetch the orders using the userId from the token
	res, err := oc.OrderUsecase.GetOrders(page, userId.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if res != nil {
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "msg": "No orders found for this page"})
	}
}

// GetById godoc
// @Summary      Get an order by ID
// @Description  Retrieve order details by order ID
// @Tags         Orders
// @Param        id   path      string  true  "Order ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /order/{id} [get]
// @Security apiKey
func (oc *OrderController) GetById(c *gin.Context) {

	var uri entities.IdRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	res, err := oc.OrderUsecase.GetById(uri.Id)
	if err != nil || !utils.IsValidUUID(res.Id) {
		c.JSON(http.StatusNotFound, gin.H{"status": "failed", "msg": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
}

// Create godoc
// @Summary      Create a new order
// @Description  Add a new order
// @Tags         Orders
// @Param        order  body      entities.OrderRequest  true  "Order data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /order [post]
// @Security apiKey
func (oc *OrderController) Create(c *gin.Context) {

	var post *entities.OrderRequest
	if err := c.ShouldBind(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	insertedId, err := oc.OrderUsecase.Create(post)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success", "id": insertedId, "msg": nil})
}

// UpdateStatus godoc
// @Summary      Update order status
// @Description  Update the status of an order
// @Tags         Orders
// @Param        id      path      string  true  "Order ID"
// @Param        status  body      int     true  "New status"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Router       /order/{id}/status [put]
// @Security apiKey
func (oc *OrderController) UpdateStatus(c *gin.Context) {

	var uri entities.IdRequest
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	var status struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "failed", "msg": err})
		return
	}

	res, err := oc.OrderUsecase.UpdateStatus(uri.Id, status.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "failed", "msg": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": res, "msg": nil})
}
