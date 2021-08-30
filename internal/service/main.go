package service

import (
	"net"
	"net/http"

	"github.com/Alien-Worlds/missions-api/internal/config"
	"gitlab.com/distributed_lab/logan/v3"
)

type service struct {
	log      *logan.Entry
	listener net.Listener
}

func (s *service) run(cfg config.Config) error {
	s.log.Info("Running api service")

	r := s.router(cfg)

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {
	log := cfg.Log().WithField("api-service", "")

	return &service{
		log:      log,
		listener: cfg.Listener(),
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(cfg); err != nil {
		panic(err)
	}
}
