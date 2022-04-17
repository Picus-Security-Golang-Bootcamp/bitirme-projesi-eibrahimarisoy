package cart

import (
	"errors"
	"fmt"
	"patika-ecommerce/internal/api"
	"patika-ecommerce/internal/model"
	product "patika-ecommerce/internal/product"
	paginationHelper "patika-ecommerce/pkg/pagination"
	"reflect"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
)

func TestCartService_GetOrCreateCart(t *testing.T) {
	userId := uuid.New()
	type fields struct {
		cartRepo     MockCartRepository
		cartItemRepo MockCartItemRepository
		productRepo  product.MockProductRepository
	}
	type args struct {
		user *model.User
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Cart
		wantErr bool
	}{
		{
			name: "createdCartSuccess",
			fields: fields{
				cartRepo: &mockCartRepo{
					items: []model.Cart{},
					users: []model.User{
						{
							Base: model.Base{ID: userId},
						},
					},
				},
				cartItemRepo: &mockCartItemRepo{items: []model.CartItem{}},
				// productRepo:  &mockProductRepo{items: []model.Product{}},
			},
			args: args{
				user: &model.User{
					Base: model.Base{ID: userId},
				},
			},
			want: &model.Cart{
				User: model.User{
					Base: model.Base{ID: userId},
				},
				Status: model.CartStatusCreated,
			},
			wantErr: false,
		},
		{
			name: "gettingCartSuccess",
			fields: fields{
				cartRepo: &mockCartRepo{
					items: []model.Cart{
						{
							Base: model.Base{ID: uuid.New()},
							User: model.User{
								Base: model.Base{ID: userId},
							},
							Status: model.CartStatusCreated,
						},
					},
					users: []model.User{
						{
							Base: model.Base{ID: userId},
						},
					},
				},
				cartItemRepo: &mockCartItemRepo{items: []model.CartItem{}},
				// productRepo:  &mockProductRepo{items: []model.Product{}},
			},
			args: args{
				user: &model.User{
					Base: model.Base{ID: userId},
				},
			},
			want: &model.Cart{
				User: model.User{
					Base: model.Base{ID: userId},
				},
				Status: model.CartStatusCreated,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CartService{
				cartRepo:     tt.fields.cartRepo,
				cartItemRepo: tt.fields.cartItemRepo,
				// productRepo:  tt.fields.productRepo,
			}
			got, err := r.GetOrCreateCart(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.GetOrCreateCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CartService.GetOrCreateCart() = %v, want %v", got, tt.want)
			}

			assert.Equal(t, tt.want.User.ID, got.User.ID)
			assert.Equal(t, model.CartStatusCreated, got.Status)
		})
	}
}

var (
	productOneID          = uuid.New()
	productTwoID          = uuid.New()
	productOneStock int64 = 10

	productOneName        = "productOne"
	productTwoName        = "productTwo"
	productTwoStock int64 = 10

	productOne = model.Product{
		Base:  model.Base{ID: productOneID},
		Name:  &productOneName,
		Stock: &productOneStock,
		Price: 10,
	}

	productTwo = model.Product{
		Base:  model.Base{ID: productTwoID},
		Name:  &productTwoName,
		Stock: &productTwoStock,
		Price: 20,
	}

	products = []model.Product{
		productOne,
		productTwo,
	}
)

func TestCartService_AddToCart(t *testing.T) {
	userId := uuid.New()
	type fields struct {
		cartRepo     MockCartRepository
		cartItemRepo MockCartItemRepository
		productRepo  product.MockProductRepository
	}
	type args struct {
		user *model.User
		req  *api.AddToCartRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Cart
		wantErr bool
	}{
		{
			name: "addToCartSuccess",
			fields: fields{
				cartRepo: &mockCartRepo{
					items: []model.Cart{
						{
							Base:   model.Base{ID: uuid.New()},
							UserID: userId,
							Status: model.CartStatusCreated,
						},
					},
				},
				cartItemRepo: &mockCartItemRepo{items: []model.CartItem{}},
				productRepo:  &mockProductRepo{items: products},
			},
			args: args{
				user: &model.User{
					Base: model.Base{ID: userId},
				},
				req: &api.AddToCartRequest{
					ProductID: strfmt.UUID(productOneID.String()),
					Quantity:  1,
				},
			},
			want:    &model.Cart{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CartService{
				cartRepo:     tt.fields.cartRepo,
				cartItemRepo: tt.fields.cartItemRepo,
				productRepo:  tt.fields.productRepo,
			}
			got, err := r.AddToCart(tt.args.user, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.AddToCart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("CartService.AddToCart() = %v, want %v", got, tt.want)
			// }
			assert.Equal(t, 1, len(got.Items))
		})
	}
}

func TestCartService_UpdateCartItem(t *testing.T) {
	cartID := uuid.New()
	cartItemID := uuid.New()
	userID := uuid.New()
	productID := uuid.New()

	type fields struct {
		cartRepo     MockCartRepository
		cartItemRepo MockCartItemRepository
		productRepo  product.MockProductRepository
	}
	type args struct {
		user *model.User
		id   uuid.UUID
		req  *api.CartItemUpdateRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.CartItem
		wantErr bool
	}{
		{
			name: "updateCartItemSuccess",
			fields: fields{
				cartRepo: &mockCartRepo{
					items: []model.Cart{
						{
							Base:   model.Base{ID: cartID},
							UserID: userID,
							Status: model.CartStatusCreated,
							Items: []model.CartItem{
								{
									Base:   model.Base{ID: cartItemID},
									CartID: cartID,
									Product: model.Product{
										Base:  model.Base{ID: productID},
										Name:  &productOneName,
										Stock: &productOneStock,
										Price: 10,
									},
									ProductID: productID,
									Quantity:  1,
								},
							},
						},
					},
				},
				cartItemRepo: &mockCartItemRepo{
					items: []model.CartItem{
						{
							Base:      model.Base{ID: cartItemID},
							CartID:    cartID,
							ProductID: productID,
							Quantity:  1,
						},
					},
				},
				productRepo: &mockProductRepo{
					items: []model.Product{
						{
							Base:  model.Base{ID: productID},
							Name:  &productOneName,
							Stock: &productOneStock,
							Price: 10,
						},
					},
				},
			},
			args: args{
				user: &model.User{
					Base: model.Base{ID: userID},
				},
				id: cartItemID,
				req: &api.CartItemUpdateRequest{
					Quantity: 3,
				},
			},
			want:    &model.CartItem{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CartService{
				cartRepo:     tt.fields.cartRepo,
				cartItemRepo: tt.fields.cartItemRepo,
				productRepo:  tt.fields.productRepo,
			}
			got, err := r.UpdateCartItem(tt.args.user, tt.args.id, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("CartService.UpdateCartItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, int64(3), got.Quantity)

			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("CartService.UpdateCartItem() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestCartService_DeleteCartItem(t *testing.T) {
	cartID := uuid.New()
	cartItemID := uuid.New()
	userID := uuid.New()
	productID := uuid.New()

	type fields struct {
		cartRepo     MockCartRepository
		cartItemRepo MockCartItemRepository
		productRepo  product.MockProductRepository
	}
	type args struct {
		user *model.User
		id   uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "deleteCartItemSuccess",
			fields: fields{
				cartRepo: &mockCartRepo{
					items: []model.Cart{
						{
							Base:   model.Base{ID: cartID},
							UserID: userID,
							Status: model.CartStatusCreated,
							Items: []model.CartItem{
								{
									Base:   model.Base{ID: cartItemID},
									CartID: cartID,
									Product: model.Product{
										Base:  model.Base{ID: productID},
										Name:  &productOneName,
										Stock: &productOneStock,
										Price: 10,
									},
									ProductID: productID,
									Quantity:  1,
								},
							},
						},
					},
				},
				cartItemRepo: &mockCartItemRepo{
					items: []model.CartItem{
						{
							Base:      model.Base{ID: cartItemID},
							CartID:    cartID,
							ProductID: productID,
							Quantity:  1,
						},
					},
				},
				productRepo: &mockProductRepo{
					items: []model.Product{
						{
							Base:  model.Base{ID: productID},
							Name:  &productOneName,
							Stock: &productOneStock,
							Price: 10,
						},
					},
				},
			},
			args: args{
				user: &model.User{
					Base: model.Base{ID: userID},
				},
				id: cartItemID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &CartService{
				cartRepo:     tt.fields.cartRepo,
				cartItemRepo: tt.fields.cartItemRepo,
				productRepo:  tt.fields.productRepo,
			}
			if err := r.DeleteCartItem(tt.args.user, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("CartService.DeleteCartItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

type mockCartRepo struct {
	items []model.Cart
	users []model.User
}

type mockCartItemRepo struct {
	items []model.CartItem
}

type mockProductRepo struct {
	items      []model.Product
	categories []model.Category
}

var (
	CartNotFoundError     = fmt.Errorf("cart not found")
	CartItemNotFoundError = fmt.Errorf("cart item not found")
)

// GetOrCreateCart if cart is exists returns it otherwise create cart and return it
func (r *mockCartRepo) GetOrCreateCart(user *model.User) (*model.Cart, error) {
	cart := &model.Cart{}

	for _, item := range r.users {
		if user.ID == item.ID {
			for _, item := range r.items {
				if item.UserID == user.ID {
					return &item, nil
				}
			}
			cart.User = *user
			cart.Status = model.CartStatusCreated
			r.items = append(r.items, *cart)
			return cart, nil
		}
	}

	return nil, errors.New("user not found")
}

// GetCreatedCart returns a created cart by user id
func (r *mockCartRepo) GetCreatedCart(user *model.User) (*model.Cart, error) {
	for _, item := range r.items {
		if user.ID == item.UserID && item.Status == model.CartStatusCreated {
			return &item, nil
		}
	}
	return nil, CartNotFoundError
}

// GetCreatedCartWithItemsAndProducts returns a cart by user id
func (r *mockCartRepo) GetCreatedCartWithItemsAndProducts(user *model.User) (*model.Cart, error) {
	for _, item := range r.items {
		if user.ID == item.UserID && item.Status == model.CartStatusCreated {
			return &item, nil
		}
	}
	return nil, CartNotFoundError
}

// GetCreatedCartWithItems returns a cart by user id
func (r *mockCartRepo) GetCreatedCartWithItems(user *model.User) (*model.Cart, error) {
	for _, item := range r.items {
		if user.ID == item.UserID && item.Status == model.CartStatusCreated {
			return &item, nil
		}
	}
	return nil, CartNotFoundError
}

// GetCreatedCartByUserAndCart returns a cart by user id
func (r *mockCartRepo) GetCreatedCartByUserAndCart(user *model.User, cartId strfmt.UUID) (*model.Cart, error) {
	cart := model.Cart{}
	// if err := r.db.Preload("Items.Product").
	// 	Where("user_id = ? AND status = ? AND id = ?", user.ID, model.CartStatusCreated, cartId).
	// 	First(&cart).Error; err != nil {
	// 	return nil, err
	// }
	// fmt.Println("cart", cart)
	return &cart, nil
}

// GetCartByID returns a cart by id
func (r *mockCartRepo) GetCartByID(id uuid.UUID) (*model.Cart, error) {
	for _, item := range r.items {
		if item.ID == id {
			return &item, nil
		}
	}

	return nil, errors.New("cart not found2")
}

// UpdateCart updates a cart
func (r *mockCartRepo) UpdateCart(cart *model.Cart) error {

	// return r.db.Model(&cart).Updates(cart).Error
	return nil
}

// ###### CART ITEM ######

func (r *mockCartItemRepo) Create(cart *model.Cart, product *model.Product) error {
	cartItem := model.CartItem{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  1,
		Price:     product.Price,
	}

	cart.Items = append(cart.Items, cartItem)
	r.items = append(r.items, cartItem)
	return nil
}

// UpdateCartItem updates a cart item
func (r *mockCartItemRepo) UpdateCartItem(cartItem *model.CartItem) error {
	for i, item := range r.items {
		if item.ID == cartItem.ID {
			r.items[i] = *cartItem
			return nil
		}
	}

	return CartItemNotFoundError
}

// GetCartItemByID returns a cart item by id
func (r *mockCartItemRepo) GetCartItemByCartAndIDWithProduct(cart *model.Cart, id uuid.UUID) (*model.CartItem, error) {
	cartItem := &model.CartItem{}
	// if err := r.db.Model(&cartItem).Preload("Product").Where("cart_id = ? AND id = ?", cart.ID, id).First(cartItem).Error; err != nil {
	// 	if err == gorm.ErrRecordNotFound {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }

	return cartItem, nil
}

// DeleteCartItem deletes a cart item
func (r *mockCartItemRepo) DeleteCartItem(cartItem *model.CartItem) error {
	for i, item := range r.items {
		if item.ID == cartItem.ID {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}

	// return r.db.Delete(cartItem).Error
	return CartItemNotFoundError
}

// Product
func (r *mockProductRepo) Insert(product *model.Product) error {

	// newCategories := []model.Category{}

	// for _, item := range product.Categories {
	// 	for _, category := range r.categories {
	// 		if item.ID == category.ID {
	// 			newCategories = append(newCategories, category)
	// 		}
	// 	}
	// }

	// if len(newCategories) != len(product.Categories) {
	// 	return errors.New("category not found")
	// }

	// product.Categories = newCategories

	// r.items = append(r.items, *product)

	return nil
}

// GetProducts get all products
func (r *mockProductRepo) GetAll(pagination *paginationHelper.Pagination) (*paginationHelper.Pagination, error) {
	// var products []model.Product
	// pagination.TotalRows = int64(len(r.items))
	// pagination.Rows = ProductsToResponse(&products)

	return pagination, nil
}

// GetProduct get a single product
func (r *mockProductRepo) Get(id uuid.UUID) (*model.Product, error) {
	// for _, item := range r.items {
	// 	if item.ID == id {
	// 		return &item, nil
	// 	}
	// }
	return nil, nil
}

// GetProductWithoutCategories get a single product
func (r *mockProductRepo) GetProductWithoutCategories(id uuid.UUID) (*model.Product, error) {
	for _, item := range r.items {
		if item.ID == id {
			return &item, nil
		}
	}
	return nil, errors.New("product not found")
}

// DeleteProduct delete a single product
func (r *mockProductRepo) Delete(product *model.Product) error {
	for i, item := range r.items {
		if item.ID == product.ID {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return errors.New("product not found")
}

// UpdateProduct update a single product
func (r *mockProductRepo) Update(product *model.Product) error {
	for i, item := range r.items {
		if item.ID == product.ID {
			r.items[i] = *product
			return nil
		}
	}
	return errors.New("product not found")
}
