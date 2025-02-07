package resty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// ParamHandler provides access tp parameters in the http.Request object
type ParamHandler interface {
	Param(key string) string
	Query(key string) string
	Queries(key string) []string
	ScanBody(any) error
}

type Errors interface {
	AddError(err string)
	HasError() bool
}

// ResourceFunc is used as a result of ActionFunc when returning long binary streams
type ResourceFunc func() (contenType string, contentLength uint64, reader io.Reader)

// ActionHandler performs business logic
// it is not intended this function to has access to http.Request or http.ResponseWriter
// Returns:
//  1. The result data. It may be struct or array or nil
//  2. HTTP status code
//  3. error (if any) or nil
type ActionFunc func(ParamHandler, context.Context) (interface{}, int, error)

// ValidationFunk represents validation logic, executed before given ActionFunc
type ValidationFunk func(ParamHandler, context.Context, Errors) error

type errorsList struct {
	Errors []string `json:"errors"`
}

func (el *errorsList) AddError(err string) {
	el.Errors = append(el.Errors, err)
}

func (el *errorsList) HasError() bool {
	return len(el.Errors) > 0
}

func validate(action ActionFunc, w http.ResponseWriter, ph ParamHandler, ctx context.Context) bool {
	validators := getValidators(action)
	// perform Validate logic
	if validators != nil {
		el := NewErrors()
		for _, validator := range validators {
			if err := validator(ph, ctx, el); err != nil {
				// there is fatal error. Returning code 500
				log.Printf("Error executing validation function %#v for %#v: %v", validator, action, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return false
			}
		}
		if el.HasError() {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(el)
			return false
		}
	}
	return true
}

func sendReader(code int, reader io.Reader, w http.ResponseWriter) {
	buffer := make([]byte, 512)
	if count, err := reader.Read(buffer); err != nil && err != io.EOF {
		// react on error
		log.Printf("Error reading from reader object: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else if count > 0 {
		w.Header().Set("Content-Type", http.DetectContentType(buffer))
		w.WriteHeader(code)
		w.Write(buffer[:count])
		io.Copy(w, reader)
	} else {
		http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
	}
}

// HandleAction takes ActionHandler and converts it to http.HandlerFunc
func HandleAction(action ActionFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		ph := &paramHandler{r: r}
		ctx := r.Context()
		if validate(action, w, ph, ctx) {
			result, code, err := action(ph, ctx)
			if err != nil {
				// there is fatal error. Returning code 500, ignoring the code returned by the function
				log.Printf("Error executing action %v: %v", action, err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			} else if result != nil {
				if code <= 0 {
					code = http.StatusOK
				}
				if s, ok := result.(string); ok {
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
					w.WriteHeader(code)
					w.Write([]byte(s))
				} else if resource, ok := result.(func() (contenType string, contentLength uint64, reader io.Reader)); ok {
					contentType, contentLength, reader := resource()
					w.Header().Set("Content-Type", contentType)
					if contentLength > 0 {
						w.Header().Set("Content-Length", fmt.Sprintf("%v", contentLength))
					}
					io.Copy(w, reader)
				} else if reader, ok := result.(io.Reader); ok {
					sendReader(code, reader, w)
				} else {
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
					w.WriteHeader(code)
					json.NewEncoder(w).Encode(result)
				}
			} else {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				if code > 0 {
					w.WriteHeader(code)
				} else {
					w.WriteHeader(http.StatusNoContent)
				}
			}
		}
	}
}

func NewErrors(errs ...string) Errors {
	r := &errorsList{Errors: make([]string, 0)}
	r.Errors = append(r.Errors, errs...)
	return r
}

func StatusError(status int) (interface{}, int, error) {
	return http.StatusText(status), status, nil
}

func StatusErrorNotFound() (interface{}, int, error) {
	return StatusError(http.StatusNotFound)
}

func StatusErrorBadRequest() (interface{}, int, error) {
	return StatusError(http.StatusBadRequest)
}

func StatusErrorInternalServerError() (interface{}, int, error) {
	return StatusError(http.StatusInternalServerError)
}

func StatusOK(result interface{}) (interface{}, int, error) {
	if result != nil {
		return result, http.StatusOK, nil
	}
	return nil, http.StatusNoContent, nil
}
