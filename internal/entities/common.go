// internal/entities/common.go
package entities

// IDRequest represents a request to get any entity by its ID.
type IDRequest struct {
	ID string `uri:"id" binding:"required" example:"451fa817-41f4-40cf-8dc2-c9f22aa98a4f" minLength:"36"`
}