package utils

import (
	"restful/model/domain"
	"restful/model/web"
)

func ConvertUserToSlice(user *domain.Users) *web.Users {
	return &web.Users{
		Id:       user.Id,
		Name:     user.Name,
		Password: user.Password,
		Email:    user.Email,
	}
}

func ConvertProductToWeb(pr *domain.Product) *web.Product {
	return &web.Product{
		Id:        pr.Id,
		Name:      pr.Name,
		Deskripsi: pr.Deskripsi,
		Category:  pr.Category,
		Quantity:  pr.Quantity,
		Price:     pr.Price,
		Userid:    pr.Userid,
	}
}
func ConvertProductToSlice(pr []domain.Product) []web.Product {
	p := []web.Product{}
	for _, v := range pr {
		p = append(p, *ConvertProductToWeb(&v))
	}
	return p
}

func ConvertOrdersToWeb(ord *domain.Orders) *web.Orders {
	return &web.Orders{
		Id:     ord.Id,
		Total:  ord.Total,
		Status: ord.Status,
		Userid: ord.Userid,
	}
}

func ConvertOrdersToSlice(ord []domain.Orders) []web.Orders {
	or := []web.Orders{}
	for _, v := range ord {
		or = append(or, *ConvertOrdersToWeb(&v))
	}
	return or
}

func ConvertOrderItemsToWeb(ord *domain.OrdersItems) *web.OrdersItems {
	return &web.OrdersItems{
		Id:         ord.Id,
		Total:      ord.Total,
		TotalPrice: ord.TotalPrice,
		OrderId:    ord.OrderId,
		ProductId:  ord.ProductId,
		Userid:     ord.UserId,
	}
}

func ConvertOrderItemsToSlice(ord []domain.OrdersItems) []web.OrdersItems {
	o := []web.OrdersItems{}
	for _, v := range ord {
		o = append(o, *ConvertOrderItemsToWeb(&v))
	}
	return o
}
