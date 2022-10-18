package helpers

import (
	"context"
	"github.com/Digital-Voting-Team/staff-service/internal/data"

	"net/http"

	"gitlab.com/distributed_lab/logan/v3"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	addressesQCtxKey
	personsQCtxKey
	staffQCtxKey
	positionQCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxAddressesQ(entry data.AddressesQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, addressesQCtxKey, entry)
	}
}

func AddressesQ(r *http.Request) data.AddressesQ {
	return r.Context().Value(addressesQCtxKey).(data.AddressesQ).New()
}

func CtxPersonsQ(entry data.PersonsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, personsQCtxKey, entry)
	}
}

func PersonsQ(r *http.Request) data.PersonsQ {
	return r.Context().Value(personsQCtxKey).(data.PersonsQ).New()
}

func CtxStaffQ(entry data.StaffQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, staffQCtxKey, entry)
	}
}

func StaffQ(r *http.Request) data.StaffQ {
	return r.Context().Value(staffQCtxKey).(data.StaffQ).New()
}

func CtxPositionsQ(entry data.PositionsQ) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, positionQCtxKey, entry)
	}
}

func PositionsQ(r *http.Request) data.PositionsQ {
	return r.Context().Value(positionQCtxKey).(data.PositionsQ).New()
}
