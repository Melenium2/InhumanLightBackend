package apiserver

import (
	"net/http"
	"time"

	"github.com/inhumanLightBackend/app/apiserver/handlers"
	"github.com/inhumanLightBackend/app/store"
	"github.com/sirupsen/logrus"
)

// Init new server
func NewServer(store store.Store, config *Config) *http.Server {
	l := logrus.New()
	l.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp: true,
	})
	h := handlers.New(store, l)
	h.SetupRoutes()
	s := &http.Server{
		Addr: config.Port,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
		Handler: h,
	}

	return s
}
