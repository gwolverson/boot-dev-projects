package main

import (
	"html/template"
	"net/http"
)

var tmpl = template.Must(template.New("metrics").Parse(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited {{.}} times!</p>
  </body>
</html>
`))

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	err := tmpl.Execute(w, cfg.fileserverHits.Load())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
