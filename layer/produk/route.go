package produk

import (
	"errors"
	"net/http"
	"restful/exception"
	"restful/layer/auth"
	"restful/model/web"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

type Route struct {
	service  Service
	validate *validator.Validate
}

func NewRoute(service Service, validate *validator.Validate) *Route {
	return &Route{
		service:  service,
		validate: validate,
	}
}

func (h *Route) RegisterRoute(router *httprouter.Router) {
	router.GET("/product/:id", h.handleProduct)
	router.GET("/product/", h.handleAllproduct)

	router.POST("/product/create", h.handleCreate)
}

func (h *Route) handleProduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	nm := params.ByName("id")

	id, err := strconv.Atoi(nm)
	if err != nil {
		errMsg := err.Error()
		if !strings.Contains(errMsg, "create") {
			exception.WriteBadRequest(w, "yang anda masukkan bukanlah angka ")
			return
		}

		pr, err := h.service.GetByName(r.Context(), nm)
		if err != nil {
			exception.WriteInternalError(w, err.Error())
			return
		}
		for _, v := range pr {
			if v.Id == 0 {
				exception.WriteBadRequest(w, "data yg kamu minta tdk ada")
				return
			}
		}
		err = exception.SuccesWriteJson(w, "sukses", pr)
		if err != nil {
			exception.WriteInternalError(w, err.Error())
			return
		}
		return
	}
	pr, err := h.service.GetById(r.Context(), id)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) {
			exception.WriteBadRequest(w, err.Error())
			return
		}
		exception.WriteInternalError(w, err.Error())
		return
	}

	if pr.Id == 0 {
		exception.WriteBadRequest(w, "data yg kamu minta tdk ada")
		return
	}
	err = exception.SuccesWriteJson(w, "sukses", pr)
	if err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}
}

func (h *Route) handleAllproduct(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	pr, err := h.service.GetAllProduk(r.Context())
	if err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}
	if len(pr) <= 0 {
		exception.WriteBadRequest(w, "produk masih kosong")
		return
	}

	err = exception.SuccesWriteJson(w, "sukses", pr)
	if err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}
}

func (h *Route) handleCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	w.Header().Add("Content-Type", "application/json")

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

	//get user id
	userId, ok := ses.Values["data"].(int)
	if !ok {
		exception.WriteBadRequest(w, exception.ErrNoSessionFound.Error())
		return
	}

	payload := &web.ProductCreatePayload{}
	newUserid := userId
	payload.Userid = newUserid

	if err := exception.ParseJson(r, payload); err != nil {
		exception.WriteBadRequest(w, err.Error())
		return
	}

	if err := h.validate.Struct(*payload); err != nil {
		errors := err.(validator.ValidationErrors)
		exception.WriteBadRequest(w, "error validation", map[string]any{
			"error": errors.Error(),
		})
		return
	}
	if err := h.service.CreateProduct(r.Context(), payload); err != nil {
		exception.WriteBadRequest(w, err.Error())
		return
	}

	if err := exception.SuccesWriteJson(w, "sukses menambahkan produk"); err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}

}
