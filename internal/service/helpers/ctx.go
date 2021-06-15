package helpers

import (
    "context"
    "github.com/redcuckoo/bsc-checker-events/internal/data"
    "net/http"

    "gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
    logCtxKey ctxKey = iota
    missionCtxKey
    explorerMissionCtxKey
    explorerCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
    return func(ctx context.Context) context.Context {
        return context.WithValue(ctx, logCtxKey, entry)
    }
}

func Log(r *http.Request) *logan.Entry {
    return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxMission(q data.MissionQ) func(context.Context) context.Context {
    return func(ctx context.Context) context.Context {
        return context.WithValue(ctx, missionCtxKey, q)
    }
}

func Mission(r *http.Request) data.MissionQ {
    return r.Context().Value(missionCtxKey).(data.MissionQ).New()
}

func CtxExplorerMission(q data.ExplorerMissionQ) func(context.Context) context.Context {
    return func(ctx context.Context) context.Context {
        return context.WithValue(ctx, explorerMissionCtxKey, q)
    }
}

func ExplorerMission(r *http.Request) data.ExplorerMissionQ {
    return r.Context().Value(explorerMissionCtxKey).(data.ExplorerMissionQ).New()
}

func CtxExplorer(q data.ExplorerQ) func(context.Context) context.Context {
    return func(ctx context.Context) context.Context {
        return context.WithValue(ctx, explorerCtxKey, q)
    }
}

func Explorer(r *http.Request) data.ExplorerQ {
    return r.Context().Value(explorerCtxKey).(data.ExplorerQ).New()
}