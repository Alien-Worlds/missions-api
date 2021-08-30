package cli

import (
    "github.com/Alien-Worlds/missions-api/internal/assets"
    "github.com/Alien-Worlds/missions-api/internal/config"
    migrate "github.com/rubenv/sql-migrate"
    "gitlab.com/distributed_lab/logan/v3/errors"
)

var migrations = &migrate.PackrMigrationSource{
    Box: assets.Migrations,
}

func MigrateUp(cfg config.Config) error {
    applied, err := migrate.Exec(cfg.DB().RawDB(), "postgres", migrations, migrate.Up)
    if err != nil {
      return errors.Wrap(err, "failed to apply migrations")
    }
    cfg.Log().WithField("applied", applied).Info("migrations applied")
    return nil
}

func MigrateDown(cfg config.Config) error {
    applied, err := migrate.Exec(cfg.DB().RawDB(), "postgres", migrations, migrate.Down)
    if err != nil {
      return errors.Wrap(err, "failed to apply migrations")
    }
    cfg.Log().WithField("applied", applied).Info("migrations applied")
    return nil
}
