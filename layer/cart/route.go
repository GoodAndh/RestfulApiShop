package cart

import (
	"fmt"
	"log"
	"net/http"
	"restful/exception"
	"restful/layer/auth"
	"restful/layer/produk"
	"restful/model/web"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	service  *ServiceImpl
	validate *validator.Validate
	produk   produk.Service
}

func NewRoute(service *ServiceImpl, validate *validator.Validate, produk produk.Service) *Route {
	return &Route{
		service:  service,
		validate: validate,
		produk:   produk,
	}
}

func (h *Route) RegisterRoute(router *httprouter.Router) {
	router.POST("/order/product/:productId", h.handleOrder)
	router.GET("/order/find", h.handleFind)
	router.GET("/order/find/:orderid", h.handleOrderid)
	router.POST("/order/find/:orderid/checkout", h.handleOrderidPost)

}

func (h *Route) handleOrder(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	w.Header().Add("Content-Type", "application/json")

	idParams := params.ByName("productId")

	id, err := strconv.Atoi(idParams)
	if err != nil {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "find") {
			exception.WriteBadRequest(w, "yang anda masukkan bukanlah angka ")
			return
		}
	}

	//validate session
	ses, err := auth.SessionStore.Get(r, "lg-ses")
	if err != nil {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	if auten, ok := ses.Values["auten"].(bool); !auten || !ok {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	//get userid
	userId, ok := ses.Values["data"].(int)
	if !ok {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	//get product id
	pR, err := h.produk.GetById(r.Context(), id)
	if err != nil {

		exception.WriteBadRequest(w, err.Error())
		return
	}

	ps := []web.Product{}
	ps = append(ps, *pR)

	cartItems := &web.CartItemsCheckout{}

	if err := exception.ParseJson(r, cartItems); err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}
	cartCheck := []web.CartItems{}

	for _, cart := range cartItems.Items {
		cart.ProductId = id
		cartCheck = append(cartCheck, cart)

	}

	if err := h.validate.Struct(*cartItems); err != nil {
		errors := err.(validator.ValidationErrors)
		exception.WriteBadRequest(w, "error validation", map[string]any{
			"error": errors.Error(),
		})
		return
	}

	orderId, TotalPrice, err := h.service.CreateOrder(r.Context(), ps, cartCheck, userId)
	if err != nil {
		exception.WriteBadRequest(w, err.Error())
		return
	}

	if err := exception.SuccesWriteJson(w, "order sukses dibuat ", map[string]any{
		"total_price": fmt.Sprintf("total harga yang perlu kamu order : %v", TotalPrice),
		"order_id":    fmt.Sprintf("order id :%d", orderId),
	}); err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}

}

func (h *Route) handleFind(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//validate session
	w.Header().Add("Content-Type", "application/json")

	ses, err := auth.SessionStore.Get(r, "lg-ses")
	if err != nil {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	if auten, ok := ses.Values["auten"].(bool); !auten || !ok {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	//get userid
	userId, ok := ses.Values["data"].(int)
	if !ok {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	ord, err := h.service.order.GetOrders(r.Context(), userId, "")
	if err != nil {
		exception.WriteBadRequest(w, err.Error())
		return
	}

	for _, v := range ord {
		ordItem, err := h.service.order.GetOrderItems(r.Context(), userId, v.Id)
		if err != nil {
			exception.WriteBadRequest(w, err.Error())
			return
		}
		for _, items := range ordItem {
			if err := exception.SuccesWriteJson(w, "silahkan lakukan pembayaran pada status pending", map[string]any{
				"order_id":      v.Id,
				"total_pesanan": v.Total,
				"status":        v.Status,
				"total_harga":   items.TotalPrice,
			}); err != nil {
				exception.WriteInternalError(w, err.Error())
				return
			}

		}

	}
}

func (h *Route) handleOrderid(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//validate session
	w.Header().Add("Content-Type", "application/json")

	ses, err := auth.SessionStore.Get(r, "lg-ses")
	if err != nil {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	if auten, ok := ses.Values["auten"].(bool); !auten || !ok {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	//get userid
	userId, ok := ses.Values["data"].(int)
	if !ok {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	idParams := params.ByName("orderid")
	id, err := strconv.Atoi(idParams)
	if err != nil {
		exception.WriteBadRequest(w, "yang anda masukkan bukanlah angka ")
		return
	}
	ord, err := h.service.order.GetOrdersById(r.Context(), userId, "s", id)
	if err != nil {
		exception.WriteBadRequest(w, err.Error())
		return
	}
	for _, v := range ord {
		ordItem, err := h.service.order.GetOrderItems(r.Context(), userId, v.Id)
		if err != nil {
			exception.WriteBadRequest(w, err.Error())
			return
		}
		for _, items := range ordItem {
			pr, err := h.produk.GetById(r.Context(), items.ProductId)
			if err != nil {
				exception.WriteBadRequest(w, err.Error())
				return
			}
			if err := exception.SuccesWriteJson(w, "silahkan lakukan pembayaran", map[string]any{
				"order_id":      v.Id,
				"total_pesanan": v.Total,
				"status":        v.Status,
				"total_harga":   items.TotalPrice,
				"nama_produk":   pr.Name,
				"id_penjual":    pr.Userid,
				"harga_satuan":  pr.Price,
				"stock_tersisa": pr.Quantity,
			}); err != nil {
				exception.WriteInternalError(w, err.Error())
				return
			}
		}
	}
}

func (h *Route) handleOrderidPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	//validate session
	w.Header().Add("Content-Type", "application/json")

	ses, err := auth.SessionStore.Get(r, "lg-ses")
	if err != nil {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	if auten, ok := ses.Values["auten"].(bool); !auten || !ok {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	//get userid
	userId, ok := ses.Values["data"].(int)
	if !ok {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	idParams := params.ByName("orderid")

	id, err := strconv.Atoi(idParams)
	if err != nil {
		log.Println("error=",err," idparams = ",idParams)
		exception.WriteBadRequest(w, "yang anda masukkan bukanlah angka ")
		return
	}

	ord, err := h.service.order.GetOrdersById(r.Context(), userId, "s", id)
	if err != nil {
		exception.WriteBadRequest(w, err.Error())
		return
	}

	for _, v := range ord {
		ordItem, err := h.service.order.GetOrderItems(r.Context(), userId, v.Id)
		if err != nil {
			exception.WriteBadRequest(w, err.Error())
			return
		}
		for _, items := range ordItem {
			pr, err := h.produk.GetById(r.Context(), items.ProductId)
			if err != nil {
				exception.WriteBadRequest(w, err.Error())
				return
			}
			err = h.produk.UpdateProduct(r.Context(), &web.ProductUpdatePayload{
				Id:        pr.Id,
				Name:      pr.Name,
				Deskripsi: pr.Deskripsi,
				Category:  pr.Category,
				Quantity:  pr.Quantity - v.Total,
				Price:     pr.Price,
			})
			if err != nil {
				exception.WriteBadRequest(w, err.Error())
				return
			}
			err = h.service.order.UpdateOrders(r.Context(), &web.OrdersUpdatePayload{
				Id:     v.Id,
				Total:  v.Total,
				Status: "sukses",
				Userid: userId,
			})
			if err != nil {
				exception.WriteBadRequest(w, err.Error())
				return
			}
			if err := exception.SuccesWriteJson(w, "pembayaran berhasil"); err != nil {
				exception.WriteInternalError(w, err.Error())
				return
			}
		}

	}
}
