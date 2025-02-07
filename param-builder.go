package resty

import (
	"encoding/json"
	"fmt"
	"io"
)

type paramBuilder struct {
	key, value string
	params     map[string]string
	queries    map[string]string
	body       io.Reader
}

func (pb *paramBuilder) Param(key string) string {
	if pb.params != nil {
		return pb.params[key]
	}
	if pb.key == key {
		return pb.value
	}
	panic("this ParamHandler does not param query parameters")
}

func (pb *paramBuilder) Query(key string) string {
	if pb.queries != nil {
		return pb.queries[key]
	}
	panic("this ParamHandler does not support query parameters")
}

func (pb *paramBuilder) Queries(key string) []string {
	panic("this ParamHandler does not support multiple query parameters")
}

func (pb *paramBuilder) ScanBody(dst any) error {
	if pb.body != nil {
		return json.NewDecoder(pb.body).Decode(dst)
	}
	panic("this ParamHandler does not support body scanning")
}

// NewParamHandlerWithSingleParam създава нов ParamHandler, който съдържа списък от параметри
func NewParamHandlerWithSingleParam(key string, value any) ParamHandler {
	return &paramBuilder{
		key:   key,
		value: fmt.Sprintf("%v", value),
	}
}
