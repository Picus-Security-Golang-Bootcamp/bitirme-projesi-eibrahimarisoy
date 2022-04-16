package cart

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func getAddToCartPayload(productID string, quantity int) []byte {
	q := strconv.Itoa(quantity)
	var jsonStr = []byte(
		`{"productId": "` + productID + `","quantity": ` + q + `}`)

	return jsonStr
}

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
	cartHandler.getOrCreateCart(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func Test_cartHandler_AddToCart(t *testing.T) {
	// name, description := "test", "test"
	cartId := uuid.New()
	productId := uuid.New()
	userId := uuid.New()

	gin.SetMode(gin.TestMode)

	user := model.User{
		Base: model.Base{ID: userId},
	}

	mockService := &mockCartService{
		carts: []model.Cart{
			{
				Base:   model.Base{ID: cartId},
				UserID: userId,
				Items:  []model.CartItem{},
			},
		},
		users: []model.User{user},
	}
	w := httptest.NewRecorder()
	cartHandler := &cartHandler{
		cartService: mockService,
	}
	c, r := gin.CreateTestContext(w)

	r.POST("/cart/add", cartHandler.addToCart)
	c.Request, _ = http.NewRequest("POST", "/cart/add", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", &user)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getAddToCartPayload(productId.String(), 2)))
	cartHandler.addToCart(c)

	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body)
}

func Test_cartHandler_listCartItems(t *testing.T) {
	// name, description := "test", "test"
	cartId := uuid.New()
	productId := uuid.New()
	userId := uuid.New()

	productName := "product name"
	var productStock int64 = 10
	gin.SetMode(gin.TestMode)

	user := model.User{
		Base: model.Base{ID: userId},
	}

	mockService := &mockCartService{
		carts: []model.Cart{
			{
				Base:   model.Base{ID: cartId},
				UserID: userId,
				Status: model.CartStatusCreated,
				Items: []model.CartItem{
					{
						ProductID: productId,
						Quantity:  2,
						Price:     100,
						Product: model.Product{
							Base:        model.Base{ID: productId},
							Name:        &productName,
							Description: "description",
							Price:       100,
							Stock:       &productStock,
						},
					},
				},
			},
		},
		users: []model.User{user},
	}
	w := httptest.NewRecorder()
	cartHandler := &cartHandler{
		cartService: mockService,
	}
	c, r := gin.CreateTestContext(w)

	r.POST("/cart/items", cartHandler.listCartItems)
	c.Request, _ = http.NewRequest("POST", "/cart/items", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("user", &user)
	cartHandler.listCartItems(c)

	assert.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body)
}

var (
	UserNotFoundError    = fmt.Errorf("user not found")
	ProductNotFoundError = fmt.Errorf("product not found")
)

type mockCartService struct {
	carts    []model.Cart
	users    []model.User
	products []model.Product
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
	product := &model.Product{}
	reqProductId, _ := uuid.Parse(req.ProductID.String())

	for _, item := range r.products {
		if item.ID == reqProductId {
			product = &item
		}
	}
	if product == nil {
		return nil, ProductNotFoundError
	}
	for _, item := range r.carts {
		if item.UserID == user.ID && item.Status == model.CartStatusCreated {
			cart = &item
		}
	}

	if cart == nil {
		return nil, CartNotFoundError
	}

	for _, item := range cart.Items {
		if item.ProductID == reqProductId {
			item.Quantity += req.Quantity
			return cart, nil
		}
	}

	cart.Items = append(cart.Items, model.CartItem{
		ProductID: reqProductId,
		Quantity:  req.Quantity,
	})

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
