package web

type OrdersCreatePayload struct {
	Total  int    `json:"total_order" validate:"required"`
	Status string `json:"-" validate:"required,oneof=pending sukses cancel"`
	Userid int    `json:"userid" validate:"required"`
}

type OrderItemsCreatePayload struct {
	Total      int     `json:"total_order" validate:"required"`
	TotalPrice float64 `json:"total_price" validate:"required"`
	OrderId    int     `json:"orderid" validate:"required"`
	ProductId  int     `json:"productID" validate:"required"`
	Userid     int     `json:"userid" validate:"required"`
}

type Orders struct {
	Id     int    `json:"id"`
	Total  int    `json:"total"`
	Status string `json:"status"`
	Userid int    `json:"userid"`
}

type OrdersItems struct {
	Id         int     `json:"id"`
	Total      int     `json:"total"`
	TotalPrice float64 `json:"total_price"`
	OrderId    int     `json:"orderid"`
	ProductId  int     `json:"productID"`
	Userid     int     `json:"userid"`
}

type OrdersUpdatePayload struct {
	Id int `json:"id" validate:"required"`
	Total  int    `json:"total_order" validate:"required"`
	Status string `json:"-" validate:"required,oneof=pending sukses cancel"`
	Userid int    `json:"userid" validate:"required"`
}
