package model

import (
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestCart_GetCartItemByID(t *testing.T) {
	cartId := uuid.New()
	cartItemId := uuid.New()
	productId := uuid.New()

	type fields struct {
		Base   Base
		Status CartStatus
		UserID uuid.UUID
		User   User
		Items  []CartItem
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *CartItem
		wantErr bool
	}{
		{
			name: "getCartItemByID_Succeed",
			fields: fields{
				Items: []CartItem{
					{
						Base:      Base{ID: cartItemId},
						CartID:    cartId,
						ProductID: productId,
						Quantity:  1,
						Price:     100,
					},
				},
			},
			args: args{
				id: cartItemId,
			},
			want: &CartItem{
				Base:      Base{ID: cartItemId},
				CartID:    cartId,
				ProductID: productId,
				Quantity:  1,
				Price:     100,
			},
			wantErr: false,
		},
		{
			name: "getCartItemByID_Failed",
			fields: fields{
				Items: []CartItem{
					{
						Base:      Base{ID: cartItemId},
						CartID:    cartId,
						ProductID: productId,
						Quantity:  1,
						Price:     100,
					},
				},
			},
			args: args{
				id: uuid.New(),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cart{
				Base:   tt.fields.Base,
				Status: tt.fields.Status,
				UserID: tt.fields.UserID,
				User:   tt.fields.User,
				Items:  tt.fields.Items,
			}
			got, err := c.GetCartItemByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cart.GetCartItemByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cart.GetCartItemByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCart_GetTotalPrice(t *testing.T) {
	type fields struct {
		Base   Base
		Status CartStatus
		UserID uuid.UUID
		User   User
		Items  []CartItem
	}
	tests := []struct {
		name   string
		fields fields
		want   float64
	}{
		{
			name: "getTotalPrice_Succeed",
			fields: fields{
				Items: []CartItem{
					{
						Base:      Base{ID: uuid.New()},
						CartID:    uuid.New(),
						ProductID: uuid.New(),
						Quantity:  1,
						Price:     100,
					},
					{
						Base:      Base{ID: uuid.New()},
						CartID:    uuid.New(),
						ProductID: uuid.New(),
						Quantity:  1,
						Price:     100,
					},
				},
			},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cart{
				Base:   tt.fields.Base,
				Status: tt.fields.Status,
				UserID: tt.fields.UserID,
				User:   tt.fields.User,
				Items:  tt.fields.Items,
			}
			if got := c.GetTotalPrice(); got != tt.want {
				t.Errorf("Cart.GetTotalPrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCart_BeforeCreate(t *testing.T) {
	type fields struct {
		Base   Base
		Status CartStatus
		UserID uuid.UUID
		User   User
		Items  []CartItem
	}
	type args struct {
		tx *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "BeforeCreate_Succeed",
			fields: fields{
				Base: Base{
					ID: uuid.New(),
				},
				UserID: uuid.New(),
				User:   User{},
				Items:  []CartItem{},
			},
			args: args{
				tx: nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Cart{
				Base:   tt.fields.Base,
				Status: tt.fields.Status,
				UserID: tt.fields.UserID,
				User:   tt.fields.User,
				Items:  tt.fields.Items,
			}
			if err := c.BeforeCreate(tt.args.tx); (err != nil) != tt.wantErr {
				t.Errorf("Cart.BeforeCreate() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, c.Status, CartStatusCreated)
		})
	}
}
