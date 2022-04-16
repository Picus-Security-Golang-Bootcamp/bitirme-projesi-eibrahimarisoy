package order

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"
	common "patika-ecommerce/pkg/utils"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo     MockOrderRepository
	orderItemRepo *OrderItemRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(orderRepo *OrderRepository, orderItemRepo *OrderItemRepository) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
	}
}

// CompleteOrder completes an order
func (r *OrderService) CompleteOrder(user *model.User, req *api.OrderRequest) (*model.Order, error) {
	cartId, err := common.StrfmtToUUID(*req.CartID)
	if err != nil {
		return nil, err
	}

	order, err := r.orderRepo.CompleteOrder(user, cartId)
	if err != nil {
		return nil, err
	}

	return order, nil
}

// GetOrdersByUser returns all orders of a user
func (r *OrderService) GetOrdersByUser(user *model.User, pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error) {
	return r.orderRepo.GetOrdersByUser(user, pagination)
}

// CancelOrder cancels an order
func (r *OrderService) CancelOrder(user *model.User, id uuid.UUID) error {
	return r.orderRepo.CancelOrder(id, user)
}
