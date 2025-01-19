package services

import (
	"context"
	"purchase-cart-service/errors"
	"purchase-cart-service/models"
	"purchase-cart-service/repositories"
)

type OrderService struct {
	OrderRepo     repositories.OrderRepository
	OrderItemRepo repositories.OrderItemRepository
	ProductRepo   repositories.ProductRepository
}

func (ordsrv OrderService) CreateOrder(ctx context.Context, orderRequest *models.OrderRequest) (*models.OrderResponse, error) {

	if len(orderRequest.Order.Items) > 50 {
		return nil, errors.EXCEEDED_MAX_ITEM
	}

	productIds, err := getProductsIds(orderRequest)
	if err != nil {
		return nil, err
	}

	products, err := ordsrv.ProductRepo.GetByIDs(ctx, productIds)
	if err != nil {
		return nil, errors.PrintAndReturnErr(err, errors.INTERNAL_SERVER_ERROR)
	}

	if len(products) == 0 {
		return nil, errors.CANNOT_FIND_PRODUCTS
	}

	productMap := getProductMapByModel(products)

	orderId, err := ordsrv.OrderRepo.Insert(ctx)
	if err != nil {
		return nil, errors.PrintAndReturnErr(err, errors.INTERNAL_SERVER_ERROR)
	}

	// var items []*models.OrderItem
	// for _, request := range orderRequest.Order.Items {
	// 	product, ok := productMap[request.ProductId]

	// 	if ok {
	// 		item, err := ordsrv.OrderItemRepo.Insert(ctx, *orderId, request, product)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		items = append(items, item)
	// 	} else {
	// 		return nil, err
	// 	}
	// }

	items, err := ordsrv.OrderItemRepo.InsertBatch(ctx, *orderId, orderRequest.Order.Items, productMap)
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

func getProductsIds(orderRequest *models.OrderRequest) ([]int, error) {
	var rv []int

	alreadySeen := make(map[int]bool)

	for _, item := range orderRequest.Order.Items {
		if alreadySeen[item.ProductId] {
			return nil, errors.DUPLICATE_PRODUCT_ID
		}

		if item.ProductId < 0 {
			return nil, errors.INVALID_PRODUCT_ID
		}

		rv = append(rv, item.ProductId)
		alreadySeen[item.ProductId] = true
	}

	return rv, nil
}

func getProductMapByModel(products []models.Product) map[int]models.Product {
	rv := make(map[int]models.Product)

	for _, product := range products {
		rv[product.Id] = product
	}

	return rv
}
