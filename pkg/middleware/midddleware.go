package middleware

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Function func(next httprouter.Handle) httprouter.Handle

type Onion struct {
	middlewares []Function
	log         *zap.Logger
	duration    time.Duration
}

func NewOnion(lg *zap.Logger) *Onion {
	return &Onion{
		log: lg,
	}
}

func (o *Onion) Apply(h httprouter.Handle) httprouter.Handle {
	for i := range o.middlewares {
		h = o.middlewares[i](h)
	}
	return h
}

func (o *Onion) AppendMiddleware(mw ...Function) {
	o.middlewares = append(o.middlewares, mw...)
}

func (o *Onion) LogRequestResponse(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		o.log.Info("Request recieved",
			zap.String("method", r.Method),
			zap.String("requestURI", r.RequestURI),
			zap.String("host", r.Host),
			zap.Time("time", time.Now()),
		)

		next(w, r, ps)

		o.log.Info("Reposne sent",
			zap.Duration("time to handle", o.duration),
		)
	}
}

func (o *Onion) Timer(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		start := time.Now()
		next(w, r, ps)
		o.duration = time.Since(start)
	}
}
