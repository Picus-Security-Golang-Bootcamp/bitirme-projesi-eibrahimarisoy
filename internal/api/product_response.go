// Code generated by go-swagger; DO NOT EDIT.

package api

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ProductResponse product response
//
// swagger:model ProductResponse
type ProductResponse struct {

	// categories
	Categories []strfmt.UUID `json:"categories"`

	// description
	Description string `json:"description,omitempty"`

	// id
	// Format: uuid
	ID strfmt.UUID `json:"id,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// price
	Price float64 `json:"price,omitempty"`

	// sku
	Sku string `json:"sku,omitempty"`

	// slug
	Slug string `json:"slug,omitempty"`

	// stock
	Stock int64 `json:"stock,omitempty"`
}

// Validate validates this product response
func (m *ProductResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCategories(formats); err != nil {
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

func (m *ProductResponse) validateCategories(formats strfmt.Registry) error {
	if swag.IsZero(m.Categories) { // not required
		return nil
	}

	for i := 0; i < len(m.Categories); i++ {

		if err := validate.FormatOf("categories"+"."+strconv.Itoa(i), "body", "uuid", m.Categories[i].String(), formats); err != nil {
			return err
		}

	}

	return nil
}

func (m *ProductResponse) validateID(formats strfmt.Registry) error {
	if swag.IsZero(m.ID) { // not required
		return nil
	}

	if err := validate.FormatOf("id", "body", "uuid", m.ID.String(), formats); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this product response based on context it is used
func (m *ProductResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ProductResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProductResponse) UnmarshalBinary(b []byte) error {
	var res ProductResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}