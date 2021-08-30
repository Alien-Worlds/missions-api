package cli

import (
    "context"
    "github.com/Alien-Worlds/missions-api/internal/config"
    "github.com/Alien-Worlds/missions-api/internal/service"
    "github.com/Alien-Worlds/missions-api/internal/service/checker-svc/checker"
    "github.com/alecthomas/kingpin"
    "gitlab.com/distributed_lab/kit/kv"
    "gitlab.com/distributed_lab/logan/v3"
)

func Run(args []string) bool {
    log := logan.New()

    defer func() {
        if rvr := recover(); rvr != nil {
            log.WithRecover(rvr).Error("app panicked")
        }
    }()

    cfg := config.New(kv.MustFromEnv())
    log = cfg.Log()

    app := kingpin.New("missions-api", "")

    runCmd := app.Command("run", "run command")
    serviceCmd := runCmd.Command("service", "run service") // you can insert custom help
    apiCmd := runCmd.Command("api", "run api")

    migrateCmd := app.Command("migrate", "migrate command")
    migrateUpCmd := migrateCmd.Command("up", "migrate db up")
    migrateDownCmd := migrateCmd.Command("down", "migrate db down")

    // custom commands go here...

    cmd, err := app.Parse(args[1:])
    if err != nil {
        log.WithError(err).Error("failed to parse arguments")
        return false
    }

    ctx := context.Background()

    switch cmd {
    case apiCmd.FullCommand():
        service.Run(cfg)
        return true
    case serviceCmd.FullCommand():
        svc := checker.New(cfg)
        svc.Run(ctx)
        return true
    case migrateUpCmd.FullCommand():
        err = MigrateUp(cfg)
    case migrateDownCmd.FullCommand():
        err = MigrateDown(cfg)
    default:
        log.Errorf("unknown command %s", cmd)
        return false
    }

    if err != nil {
        log.WithError(err).Error("failed to exec cmd")
        return false
    }

    //<- ctx.Done()
    return true
}
