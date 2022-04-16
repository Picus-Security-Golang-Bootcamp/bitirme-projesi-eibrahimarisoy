package order

import (
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	paginationHelper "patika-ecommerce/pkg/pagination"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func TestOrderService_CompleteOrder(t *testing.T) {
	cartId := uuid.New()
	reqCartId := strfmt.UUID(cartId.String())

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

	type fields struct {
		orderRepo *mockOrderRepo
	}
	type args struct {
		user *model.User
		req  *api.OrderRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Order
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				orderRepo: orderRepo,
			},
			args: args{
				user: &model.User{
					Base: model.Base{ID: userId},
				},
				req: &api.OrderRequest{
					CartID: &reqCartId,
				},
			},
			want: &model.Order{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrderService{
				orderRepo: tt.fields.orderRepo,
			}
			_, err := r.CompleteOrder(tt.args.user, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderService.CompleteOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, int64(9), *orderRepo.products[0].Stock)
			assert.Equal(t, 1, len(orderRepo.orders))

		})
	}
}

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
	return nil, nil
}

// GetOrderByIdAndUser returns an order by id and user
func (r *mockOrderRepo) GetOrderByIdAndUser(user *model.User, id uuid.UUID) (*model.Order, error) {
	return nil, nil
}

// CancelOrder cancels an order
func (r *mockOrderRepo) CancelOrder(id uuid.UUID, user *model.User) error {
	return nil
}
