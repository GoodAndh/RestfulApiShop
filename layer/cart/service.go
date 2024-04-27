package cart

import (
	"context"
	"fmt"
	"restful/layer/orders"
	"restful/model/web"
)

type ServiceImpl struct {
	order orders.Service
}

func NewService(order orders.Service) *ServiceImpl {
	return &ServiceImpl{
		order: order,
	}
}

func (s *ServiceImpl) CreateOrder(ctx context.Context, ps []web.Product, cartItems []web.CartItems, userId int) (int, float64, error) {

	//store product into map for easier access
	pMap := make(map[int]web.Product)
	for _, items := range ps {
		pMap[items.Id] = items
	}

	//check if all product is available

	if err := checkIfIsAvailable(cartItems, pMap); err != nil {
		return 0, 0, err
	}

	//calculate total price
	totalPrice, quantity := calculateTotalPrice(cartItems, pMap)

	//create Order
	orderId, err := s.order.CreateOrders(ctx, &web.OrdersCreatePayload{
		Total:  quantity,
		Status: "pending",
		Userid: userId,
	})
	if err != nil {
		return 0, 0, err
	}

	for _, items := range cartItems {
		err := s.order.CreateOrderItems(ctx, &web.OrderItemsCreatePayload{
			Total:      quantity,
			TotalPrice: totalPrice,
			OrderId:    orderId,
			ProductId:  items.ProductId,
			Userid:     userId,
		})
		if err != nil {
			return 0, 0, err
		}
	}

	return orderId, totalPrice, nil

}

func checkIfIsAvailable(cart []web.CartItems, pMap map[int]web.Product) error {
	if len(cart) <= 0 {
		return fmt.Errorf("cart is empty,try order something")
	}

	for _, items := range cart {
		product, ok := pMap[items.ProductId]
		if !ok {
			return fmt.Errorf(" product id : %d tidak ditemukan ", items.ProductId)
		}

		if product.Quantity < items.Quantity {
			return fmt.Errorf("pesanan anda melebihi jumlah stock yang tersedia ,pesanan anda : %d sisa stock tersedia : %d ", items.Quantity, product.Quantity)
		}
	}

	return nil

}

func calculateTotalPrice(cartItem []web.CartItems, Pmap map[int]web.Product) (float64, int) {
	var total float64
	var quantity int

	for _, items := range cartItem {
		product := Pmap[items.ProductId]
		total += product.Price * float64(items.Quantity)
		quantity = items.Quantity
	}

	return total, quantity
}
