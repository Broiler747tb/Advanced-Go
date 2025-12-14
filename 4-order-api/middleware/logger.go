package middleware

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logrus.SetOutput(os.Stdout)
		wrapper := &WrapperWriter{
			ResponseWriter: w,
		}
		next.ServeHTTP(wrapper, r)
		logrus.WithFields(logrus.Fields{
			"Method":     r.Method,
			"Path":       r.URL.Path,
			"Statuscode": wrapper.StatusCode,
		})
	})
}
