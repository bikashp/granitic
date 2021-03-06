// Copyright 2016 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package ws

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

// NewWsParamsForPath creates a WsParams used to store the elements of a request
// path extracted using regular expression groups.
func NewWsParamsForPath(targets []string, values []string) *WsParams {

	contents := make(url.Values)
	v := len(values)
	var names []string

	for i, k := range targets {

		if i < v {
			contents[strconv.Itoa(i)] = []string{values[i]}
			names = append(names, k)
		}

	}

	p := new(WsParams)
	p.values = contents
	p.paramNames = names

	return p

}

// NewWsParamsForQuery creates a WsParams storing the HTTP query parameters from a request.
func NewWsParamsForQuery(values url.Values) *WsParams {

	wp := new(WsParams)
	wp.values = values

	var names []string

	for k, _ := range values {
		names = append(names, k)
	}

	wp.paramNames = names

	return wp

}

// An abstraction of the HTTP query parameters or path parameters with type-safe accessors.
type WsParams struct {
	values     url.Values
	paramNames []string
}

// ParamNames returns the names of all of the parameters stored
func (wp *WsParams) ParamNames() []string {
	return wp.paramNames
}

// NotEmpty returns true if a parameter with the supplied name exists and has a non-empty string representation.
func (wp *WsParams) NotEmpty(key string) bool {

	if !wp.Exists(key) {
		return false
	}

	s, _ := wp.StringValue(key)

	return s != ""

}

// Exists returns true if a parameter with the supplied name exists, even if that parameter's value is an empty string.
func (wp *WsParams) Exists(key string) bool {
	return wp.values[key] != nil
}

// MultipleValues returns true if the parameter with the supplied name was set more than once (allowed for HTTP query parameters).
func (wp *WsParams) MultipleValues(key string) bool {

	value := wp.values[key]

	return value != nil && len(value) > 1

}

// StringValue returns the string representation of the specified parameter or an error if no value exists for that parameter.
func (wp *WsParams) StringValue(key string) (string, error) {

	s := wp.values[key]

	if s == nil {
		return "", wp.noVal(key)
	}

	return s[len(s)-1], nil

}

// BoolValue returns the bool representation of the specified parameter (using Go's bool conversion rules) or an error if no value exists for that parameter.
func (wp *WsParams) BoolValue(key string) (bool, error) {

	v := wp.values[key]

	if v == nil {
		return false, wp.noVal(key)
	}

	b, err := strconv.ParseBool(v[len(v)-1])

	return b, err

}

// FloatNValue returns a float representation of the specified parameter with the specified bit size, or an error if no value exists for that parameter or
// if the value could not be converted to a float.
func (wp *WsParams) FloatNValue(key string, bits int) (float64, error) {

	v := wp.values[key]

	if v == nil {
		return 0.0, wp.noVal(key)
	}

	i, err := strconv.ParseFloat(v[len(v)-1], bits)

	return i, err

}

// IntNValue returns a signed int representation of the specified parameter with the specified bit size, or an error if no value exists for that parameter or
// if the value could not be converted to an int.
func (wp *WsParams) IntNValue(key string, bits int) (int64, error) {

	v := wp.values[key]

	if v == nil {
		return 0, wp.noVal(key)
	}

	i, err := strconv.ParseInt(v[len(v)-1], 10, bits)

	return i, err

}

// UIntNValue returns an unsigned int representation of the specified parameter with the specified bit size, or an error if no value exists for that parameter or
// if the value could not be converted to an unsigned int.
func (wp *WsParams) UIntNValue(key string, bits int) (uint64, error) {

	v := wp.values[key]

	if v == nil {
		return 0, wp.noVal(key)
	}

	i, err := strconv.ParseUint(v[len(v)-1], 10, bits)

	return i, err

}

func (wp *WsParams) noVal(key string) error {
	message := fmt.Sprintf("No value available for key %s", key)
	return errors.New(message)
}
