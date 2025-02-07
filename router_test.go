package resty_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mmitevski/resty"
)

type Data struct {
	Id    string
	Count int
}

func TestRouterGet(t *testing.T) {
	router := resty.NewRouter().
		Get("/api/data", func(ph resty.ParamHandler, ctx context.Context) (interface{}, int, error) {
			return &Data{
				Id:    "ID1",
				Count: 1,
			}, http.StatusOK, nil
		})
	r, _ := http.NewRequest(http.MethodGet, "/api/data", nil)
	{
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, r)
		if rr.Code != http.StatusOK {
			t.Fatalf("Unexpected response code for GET method: %v", rr.Code)
		}
		if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Fatalf("Unexpected response Content-Type for GET method: %v", ct)
		}
		var data Data
		if err := json.NewDecoder(rr.Body).Decode(&data); err != nil {
			t.Fatalf("Error decoding JSON object")
		}
		if data.Id != "ID1" {
			t.Errorf("Wrong ID returned from GET method: %v", data.Id)
		}
	}
}

func TestRouterPost(t *testing.T) {
	router := resty.NewRouter().
		Post("/api/data", func(ph resty.ParamHandler, ctx context.Context) (interface{}, int, error) {
			var data Data
			if err := ph.ScanBody(&data); err != nil {
				return err, http.StatusBadRequest, nil
			}
			data.Id = "ID1"
			return &data, http.StatusCreated, nil
		})
	var data Data
	data.Count = 1
	jd, _ := json.Marshal(data)
	r, _ := http.NewRequest(http.MethodPost, "/api/data", bytes.NewBuffer(jd))
	{
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, r)
		if rr.Code != http.StatusCreated {
			t.Fatalf("Unexpected response code for POST method: %v", rr.Code)
		}
		if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Fatalf("Unexpected response Content-Type for POST method: %v", ct)
		}
		if err := json.NewDecoder(rr.Body).Decode(&data); err != nil {
			t.Fatalf("Error decoding JSON object")
		}
		if data.Id != "ID1" {
			t.Errorf("Wrong ID returned from POST method: %v", data.Id)
		}
		if data.Count != 1 {
			t.Errorf("Wrong Count returned from POST method: %v", data.Count)
		}
	}
}

func TestRouterPut(t *testing.T) {
	router := resty.NewRouter().
		Put("/api/data", func(ph resty.ParamHandler, ctx context.Context) (interface{}, int, error) {
			var data Data
			if err := ph.ScanBody(&data); err != nil {
				return err, http.StatusBadRequest, nil
			}
			data.Count++
			return &data, http.StatusOK, nil
		})
	data := Data{Id: "ID2", Count: 1}
	jd, _ := json.Marshal(data)
	r, _ := http.NewRequest(http.MethodPut, "/api/data", bytes.NewBuffer(jd))
	{
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, r)
		if rr.Code != http.StatusOK {
			t.Fatalf("Unexpected response code for PUT method: %v", rr.Code)
		}
		if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
			t.Fatalf("Unexpected response Content-Type for PUT method: %v", ct)
		}
		if err := json.NewDecoder(rr.Body).Decode(&data); err != nil {
			t.Fatalf("Error decoding JSON object")
		}
		if data.Id != "ID2" {
			t.Errorf("Wrong ID returned from PUT method: %v", data.Id)
		}
		if data.Count != 2 {
			t.Errorf("Wrong Count returned from PUT method: %v", data.Count)
		}
	}
}
