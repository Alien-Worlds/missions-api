package service

import (
	"net"
	"net/http"

	"github.com/Alien-Worlds/missions-api/internal/config"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

type service struct {
	log      *logan.Entry
	copus    types.Copus
	listener net.Listener
}

func (s *service) run(cfg config.Config) error {
	s.log.Info("Running api service")

	r := s.router(cfg)

	if err := s.copus.RegisterChi(r); err != nil {
		s.log.Info("errored while running api service")

		return errors.Wrap(err, "cop failed")
	}

	return http.Serve(s.listener, r)
}

func newService(cfg config.Config) *service {
	log := cfg.Log().WithField("api-service", "")

	return &service{
		log:      log,
		copus:    cfg.Copus(),
		listener: cfg.Listener(),
	}
}

func Run(cfg config.Config) {
	if err := newService(cfg).run(cfg); err != nil {
		panic(err)
	}
}
