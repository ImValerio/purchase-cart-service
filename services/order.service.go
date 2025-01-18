package services

import (
	"context"
	"errors"
	"purchase-cart-service/models"
	"purchase-cart-service/repositories"
)

type OrderService struct {
	OrderRepo     repositories.OrderRepository
	OrderItemRepo repositories.OrderItemRepository
	ProductRepo   repositories.ProductRepository
}

func (os OrderService) CreateOrder(ctx context.Context, orderRequest *models.OrderRequest) (*models.OrderResponse, error) {

	if len(orderRequest.Order.Items) > 50 {
		return nil, errors.New("you can order a maximum of 50 different products per request. Please reduce the number of items and try again")
	}

	productIds := getProductsIds(orderRequest)

	products, err := os.ProductRepo.GetByIDs(ctx, productIds)
	if err != nil {
		return nil, err
	}

	productMap := getProductMapByModel(products)

	orderId, err := os.OrderRepo.Insert(ctx)
	if err != nil {
		return nil, err
	}

	// var items []*models.OrderItem
	// for _, request := range orderRequest.Order.Items {
	// 	product, ok := productMap[request.ProductId]

	// 	if ok {
	// 		item, err := os.OrderItemRepo.Insert(ctx, *orderId, request, product)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		items = append(items, item)
	// 	} else {
	// 		return nil, err
	// 	}
	// }

	items, err := os.OrderItemRepo.InsertBatch(ctx, *orderId, orderRequest.Order.Items, productMap)
	if err != nil {
		return nil, err
	}

	itemsDetail := []models.ItemDetail{}

	totalPrice, totalVAT := 0.0, 0.0
	for _, item := range items {
		var itemDetail models.ItemDetail

		itemDetail.Quantity = item.Quantity
		itemDetail.ProductID = item.ProductId
		itemDetail.Price = item.Price
		itemDetail.VAT = item.VAT

		itemsDetail = append(itemsDetail, itemDetail)
		totalPrice += itemDetail.Price
		totalVAT += itemDetail.VAT
	}

	rv := &models.OrderResponse{}

	rv.OrderID = *orderId
	rv.OrderPrice = totalPrice
	rv.OrderVAT = totalVAT
	rv.Items = itemsDetail

	return rv, nil

}

func getProductsIds(orderRequest *models.OrderRequest) []int {
	var rv []int

	for _, item := range orderRequest.Order.Items {

		rv = append(rv, item.ProductId)
	}

	return rv
}

func getProductMapByModel(products []models.Product) map[int]models.Product {
	rv := make(map[int]models.Product)

	for _, product := range products {
		rv[product.Id] = product
	}

	return rv
}
