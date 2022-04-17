package cart

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"patika-ecommerce/internal/api"
	httpErr "patika-ecommerce/internal/httpErrors"
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

func getUpdateCartItemPayload(quantity int) []byte {
	q := strconv.Itoa(quantity)
	var jsonStr = []byte(`{"quantity":` + q + `}`)

	return jsonStr
}

func Test_cartHandler_getOrCreateCart(t *testing.T) {
	id := uuid.New()
	user := model.User{
		Base: model.Base{ID: id},
	}
	t.Run("createCart_Success", func(t *testing.T) {
		mockService := &mockCartService{
			carts: []model.Cart{},
			users: []model.User{
				{
					Base: model.Base{ID: id},
				},
			},
		}
		cartHandler := &cartHandler{
			cartService: mockService,
		}
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		cartHandler.getOrCreateCart(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("getCart_Success", func(t *testing.T) {
		mockService := &mockCartService{
			carts: []model.Cart{
				{
					Base:   model.Base{ID: id},
					UserID: id,
					User:   user,
				},
			},
			users: []model.User{
				{
					Base: model.Base{ID: id},
				},
			},
		}
		cartHandler := &cartHandler{
			cartService: mockService,
		}
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		cartHandler.getOrCreateCart(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("getCart_Failed_userNotFound", func(t *testing.T) {
		newUser := model.User{
			Base: model.Base{ID: uuid.New()},
		}
		mockService := &mockCartService{
			carts: []model.Cart{
				{
					Base:   model.Base{ID: id},
					UserID: id,
					User:   user,
				},
			},
			users: []model.User{
				{
					Base: model.Base{ID: id},
				},
			},
		}
		cartHandler := &cartHandler{
			cartService: mockService,
		}
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &newUser)
		cartHandler.getOrCreateCart(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func Test_cartHandler_addToCart(t *testing.T) {
	cartId, productId, userId := uuid.New(), uuid.New(), uuid.New()
	user := model.User{
		Base: model.Base{ID: userId},
	}
	userTwo := model.User{
		Base: model.Base{ID: uuid.New()},
	}

	mockService := &mockCartService{
		carts: []model.Cart{
			{
				Base:   model.Base{ID: cartId},
				UserID: userId,
				Items:  []model.CartItem{},
				Status: model.CartStatusCreated,
			},
		},
		users: []model.User{user},
	}
	cartHandler := &cartHandler{
		cartService: mockService,
	}

	t.Run("addToCart_Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/add", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getAddToCartPayload(productId.String(), 2)))
		cartHandler.addToCart(c)

		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("addToCart_Fault_unAuthorizedUser", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/add", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &userTwo)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getAddToCartPayload(productId.String(), 2)))
		cartHandler.addToCart(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("addToCart_Fault_reqBodyNotValid", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/add", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("{getAddToCartPayload(productId.String(), 2)}")))
		cartHandler.addToCart(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("addToCart_Fault_reqBodyFormat", func(t *testing.T) {
		payloadNotValid := []byte(`{"productId":test, "quantity":2}`)
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/add", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(payloadNotValid))
		cartHandler.addToCart(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func Test_cartHandler_listCartItems(t *testing.T) {
	var productStock int64 = 10
	cartId, productId, userId := uuid.New(), uuid.New(), uuid.New()
	productName := "product name"

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
	cartHandler := &cartHandler{
		cartService: mockService,
	}
	t.Run("listCartItems_Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/items", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		cartHandler.listCartItems(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})
	newUser := model.User{Base: model.Base{ID: uuid.New()}}
	t.Run("listCartItems_Failed_userNotFound", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/items", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &newUser)
		cartHandler.listCartItems(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}

func Test_cartHandler_updateCartItem(t *testing.T) {
	cartId, cartItemId, productId, userId := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	productName := "product name"
	var productStock int64 = 10

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
						Base:      model.Base{ID: cartItemId},
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
	cartHandler := &cartHandler{
		cartService: mockService,
	}
	t.Run("updateCartItem_Success", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/items/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		c.Params = []gin.Param{{Key: "id", Value: cartItemId.String()}}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getUpdateCartItemPayload(1)))
		cartHandler.updateCartItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("updateCartItem_Failed_invalidReqBody", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/items/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		c.Params = []gin.Param{{Key: "id", Value: cartItemId.String()}}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("invalid payload")))
		cartHandler.updateCartItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("updateCartItem_Failed_cartItemNotFound", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/items/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		c.Params = []gin.Param{{Key: "id", Value: "cartItemId.String()"}}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getUpdateCartItemPayload(1)))
		cartHandler.updateCartItem(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("updateCartItem_Failed_notEnoughStock", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/cart/items/:id", nil)
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("user", &user)
		c.Params = []gin.Param{{Key: "id", Value: cartItemId.String()}}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(getUpdateCartItemPayload(5000)))
		cartHandler.updateCartItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func Test_cartHandler_deleteCartItem(t *testing.T) {
	cartId, cartItemId, productId, userId := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	productName := "product name"
	var productStock int64 = 10

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
						Base:      model.Base{ID: cartItemId},
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
	cartHandler := &cartHandler{
		cartService: mockService,
	}

	t.Run("deleteCartItem_Success", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Set("user", &user)
		c.Params = []gin.Param{{Key: "id", Value: cartItemId.String()}}
		cartHandler.deleteCartItem(c)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("deleteCartItem_Failed_cartItemNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Set("user", &user)
		c.Params = []gin.Param{{Key: "id", Value: uuid.New().String()}}
		cartHandler.deleteCartItem(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("deleteCartItem_Failed_idNotValid", func(t *testing.T) {
		w := httptest.NewRecorder()
		gin.SetMode(gin.TestMode)
		c, _ := gin.CreateTestContext(w)
		c.Set("user", &user)
		c.Params = []gin.Param{{Key: "id", Value: "uuid.New().String()"}}
		cartHandler.deleteCartItem(c)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

}

var (
	UserNotFoundError    = fmt.Errorf("user not found")
	ProductNotFoundError = fmt.Errorf("product not found")
	NotEnoughStockError  = fmt.Errorf("not enough stock")
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
	}
	return nil, UserNotFoundError

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

	}
	return nil, httpErr.UnauthorizedError
}

// UpdateCartItem updates a cart item
func (r *mockCartService) UpdateCartItem(user *model.User, id uuid.UUID, req *api.CartItemUpdateRequest) (*model.CartItem, error) {

	for _, item := range r.carts {
		if item.UserID == user.ID && item.Status == model.CartStatusCreated {

			for index, cartItem := range item.Items {
				if cartItem.ID == id {
					if req.Quantity > *cartItem.Product.Stock {
						return nil, fmt.Errorf("Product stock is not enough")
					}
					if req.Quantity == 0 {
						item.Items = append(item.Items[:index], item.Items[index+1:]...)
						return &cartItem, nil
					}
					cartItem.Quantity = req.Quantity
					return &cartItem, nil
				}
			}
		}

	}
	return nil, CartItemNotFoundError
}

// DeleteCartItem deletes a cart item
func (r *mockCartService) DeleteCartItem(user *model.User, id uuid.UUID) error {

	for _, item := range r.carts {
		if item.UserID == user.ID && item.Status == model.CartStatusCreated {
			for index, cartItem := range item.Items {
				if cartItem.ID == id {
					item.Items = append(item.Items[:index], item.Items[index+1:]...)
					return nil
				}
			}
		}
	}
	return CartItemNotFoundError

}
