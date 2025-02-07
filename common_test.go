package resty

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func checkResponseCode(t *testing.T, expected int, rr *httptest.ResponseRecorder) {
	if expected != rr.Code {
		t.Errorf("Expected response code %d. Got %d\n", expected, rr.Code)
	}
}

func checkResponseBody(t *testing.T, expected string, rr *httptest.ResponseRecorder) {
	actual := strings.TrimSpace(fmt.Sprintf("%v", rr.Body))
	if expected != actual {
		t.Errorf("Expected response body '%s'. Got '%s'\n", expected, actual)
	}
}

func checkResponseType(t *testing.T, expected string, rr *httptest.ResponseRecorder) {
	actual := rr.Header().Get("Content-Type")
	if expected != actual {
		t.Errorf("Expected response Content-Type '%s'. Got '%s'\n", expected, actual)
	}
}

func TestHandleActionErrors(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	HandleAction(func(ph ParamHandler, c context.Context) (any, int, error) {
		return nil, 0, errors.New("this is simply an error, returned by the action")
	}).ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusInternalServerError, rr)
	checkResponseType(t, "text/plain; charset=utf-8", rr)
	checkResponseBody(t, http.StatusText(rr.Code), rr)
}

func TestHandleActionPlainText(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	HandleAction(func(ph ParamHandler, c context.Context) (any, int, error) {
		return "test", http.StatusOK, nil
	}).ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusOK, rr)
	checkResponseType(t, "text/plain; charset=utf-8", rr)
	checkResponseBody(t, "test", rr)
}

func TestHandleActionStruct(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	data := struct {
		Key string `json:"key"`
	}{Key: "test"}
	HandleAction(func(ph ParamHandler, c context.Context) (any, int, error) {
		return data, http.StatusOK, nil
	}).ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusOK, rr)
	checkResponseType(t, "application/json; charset=utf-8", rr)
	checkResponseBody(t, `{"key":"test"}`, rr)
}

func TestHandleActionMap(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	data := make(map[string]string)
	data["key"] = "test"
	HandleAction(func(ph ParamHandler, c context.Context) (any, int, error) {
		return &data, http.StatusOK, nil
	}).ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusOK, rr)
	checkResponseType(t, "application/json; charset=utf-8", rr)
	checkResponseBody(t, `{"key":"test"}`, rr)
}

func TestHandleActionArray(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	data := make([]int, 0)
	data = append(data, 123)
	data = append(data, 456)
	HandleAction(func(ph ParamHandler, c context.Context) (any, int, error) {
		return &data, http.StatusOK, nil
	}).ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusOK, rr)
	checkResponseType(t, "application/json; charset=utf-8", rr)
	checkResponseBody(t, "[123,456]", rr)
}

func TestHandleActionReader(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	reader := strings.NewReader("This is a test!")
	HandleAction(func(ph ParamHandler, c context.Context) (any, int, error) {
		return reader, http.StatusOK, nil
	}).ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusOK, rr)
	checkResponseType(t, "application/octet-stream", rr)
	checkResponseBody(t, "This is a test!", rr)
}

func TestHandleActionResource(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)
	const s = "This is a test!"
	const tp = "text/plain; charset=utf=8"
	f := func() (contenType string, contentLength uint64, reader io.Reader) {
		return tp, uint64(len(s)), strings.NewReader(s)
	}

	rr := httptest.NewRecorder()
	HandleAction(func(ph ParamHandler, c context.Context) (any, int, error) {
		return f, http.StatusOK, nil
	}).ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusOK, rr)
	checkResponseType(t, tp, rr)
	checkResponseBody(t, s, rr)
	if length := rr.Header().Get("Content-Length"); length != "15" {
		t.Errorf("Expected Content-Length %d. Got %s\n", 15, length)
	}
}

func TestHandleActionParams(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/api/123?q=test", nil)
	rr := httptest.NewRecorder()
	router := NewRouter()
	router.Get("/api/:id", func(ph ParamHandler, ctx context.Context) (interface{}, int, error) {
		data := make(map[string]string)
		data["id"] = ph.Param("id")
		data["query"] = ph.Query("q")
		return &data, http.StatusOK, nil
	})
	router.ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusOK, rr)
	checkResponseType(t, "application/json; charset=utf-8", rr)
	checkResponseBody(t, `{"id":"123","query":"test"}`, rr)
}

func TestHandleActionValidatorPositive(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/api/123?q=test", nil)
	action := func(ph ParamHandler, c context.Context) (any, int, error) {
		data := make(map[string]string)
		data["id"] = ph.Param("id")
		data["query"] = ph.Query("q")
		return &data, http.StatusOK, nil
	}
	router := NewRouter()
	router.Get("/api/:id", action)
	rr := httptest.NewRecorder()
	executions := 0
	AddValidator(action, func(ph ParamHandler, c context.Context, errs Errors) error {
		executions++
		return nil
	})
	router.ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusOK, rr)
	checkResponseType(t, "application/json; charset=utf-8", rr)
	checkResponseBody(t, `{"id":"123","query":"test"}`, rr)
	if executions != 1 {
		t.Errorf("Expected validator to be executed once. Got %v\n", executions)
	}
}

func TestHandleActionValidatorNegative(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/api/123?q=test", nil)
	action := func(ph ParamHandler, c context.Context) (any, int, error) {
		return nil, http.StatusOK, nil
	}
	router := NewRouter()
	router.Get("/api/:id", action)
	rr := httptest.NewRecorder()
	executions := 0
	AddValidator(action, func(ph ParamHandler, c context.Context, errs Errors) error {
		executions++
		errs.AddError("Test Failure 1")
		return nil
	})
	AddValidator(action, func(ph ParamHandler, c context.Context, errs Errors) error {
		executions++
		errs.AddError("Test Failure 2")
		return nil
	})
	router.ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusBadRequest, rr)
	checkResponseType(t, "application/json; charset=utf-8", rr)
	checkResponseBody(t, `{"errors":["Test Failure 1","Test Failure 2"]}`, rr)
	if executions != 2 {
		t.Errorf("Expected validator to be executed once. Got %v\n", executions)
	}
}

func TestHandleActionValidatorValidationCallFatalError(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/api/123?q=test", nil)
	action := func(ph ParamHandler, c context.Context) (any, int, error) {
		return nil, http.StatusOK, nil
	}
	router := NewRouter()
	router.Get("/api/:id", action)
	rr := httptest.NewRecorder()
	executions := 0
	AddValidator(action, func(ph ParamHandler, c context.Context, errs Errors) error {
		executions++
		errs.AddError("fail")
		return fmt.Errorf("test fatal error during validation call")
	})
	router.ServeHTTP(rr, r)
	checkResponseCode(t, http.StatusInternalServerError, rr)
	checkResponseType(t, "text/plain; charset=utf-8", rr)
	checkResponseBody(t, http.StatusText(http.StatusInternalServerError), rr)
	if executions != 1 {
		t.Errorf("Expected validator to be executed once. Got %v\n", executions)
	}
}
