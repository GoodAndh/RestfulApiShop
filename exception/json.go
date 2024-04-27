package exception

import (
	"encoding/json"
	"net/http"
	"restful/model/web"
)

func ParseJson(r *http.Request,v any)error {
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJson(w http.ResponseWriter,code int,status,message string,data...any)error  {
	web:=&web.WebResponse{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}
	return json.NewEncoder(w).Encode(web)
}

func SuccesWriteJson(w http.ResponseWriter,message string,data...any)error  {
	return WriteJson(w,http.StatusOK,"succes",message,data...)
}

func WriteInternalError(w http.ResponseWriter,message string,data...any)error  {
	return WriteJson(w,http.StatusInternalServerError,"internal server error",message,data...)
}
func WriteBadRequest(w http.ResponseWriter,message string,data...any)error  {
	return WriteJson(w,http.StatusBadRequest,"bad request",message,data...)
}