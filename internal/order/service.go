package order

type OrderService struct {
	orderRepo     *OrderRepository
	orderItemRepo *OrderItemRepository
	// productRepo   *product.ProductRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(orderRepo *OrderRepository, orderItemRepo *OrderItemRepository) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		// productRepo:   productRepo,
	}
}
