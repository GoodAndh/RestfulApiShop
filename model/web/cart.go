package web



type CartItems struct {
	Quantity int `json:"quantity" validate:"required"`
	ProductId int `json:"productid" validate:"required"`
}

type CartItemsCheckout struct {
	Items []CartItems `json:"items" validate:"required"`
}