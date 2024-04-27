package server

import (
	"database/sql"
	"log"
	"net/http"
	"restful/exception"
	"restful/layer/cart"
	"restful/layer/orders"
	"restful/layer/produk"
	"restful/layer/users"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type Server struct {
	Addr string
	Db   *sql.DB
}

func New(addr string, db *sql.DB) *Server {
	return &Server{
		Addr: addr,
		Db:   db,
	}
}

func (s *Server) Run() error {
	var validate *validator.Validate = validator.New()

	router := httprouter.New()
	router.NotFound = exception.NotFound(exception.NotFoundHandler())

	if validate == nil {
		log.Fatalf("validate is %v", validate)
	}

	userRepo := users.NewRepository()
	userService := users.NewService(userRepo, s.Db)
	userHandler := users.NewRoute(userService, validate)
	userHandler.RegisterRoute(router)

	prRepo := produk.NewRepository()
	prService := produk.NewService(prRepo, s.Db)
	prHandler := produk.NewRoute(prService, validate)
	prHandler.RegisterRoute(router)

	orderRepo := orders.NewRepository()
	orderService := orders.NewService(s.Db, orderRepo)
	cartService := cart.NewService(orderService)
	cartHandler := cart.NewRoute(cartService, validate, prService)
	cartHandler.RegisterRoute(router)

	return http.ListenAndServe(s.Addr, router)

}
