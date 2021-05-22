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