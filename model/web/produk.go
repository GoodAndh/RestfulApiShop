package web

type Product struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	Deskripsi string  `json:"deskripsi"`
	Category  string  `json:"category"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Userid    int     `json:"id_user"`
}

type ProductCreatePayload struct {
	Name      string  `json:"name" validate:"required"`
	Deskripsi string  `json:"deskripsi" validate:"required"`
	Category  string  `json:"category" validate:"required,oneof=electric consumable etc"`
	Quantity  int     `json:"quantity" validate:"required,number"`
	Price     float64 `json:"price" validate:"required"`
	Userid    int     `json:"id_user" validate:"required"`
}

type ProductUpdatePayload struct {
	Id        int     `json:"id" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Deskripsi string  `json:"deskripsi" validate:"required"`
	Category  string  `json:"category" validate:"required,oneof=electric consumable etc"`
	Quantity  int     `json:"quantity" validate:"required,number"`
	Price     float64 `json:"price" validate:"required"`
}
