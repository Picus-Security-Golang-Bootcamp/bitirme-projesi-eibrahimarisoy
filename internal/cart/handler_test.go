package cart

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func Test_cartHandler_getOrCreateCart(t *testing.T) {
	// name, description := "test", "test"
	id := uuid.New()

	gin.SetMode(gin.TestMode)

	user := model.User{
		Base: model.Base{ID: id},
	}

	mockService := &mockCartService{
		carts: []model.Cart{},
		users: []model.User{
			{
				Base: model.Base{ID: id},
			},
		},
	}
	w := httptest.NewRecorder()
	cartHandler := &cartHandler{
		cartService: mockService,
	}
	c, r := gin.CreateTestContext(w)

	r.POST("/cart", cartHandler.getOrCreateCart)
	c.Request, _ = http.NewRequest("POST", "/cart", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", &user)
	// c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getCategoryPOSTPayload()))
	cartHandler.getOrCreateCart(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

var (
	UserNotFoundError = fmt.Errorf("user not found")
)

type mockCartService struct {
	carts []model.Cart
	users []model.User
}

// GetOrCreateCart returns a cart by user id
func (r *mockCartService) GetOrCreateCart(user *model.User) (*model.Cart, error) {
	newUser := &model.User{}

	for _, item := range r.users {
		if item.ID == user.ID {
			newUser = &item
		}
	}

	if newUser == nil {
		return nil, UserNotFoundError
	}

	for _, item := range r.carts {
		if item.UserID == newUser.ID && item.Status == model.CartStatusCreated {
			return &item, nil
		}
	}

	cart := &model.Cart{
		Base: model.Base{
			ID: uuid.New(),
		},
		UserID: newUser.ID,
		Status: model.CartStatusCreated,
	}

	return cart, nil
}

// AddToCart adds a product to cart
func (r *mockCartService) AddToCart(user *model.User, req *api.AddToCartRequest) (*model.Cart, error) {
	cart := &model.Cart{}
	// cart, err := r.cartRepo.GetCreatedCart(user)
	// if err != nil {
	// 	return nil, err
	// }

	// pId, err := common.StrfmtToUUID(req.ProductID)
	// if err != nil {
	// 	return nil, err
	// }

	// // find product by given id
	// product, err := r.productRepo.GetProductWithoutCategories(pId)
	// if err != nil {
	// 	return nil, err
	// }

	// // check if product is available in cart
	// is_exists := false
	// for _, item := range cart.Items {
	// 	// if product is already in cart then update quantity
	// 	if item.ProductID == pId {
	// 		if int64(item.Quantity)+req.Quantity > *product.Stock {
	// 			return nil, fmt.Errorf("Product stock is not enough")
	// 		}
	// 		item.Quantity += req.Quantity
	// 		r.cartItemRepo.UpdateCartItem(&item)
	// 		is_exists = true
	// 		break
	// 	}
	// }

	// // if product not exists in cart, create new cart item
	// if !is_exists {
	// 	if *product.Stock < req.Quantity {
	// 		return nil, fmt.Errorf("Product stock is not enough")
	// 	}
	// 	if err := r.cartItemRepo.Create(cart, product); err != nil {
	// 		return nil, err
	// 	}
	// }
	return cart, nil
}

// UpdateCartItem updates a cart item
func (r *mockCartService) UpdateCartItem(user *model.User, id uuid.UUID, req *api.CartItemUpdateRequest) (*model.CartItem, error) {
	// cart, err := r.cartRepo.GetCreatedCartWithItemsAndProducts(user)
	// if err != nil {
	// 	return nil, err
	// }

	// cartItem, err := cart.GetCartItemByID(id)
	// if err != nil {
	// 	return nil, err
	// }

	// if req.Quantity == 0 {
	// 	r.cartItemRepo.DeleteCartItem(cartItem)
	// 	cartItem.Quantity = 0
	// 	return cartItem, nil
	// }

	// if req.Quantity > *cartItem.Product.Stock {
	// 	return nil, fmt.Errorf("product stock is not enough")
	// }

	// cartItem.Quantity = req.Quantity
	// if err := r.cartItemRepo.UpdateCartItem(cartItem); err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

// DeleteCartItem deletes a cart item
func (r *mockCartService) DeleteCartItem(user *model.User, id uuid.UUID) error {
	// cart, err := r.cartRepo.GetCreatedCartWithItems(user)
	// if err != nil {
	// 	return err
	// }

	// cartItem, err := cart.GetCartItemByID(id)
	// if err != nil {
	// 	return err
	// }

	// if err := r.cartItemRepo.DeleteCartItem(cartItem); err != nil {
	// 	return err
	// }

	return nil
}
