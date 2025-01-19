package errors

import (
	"fmt"
	"log/slog"
	"net/http"
)

var (
	INTERNAL_SERVER_ERROR = NewAPIError(http.StatusInternalServerError, "internal server error")
	EXCEEDED_MAX_ITEM     = NewAPIError(http.StatusBadRequest, "you can order a maximum of 50 different products per request. Please reduce the number of items and try again")
	CANNOT_FIND_PRODUCTS  = NewAPIError(http.StatusBadRequest, "can't find products with specified ids")
	INVALID_PRODUCT_ID    = NewAPIError(http.StatusBadRequest, "invalid product id")
	DUPLICATE_PRODUCT_ID  = NewAPIError(http.StatusBadRequest, "duplicate items with the same product_id are not allowed. Please increase the quantity instead.")
)

type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func NewAPIError(statusCode int, message string) *APIError {
	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func PrintAndReturnErr(err error, apiError *APIError) (e *APIError) {
	if err != nil {
		slog.Error(err.Error())
	}
	return apiError
}
