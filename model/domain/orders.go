package domain

type Orders struct {
	Id int 
	Total int
	Status string 
	Userid int
}

type OrdersItems struct {
	Id int
	Total int
	TotalPrice float64
	OrderId int
	ProductId int
	UserId int
}

