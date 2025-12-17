package vega

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

// Validation provides form validation functionality
type Validation struct {
	Data   url.Values
	Errors map[string]string
}

// Validator creates a new validation instance
func (v *Vega) Validator(data url.Values) *Validation {
	return &Validation{
		Errors: make(map[string]string),
		Data:   data,
	}
}

// Valid returns true if there are no validation errors
func (v *Validation) Valid() bool {
	return len(v.Errors) == 0
}

// AddError adds an error message for a field
func (v *Validation) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Has checks if a field exists in the request
func (v *Validation) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}
	return true
}

// Required validates that fields are not empty
func (v *Validation) Required(r *http.Request, fields ...string) {
	for _, field := range fields {
		value := r.Form.Get(field)
		if strings.TrimSpace(value) == "" {
			v.AddError(field, "This field cannot be blank")
		}
	}
}

// Check adds an error if the condition is false
func (v *Validation) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// IsEmail validates that a value is a valid email address
func (v *Validation) IsEmail(field, value string) {
	if !govalidator.IsEmail(value) {
		v.AddError(field, "Invalid email address")
	}
}

// IsInt validates that a value is an integer
func (v *Validation) IsInt(field, value string) {
	_, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(field, "This field must be an integer")
	}
}

// IsFloat validates that a value is a floating point number
func (v *Validation) IsFloat(field, value string) {
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.AddError(field, "This field must be a floating point number")
	}
}

// IsDateISO validates that a value is a date in ISO format (YYYY-MM-DD)
func (v *Validation) IsDateISO(field, value string) {
	_, err := time.Parse("2006-01-02", value)
	if err != nil {
		v.AddError(field, "This field must be a date in the form of YYYY-MM-DD")
	}
}

// NoSpaces validates that a value contains no whitespace
func (v *Validation) NoSpaces(field, value string) {
	if govalidator.HasWhitespace(value) {
		v.AddError(field, "Spaces are not permitted")
	}
}

