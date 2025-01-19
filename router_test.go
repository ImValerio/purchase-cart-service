package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"purchase-cart-service/dtos"
	"testing"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	environment := os.Getenv("ENV")
	if environment != "prod" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	t.Helper()
	db, err := sql.Open("postgres", os.Getenv("TEST_DATABASE_URL"))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

func setupTestRouter(t *testing.T, db *sql.DB) *chi.Mux {
	t.Helper()

	router := createRouterAndSetRoutes()
	return router
}

// Helper function to send HTTP request and return response
func performRequest(router http.Handler, method, path string, body []byte) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}

// Test case: successful order creation
func TestCreateOrder_Success(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(t, db)

	orderRequest := dtos.OrderRequest{
		Order: struct {
			Items []dtos.OrderItem `json:"items"`
		}{
			Items: []dtos.OrderItem{
				{ProductId: 1, Quantity: 2},
				{ProductId: 2, Quantity: 3},
			},
		},
	}

	reqBody, _ := json.Marshal(orderRequest)
	resp := performRequest(router, "POST", "/order", reqBody)

	assert.Equal(t, http.StatusOK, resp.Code)

	var orderResponse dtos.OrderResponse
	err := json.Unmarshal(resp.Body.Bytes(), &orderResponse)
	assert.NoError(t, err)
	assert.NotZero(t, orderResponse.OrderID)
	assert.Equal(t, 2, len(orderResponse.Items))
}

// Test case: exceeding item limit (more than 50 items)
func TestCreateOrder_ExceededItemLimit(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(t, db)

	orderRequest := dtos.OrderRequest{
		Order: struct {
			Items []dtos.OrderItem `json:"items"`
		}{
			Items: make([]dtos.OrderItem, 51),
		},
	}

	reqBody, _ := json.Marshal(orderRequest)
	resp := performRequest(router, "POST", "/order", reqBody)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "you can order a maximum of 50 different products")
}

// Test case: product not found
func TestCreateOrder_ProductNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(t, db)

	orderRequest := dtos.OrderRequest{
		Order: struct {
			Items []dtos.OrderItem `json:"items"`
		}{
			Items: []dtos.OrderItem{
				{ProductId: 9999, Quantity: 1}, // Non-existent product ID
			},
		},
	}

	reqBody, _ := json.Marshal(orderRequest)
	resp := performRequest(router, "POST", "/order", reqBody)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "can't find products with specified ids")
}

// Test case: invalid product ID
func TestCreateOrder_InvalidProductID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	router := setupTestRouter(t, db)

	orderRequest := dtos.OrderRequest{
		Order: struct {
			Items []dtos.OrderItem `json:"items"`
		}{
			Items: []dtos.OrderItem{
				{ProductId: -1, Quantity: 5}, // Invalid negative product ID
			},
		},
	}

	reqBody, _ := json.Marshal(orderRequest)
	resp := performRequest(router, "POST", "/order", reqBody)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "invalid product id")
}
