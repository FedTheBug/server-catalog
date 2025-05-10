package middleware

import (
	"github.com/server-catalog/internal/config"
	"github.com/server-catalog/internal/utils"
	"net/http"
)

func AppKeyResolver(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		appKey := r.Header.Get("App-key")
		if appKey == "" {
			_ = (&utils.Response{
				Status:  http.StatusBadRequest,
				Message: "App-key missing in header",
				Error:   "App-key missing in header",
			}).Render(w)
			return
		}

		if appKey != config.App().SecretKey {
			_ = (&utils.Response{
				Status:  http.StatusUnauthorized,
				Message: "Invalid App-key",
				Error:   "Invalid App-key",
			}).Render(w)
			return
		}

		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
