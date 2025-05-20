package api

import (
    "fmt"
    "net/http"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
    HttpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Общее количество HTTP-запросов",
        },
        []string{"path", "method", "status"},
    )

    UserLoginTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "user_login_total",
            Help: "Total number of user logins",
        },
        []string{"method", "status"},
    )
)

func InitMetrics() {
    prometheus.MustRegister(HttpRequestsTotal)
    prometheus.MustRegister(UserLoginTotal)
}

// Middleware для HTTP метрик
func MetricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rw := &responseWriter{w, http.StatusOK}
        next.ServeHTTP(rw, r)
        HttpRequestsTotal.WithLabelValues(r.URL.Path, r.Method, fmt.Sprint(rw.statusCode)).Inc()
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func MetricsHandler() http.Handler {
    return promhttp.Handler()
}
