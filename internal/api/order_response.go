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

// OrderResponse order response
//
// swagger:model OrderResponse
type OrderResponse struct {

	// cart Id
	CartID *CartResponse `json:"cartId,omitempty"`

	// id
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// status
	Status string `json:"status,omitempty"`
}

// Validate validates this order response
func (m *OrderResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCartID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OrderResponse) validateCartID(formats strfmt.Registry) error {
	if swag.IsZero(m.CartID) { // not required
		return nil
	}

	if m.CartID != nil {
		if err := m.CartID.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cartId")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("cartId")
			}
			return err
		}
	}

	return nil
}

func (m *OrderResponse) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validate this order response based on the context it is used
func (m *OrderResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCartID(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OrderResponse) contextValidateCartID(ctx context.Context, formats strfmt.Registry) error {

	if m.CartID != nil {
		if err := m.CartID.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("cartId")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("cartId")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *OrderResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OrderResponse) UnmarshalBinary(b []byte) error {
	var res OrderResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
