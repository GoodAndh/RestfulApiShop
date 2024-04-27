package users

import (
	"net/http"
	"restful/exception"
	"restful/layer/auth"
	"restful/model/web"

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
	router.POST("/login/register", h.handleRegister)
	router.POST("/login/", h.handleLogin)
	
	router.POST("/logout/", h.handleLogout)
}

func (h *Route) handleRegister(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	payload := &web.UsersRegisterPayload{}

	if err := exception.ParseJson(r, payload); err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}

	if err := h.validate.Struct(*payload); err != nil {
		errors := err.(validator.ValidationErrors)
		exception.WriteBadRequest(w, "error below",map[string]any{
			"error":errors.Error(),
		})
		return
	}

	if err := h.service.CreateUsers(r.Context(), *payload); err != nil {
		exception.WriteBadRequest(w, err.Error())
		return
	}

	if err := exception.SuccesWriteJson(w, "sukses buat akun"); err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}

}

func (h *Route) handleLogin(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	w.Header().Add("Content-Type", "application/json")
	payload := &web.UsersLoginPayload{}
	if err := exception.ParseJson(r, payload); err != nil {
		exception.WriteBadRequest(w, err.Error())
		return
	}

	u, err := h.service.GetEmail(r.Context(), payload.Email)
	if err != nil {
		exception.WriteBadRequest(w, err.Error())
		return
	}

	//check if the password was correct
	if err := auth.CheckHashedPassword(u.Password, []byte(payload.Password)); err != nil {
		exception.WriteBadRequest(w, exception.ErrIncorrectPassword.Error())
		return
	}

	ses, _ := auth.SessionStore.Get(r, "lg-ses")
	if err := auth.SessionSave(r, w, ses, u.Id); err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}

	if err := exception.SuccesWriteJson(w, "sukses login"); err != nil {
		exception.WriteInternalError(w, err.Error())
		return
	}

}

func (h *Route)handleLogout(w http.ResponseWriter, r *http.Request,params httprouter.Params)  {
	ses,_:=auth.SessionStore.Get(r,"lg-ses")
	ses.Options.MaxAge=-1
	if err:=ses.Save(r,w);err!=nil{
		exception.WriteInternalError(w, err.Error())
		return
	}
}