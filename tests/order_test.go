package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"purchase-cart-service/models"
	"testing"
	"time"
)

// Mock OrderService
type MockOrderService struct {
	// Use a simple map to mock behavior
	mockBehavior func(ctx context.Context, orderRequest *models.OrderRequest) (*models.OrderResponse, error)
}

func (m *MockOrderService) CreateOrder(ctx context.Context, orderRequest *models.OrderRequest) (*models.OrderResponse, error) {
	return m.mockBehavior(ctx, orderRequest)
}

// Handler to test
func createOrderHandler(service *MockOrderService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		var body models.OrderRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			throwError(w, "invalid request body", http.StatusBadRequest)
			return
		}

		order, err := service.CreateOrder(ctx, &body)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				throwError(w, "request timeout", http.StatusRequestTimeout)
				return
			}
			throwError(w, err.Error(), http.StatusBadRequest)
			return
		}

		writeJSON(w, http.StatusOK, order)
	}
}

func TestOrderEndpoint(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := &MockOrderService{
			mockBehavior: func(ctx context.Context, orderRequest *models.OrderRequest) (*models.OrderResponse, error) {
				return &models.OrderResponse{
					OrderID:    123,
					OrderPrice: 200.0,
					OrderVAT:   40.0,
					Items: []models.ItemDetail{
						{ProductID: 1, Quantity: 2, Price: 100.0, VAT: 20.0},
					},
				}, nil
			},
		}

		handler := createOrderHandler(mockService)

		requestPayload := &models.OrderRequest{
			Order: struct {
				Items []models.OrderItemDto `json:"items"`
			}{
				Items: []models.OrderItemDto{
					{ProductId: 1, Quantity: 2},
				},
			},
		}

		body, _ := json.Marshal(requestPayload)
		req := httptest.NewRequest("POST", "/order", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status code 200, got %d", rec.Code)
		}

		var response models.OrderResponse
		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response.OrderID != 123 {
			t.Errorf("expected OrderID 123, got %d", response.OrderID)
		}
	})

	t.Run("failure - invalid productId", func(t *testing.T) {
		mockService := &MockOrderService{
			mockBehavior: func(ctx context.Context, orderRequest *models.OrderRequest) (*models.OrderResponse, error) {
				return nil, errors.New("invalid products in items list")
			},
		}

		handler := createOrderHandler(mockService)

		req := httptest.NewRequest("POST", "/order", bytes.NewBuffer([]byte("invalid json")))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status code 400, got %d", rec.Code)
		}
	})
	t.Run("failure - invalid body", func(t *testing.T) {
		mockService := &MockOrderService{}

		handler := createOrderHandler(mockService)

		req := httptest.NewRequest("POST", "/order", bytes.NewBuffer([]byte("invalid json")))
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status code 500, got %d", rec.Code)
		}
	})

}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	if err := enc.Encode(v); err != nil {
		http.Error(w, `{"error":"failed to encode JSON"}`, http.StatusInternalServerError)
	}
}

func throwError(w http.ResponseWriter, message string, status int) {
	writeJSON(w, status, map[string]string{"error": message})
}
