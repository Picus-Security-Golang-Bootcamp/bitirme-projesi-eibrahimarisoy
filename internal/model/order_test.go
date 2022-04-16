package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestOrder_IsCancelable(t *testing.T) {
	type fields struct {
		Base       Base
		UserID     uuid.UUID
		User       User
		Status     OrderStatus
		CartID     uuid.UUID
		Cart       Cart
		TotalPrice float64
		Items      []OrderItem
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "IsCancelable_Succeed",
			fields: fields{
				Base: Base{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
				},
				UserID:     uuid.New(),
				User:       User{Base: Base{ID: uuid.New()}},
				Status:     OrderStatusCompleted,
				CartID:     uuid.New(),
				Cart:       Cart{Base: Base{ID: uuid.New()}},
				TotalPrice: 100,
				Items: []OrderItem{
					{
						Base:      Base{ID: uuid.New()},
						OrderID:   uuid.New(),
						ProductID: uuid.New(),
						Price:     100,
					},
				},
			},
			want: true,
		},
		{
			name: "IsCancelable_Failed_OrderStatusNotCompleted",
			fields: fields{
				Base: Base{
					ID:        uuid.New(),
					CreatedAt: time.Now(),
				},
				UserID:     uuid.New(),
				User:       User{Base: Base{ID: uuid.New()}},
				Status:     OrderStatusCanceled,
				CartID:     uuid.New(),
				Cart:       Cart{Base: Base{ID: uuid.New()}},
				TotalPrice: 100,
				Items: []OrderItem{
					{
						Base:      Base{ID: uuid.New()},
						OrderID:   uuid.New(),
						ProductID: uuid.New(),
						Price:     100,
					},
				},
			},
			want: false,
		},
		{
			name: "IsCancelable_Failed_TimeLimit",
			fields: fields{
				Base: Base{
					ID:        uuid.New(),
					CreatedAt: time.Now().AddDate(0, 0, -15),
				},
				UserID:     uuid.New(),
				User:       User{Base: Base{ID: uuid.New()}},
				Status:     OrderStatusCanceled,
				CartID:     uuid.New(),
				Cart:       Cart{Base: Base{ID: uuid.New()}},
				TotalPrice: 100,
				Items: []OrderItem{
					{
						Base:      Base{ID: uuid.New()},
						OrderID:   uuid.New(),
						ProductID: uuid.New(),
						Price:     100,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Order{
				Base:       tt.fields.Base,
				UserID:     tt.fields.UserID,
				User:       tt.fields.User,
				Status:     tt.fields.Status,
				CartID:     tt.fields.CartID,
				Cart:       tt.fields.Cart,
				TotalPrice: tt.fields.TotalPrice,
				Items:      tt.fields.Items,
			}
			if got := o.IsCancelable(); got != tt.want {
				t.Errorf("Order.IsCancelable() = %v, want %v", got, tt.want)
			}
		})
	}
}
