package web

import (
	"net/http"
	"time"

	"github.com/chonglou/arche/web/mux"
	log "github.com/sirupsen/logrus"
)

type loggedResponse struct {
	http.ResponseWriter
	status int
}

func (l *loggedResponse) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}

// LoggerMiddleware logger middleware
func LoggerMiddleware(c *mux.Context) {
	log.Infof("%s %s %s", c.Request.Proto, c.Request.Method, c.Request.URL)
	begin := time.Now()
	lw := loggedResponse{ResponseWriter: c.Writer}
	c.Writer = &lw
	c.Next()
	log.Infof("%d %s", lw.status, time.Now().Sub(begin))
}
