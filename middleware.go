package vega

import (
	"net/http"
	"strconv"

	"github.com/justinas/nosurf"
)

// SessionLoad loads and saves session data for requests
func (v *Vega) SessionLoad(next http.Handler) http.Handler {
	return v.Session.LoadAndSave(next)
}

// NoSurf adds CSRF protection to all POST requests
func (v *Vega) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(v.config.cookie.secure)

	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   v.config.cookie.domain,
	})

	return csrfHandler
}

