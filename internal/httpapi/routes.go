package httpapi

import "net/http"

func Routes(handler *Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/gesture", handler.Gesture)
	mux.HandleFunc("/phrase", handler.CustomPhrase)
	mux.HandleFunc("/end", handler.End)
	mux.HandleFunc("/snapshot", handler.Snapshot)
	mux.HandleFunc("/panel", handler.Panel)
	return requestLogger(mux)
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Showroom-Route", r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
