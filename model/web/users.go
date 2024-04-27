package web


type Users struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Password string `json:"-"`
	Email string `json:"email"`
}

type UsersLoginPayload struct {
	Email string `json:"email"  validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
type UsersRegisterPayload struct {
	Name string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	VPassword string `json:"vpassword" validate:"required,eqfield=Password"`
	Email string `json:"email" validate:"required,email"`
}
