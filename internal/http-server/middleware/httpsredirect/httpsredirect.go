package httpsredirect

import "net/http"

func RedirectToHTTPS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, если запрос пришел не через HTTPS
		if r.TLS == nil {
			// Перенаправляем на HTTPS
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
