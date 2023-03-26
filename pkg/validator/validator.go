package validator

import (
	"encoding/json"
	"net/url"

	v "github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

func ValidatBodyAndGetData[T interface{}](data []byte) (T, error) {
	var s T
	validate := v.New()
	err := json.Unmarshal(data, &s)

	if err != nil {
		return s, err
	}

	err = validate.Struct(s)

	if err != nil {
		return s, err
	}

	return s, nil
}

func ValidatQueryeAndGetData[T interface{}](values url.Values) (*T, error) {
	var s T
	validate := v.New()
	var decoder = schema.NewDecoder()

	err := decoder.Decode(&s, values)

	if err != nil {
		return nil, err
	}

	validate.Struct(s)

	if err != nil {
		return nil, err
	}

	return &s, nil
}
