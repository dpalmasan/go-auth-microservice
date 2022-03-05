package health

import (
	"net/http"

	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func init() {
	log.Formatter = new(logrus.JSONFormatter)
}

func Routes() chi.Router {
	router := chi.NewRouter()

	router.Use(chiMiddleware.AllowContentType("application/json"))

	router.Get("/", HealthStatus)

	return router
}

func HealthStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	parent := opentracing.GlobalTracer().StartSpan("GET /health")
	defer parent.Finish()

	w.Write([]byte(`{
		"healthy": "true"
	}`))
}
