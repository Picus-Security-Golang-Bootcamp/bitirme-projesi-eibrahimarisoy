// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AddToCartRequest add to cart request
//
// swagger:model AddToCartRequest
type AddToCartRequest struct {

	// product Id
	// Format: uuid
	ProductID strfmt.UUID `json:"productId,omitempty"`

	// quantity
	Quantity int64 `json:"quantity,omitempty"`
}

// Validate validates this add to cart request
func (m *AddToCartRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateProductID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AddToCartRequest) validateProductID(formats strfmt.Registry) error {
	if swag.IsZero(m.ProductID) { // not required
		return nil
	}

	if err := validate.FormatOf("productId", "body", "uuid", m.ProductID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this add to cart request based on context it is used
func (m *AddToCartRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AddToCartRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AddToCartRequest) UnmarshalBinary(b []byte) error {
	var res AddToCartRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
