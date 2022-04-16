package order

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func getOrderCompletePayload(cartID string) []byte {
	var jsonStr = []byte(`{"cartId": "` + cartID + `"}`)

	return jsonStr
}

func Test_orderHandler_completeOrder(t *testing.T) {
	cartId := uuid.New()
	// reqCartId := strfmt.UUID(cartId.String())

	productId := uuid.New()
	userId := uuid.New()
	// cartItemId := uuid.New()

	productName := "product name"
	productStock := int64(10)

	orderRepo := &mockOrderRepo{
		orders: []model.Order{},
		carts: []model.Cart{
			{
				Base: model.Base{ID: cartId},
				Items: []model.CartItem{
					{
						ProductID: productId,
						Quantity:  1,
						Price:     100,
					},
				},
				UserID: userId,
			},
		},
		products: []model.Product{
			{
				Base:  model.Base{ID: productId},
				Name:  &productName,
				Price: 100,
				Stock: &productStock,
			},
		},
	}
	gin.SetMode(gin.TestMode)

	user := model.User{
		Base: model.Base{ID: userId},
	}

	w := httptest.NewRecorder()
	orderHandler := &orderHandler{
		orderRepo: orderRepo,
	}
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &user)
	c.Request, _ = http.NewRequest("POST", "/orders", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getOrderCompletePayload(cartId.String())))

	orderHandler.completeOrder(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(9), *orderRepo.products[0].Stock)
	assert.Equal(t, 1, len(orderRepo.orders))

}

func Test_orderHandler_listOrders(t *testing.T) {
	cartId := uuid.New()
	productId := uuid.New()
	userId := uuid.New()

	productName := "product name"
	productStock := int64(10)

	orderRepo := &mockOrderRepo{
		orders: []model.Order{
			{
				CartID: cartId,
				Status: model.OrderStatusCompleted,
				Base:   model.Base{ID: uuid.New()},
				Items: []model.OrderItem{
					{
						ProductID: productId,
						Price:     100,
						Product: model.Product{
							Base:  model.Base{ID: productId},
							Name:  &productName,
							Price: 100,
							Stock: &productStock,
						},
					},
				},
				UserID:     userId,
				TotalPrice: 100,
			},
		},
	}
	pagination := paginationHelper.Pagination{
		Limit: 2,
		Page:  1,
	}

	gin.SetMode(gin.TestMode)

	user := model.User{
		Base: model.Base{ID: userId},
	}

	w := httptest.NewRecorder()
	orderHandler := &orderHandler{
		orderRepo: orderRepo,
	}
	c, _ := gin.CreateTestContext(w)
	c.Set("pagination", &pagination)
	c.Set("user", &user)
	c.Request, _ = http.NewRequest("GET", "/orders", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	orderHandler.listOrders(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_orderHandler_cancelOrder(t *testing.T) {
	cartId := uuid.New()
	productId := uuid.New()
	userId := uuid.New()
	orderId := uuid.New()

	productName := "product name"
	productStock := int64(10)

	orderRepo := &mockOrderRepo{
		orders: []model.Order{
			{
				Base:   model.Base{ID: orderId},
				CartID: cartId,
				Status: model.OrderStatusCompleted,
				Items: []model.OrderItem{
					{
						ProductID: productId,
						Price:     100,
						Product: model.Product{
							Base:  model.Base{ID: productId},
							Name:  &productName,
							Price: 100,
							Stock: &productStock,
						},
					},
				},
				UserID:     userId,
				TotalPrice: 100,
			},
		},
	}

	gin.SetMode(gin.TestMode)

	user := model.User{
		Base: model.Base{ID: userId},
	}

	w := httptest.NewRecorder()
	orderHandler := &orderHandler{
		orderRepo: orderRepo,
	}
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &user)
	c.Params = []gin.Param{{Key: "id", Value: orderId.String()}}
	c.Request, _ = http.NewRequest("DELETE", "/orders/:id", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	orderHandler.cancelOrder(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, int64(11), *orderRepo.orders[0].Items[0].Product.Stock)
}

var (
	OrderCannotBeCanceledError = fmt.Errorf("order cannot be canceled")
	OrderNotFoundError         = fmt.Errorf("order not found")
)

type mockOrderRepo struct {
	orders     []model.Order
	carts      []model.Cart
	orderItems []model.OrderItem
	products   []model.Product
}

func (r *mockOrderRepo) CompleteOrder(user *model.User, cartId uuid.UUID) (*model.Order, error) {
	for _, item := range r.carts {
		if item.ID == cartId && item.UserID == user.ID {

			order := model.Order{
				Base:       model.Base{ID: uuid.New()},
				UserID:     user.ID,
				CartID:     item.ID,
				TotalPrice: item.GetTotalPrice(),
			}
			r.orders = append(r.orders, order)

			for _, cartItem := range item.Items {
				for _, product := range r.products {
					if product.ID == cartItem.ProductID {
						*product.Stock -= cartItem.Quantity
					}
				}

				for i := 0; i < int(cartItem.Quantity); i++ {
					orderItem := model.OrderItem{
						Base:      model.Base{ID: uuid.New()},
						OrderID:   order.ID,
						ProductID: cartItem.ProductID,
						Price:     cartItem.Price,
					}
					r.orderItems = append(r.orderItems, orderItem)
				}
			}
			return &order, nil
		}
	}

	return nil, nil
}

// GetOrdersByUser returns all orders of a user
func (r *mockOrderRepo) GetOrdersByUser(user *model.User, pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error) {
	orders := []*model.Order{}
	for _, order := range r.orders {
		orders = append(orders, &order)
	}
	pagination.TotalRows = int64(len(orders))
	pagination.Rows = OrdersToOrderDetailedResponse(orders)
	return pagination, nil
}

// GetOrderByIdAndUser returns an order by id and user
func (r *mockOrderRepo) GetOrderByIdAndUser(user *model.User, id uuid.UUID) (*model.Order, error) {
	return nil, nil
}

// CancelOrder cancels an order
func (r *mockOrderRepo) CancelOrder(id uuid.UUID, user *model.User) error {
	for _, order := range r.orders {
		if order.ID == id && order.UserID == user.ID && order.Status == model.OrderStatusCompleted {
			if order.IsCancelable() {
				order.Status = model.OrderStatusCanceled

				for _, item := range order.Items {
					*item.Product.Stock += 1
				}

				return nil
			}
			return OrderCannotBeCanceledError
		}
	}
	return OrderNotFoundError
}
