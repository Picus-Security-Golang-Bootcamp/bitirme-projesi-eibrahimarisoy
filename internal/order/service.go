package order

import (
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/cart"
	"patika-ecommerce/internal/model"
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
	fmt.Println(*r.cartRepo)
	fmt.Println(user)
	fmt.Println(req)
	cart, err := r.cartRepo.GetCreatedCartByUserAndCart(user, *req.CartID)
	if err != nil {
		return nil, err
	}

	order, err := r.orderRepo.CompleteOrder(cart)
	if err != nil {
		return nil, err
	}

	// for _, item := range req.Items {
	// 	orderItem := &model.OrderItem{
	// 		OrderID:   order.ID,
	// 		ProductID: item.ProductID,
	// 		Quantity:  item.Quantity,
	// 	}

	// 	if err := r.orderItemRepo.db.Create(orderItem).Error; err != nil {
	// 		return nil, err
	// 	}
	// }

	return order, nil
}
