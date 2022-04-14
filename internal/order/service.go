package order

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/cart"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo     *OrderRepository
	orderItemRepo *OrderItemRepository
	cartRepo      *cart.CartRepository
	// productRepo   *product.ProductRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(orderRepo *OrderRepository, orderItemRepo *OrderItemRepository, cartRepo *cart.CartRepository) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		cartRepo:      cartRepo,
		// productRepo:   productRepo,
	}
}

// CompleteOrder completes an order
func (r *OrderService) CompleteOrder(user *model.User, req *api.OrderRequest) (*model.Order, error) {
	// Check given cart is valid
	cart, err := r.cartRepo.GetCreatedCartByUserAndCart(user, *req.CartID)
	if err != nil {
		return nil, err
	}

	order, err := r.orderRepo.CompleteOrder(cart)
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
	order, err := r.orderRepo.GetOrderByIdAndUser(user, id)
	if err != nil {
		return err
	}

	return r.orderRepo.CancelOrder(order)
}
