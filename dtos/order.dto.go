package dtos

type OrderItem struct {
	ProductId int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type OrderRequest struct {
	Order struct {
		Items []OrderItem `json:"items"`
	} `json:"order"`
}

type ItemDetail struct {
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	VAT       float64 `json:"vat"`
}

type OrderResponse struct {
	OrderID    int          `json:"order_id"`
	OrderPrice float64      `json:"order_price"`
	OrderVAT   float64      `json:"order_vat"`
	Items      []ItemDetail `json:"items"`
}
