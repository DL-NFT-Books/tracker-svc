package handlers

import (
	"context"
	"net/http"

	booker "github.com/dl-nft-books/book-svc/connector"
	"github.com/dl-nft-books/tracker-svc/internal/data"
	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	dbCtxKey
	bookerCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxDB(db data.DB) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, dbCtxKey, db)
	}
}

func DB(r *http.Request) data.DB {
	return r.Context().Value(dbCtxKey).(data.DB).New()
}

func CtxBooker(entry *booker.Connector) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, bookerCtxKey, entry)
	}
}

func Booker(r *http.Request) *booker.Connector {
	return r.Context().Value(bookerCtxKey).(*booker.Connector)
}
