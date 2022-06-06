package product

import (
	"context"
	"net/http"

	"github.com/ndv6/rm-api/api/jwe"

	"github.com/ndv6/rm-api/api/middleware"
)

const (
	StorageKey = "product_storage"
)

var storage Storage

func StorageContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if storage == nil {
			storage = NewProductSqlStorage(middleware.GetDbRead(r.Context()))
		}
		ctx := context.WithValue(r.Context(), StorageKey, storage)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getProductStorage(c context.Context) Storage {
	return c.Value(StorageKey).(Storage)
}

func getClaims(c context.Context) *jwe.Claims {
	if value := c.Value(jwe.ClaimsKey); value != nil {
		return value.(*jwe.Claims)
	}
	return &jwe.Claims{}
}
