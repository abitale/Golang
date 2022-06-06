package product

import (
	"context"

	"github.com/ndv6/rm-api/proto/model"
)

type Reader interface {
	GetFacilities(ctx context.Context) ([]*model.Facility, error)
	GetProductsInit(context.Context) ([]*Product, error)
	GetProductCategoriesInit(context.Context) ([]*model.ProductCategory, error)
	GetCurrenciesInit(context.Context) ([]*model.Currency, error)
	GetSegments(ctx context.Context) ([]*model.Segment, error)
	GetProductsByIds(context.Context, []uint64) (map[uint64]*model.Product, error)
	GetProductById(context.Context, uint64) (ProductDetail, error)
	GetProductObligasiById(context.Context, uint64) (ObligasiItems, string, error)
}

type Writer interface {
}

type Storage interface {
	Reader
	Writer
}
