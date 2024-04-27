package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
)

const rndomstring = "abcdefghijklmopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRRSTUVWXY"

var (
	SessionStore *sessions.CookieStore = sessions.NewCookieStore([]byte(rndomstring))
)

func SessionSave(r *http.Request, w http.ResponseWriter, session *sessions.Session, v any) error {
	session.Options.MaxAge=3600
	session.Values["auten"] = true
	session.Values["data"] = v

	return sessions.Save(r, w)
}
