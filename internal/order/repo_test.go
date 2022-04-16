package order

import (
	"patika-ecommerce/internal/model"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func TestOrderRepository_CompleteOrder(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		user   *model.User
		cartId uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Order
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &OrderRepository{
				db: tt.fields.db,
			}
			got, err := r.CompleteOrder(tt.args.user, tt.args.cartId)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderRepository.CompleteOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderRepository.CompleteOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
