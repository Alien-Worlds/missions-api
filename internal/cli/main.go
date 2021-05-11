package cli

import (
    "context"
    "github.com/alecthomas/kingpin"
    "github.com/redcuckoo/bsc-checker-events/internal/config"
    "github.com/redcuckoo/bsc-checker-events/internal/service/checker-svc/checker"
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

    app := kingpin.New("bsc-checker-events", "")

    runCmd := app.Command("run", "run command")
    serviceCmd := runCmd.Command("service", "run service") // you can insert custom help

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
    case serviceCmd.FullCommand():
        svc := checker.New(cfg)
        svc.Run(ctx)
        return true
    case migrateUpCmd.FullCommand():
        err = MigrateUp(cfg)
        //return true
    case migrateDownCmd.FullCommand():
        err = MigrateDown(cfg)
        //return true
    // handle any custom commands here in the same way
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
