package log

import (
	"net/http"

	"github.com/konectdigital/drawbridge/log"
	"github.com/konectdigital/drawbridge/plugin"
	"github.com/konectdigital/muxinator"
)

func init() {
	plugin.RegisterPlugin("log", &Logger{})
}

type Logger struct{}

func (l *Logger) Middleware(map[string]interface{}) (muxinator.Middleware, error) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		log.Printf("%s %s", r.Method, r.RequestURI)
		next(w, r)
	}, nil
}
