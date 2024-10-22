package middleware

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type Onion struct {
	Middlewares []func(next httprouter.Handle) httprouter.Handle
	logger      *zap.Logger
}

func NewOnion() *Onion {
	logger, _ := zap.NewDevelopment()
	return &Onion{
		logger: logger,
	}
}

func (o *Onion) Apply(h httprouter.Handle) httprouter.Handle {
	for i := range o.Middlewares {
		h = o.Middlewares[i](h)
	}
	return h
}

func (o *Onion) AppendMiddleware(mw func(next httprouter.Handle) httprouter.Handle) {
	o.Middlewares = append(o.Middlewares, mw)
}

func (o *Onion) LogRequestResponse(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		o.logger.Info("Request recieved",
			zap.String("method", r.Method),
			zap.String("requestURI", r.RequestURI),
			zap.String("host", r.Host),
		)
		start := time.Now()

		next(w, r, ps)

		o.logger.Info("Reposne sent",
			zap.Duration("time to handle", time.Since(start)),
		)
	}
}
