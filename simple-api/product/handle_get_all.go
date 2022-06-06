package product

import (
	"net/http"

	"github.com/ndv6/rm-api/api/respond"
	"github.com/ndv6/rm-api/proto/model"
	"github.com/ndv6/rm-api/proto/service"
)

func HandleGetAll(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		repo     Storage
		pCats    []*model.ProductCategory
		ps       []*Product
		fs       []*model.Facility
		currency []*model.Currency
		segments []*model.Segment
		err      error
	)
	repo = getProductStorage(ctx)
	ps, err = repo.GetProductsInit(ctx)
	if err != nil {
		respond.Failed(w, r, http.StatusInternalServerError, err)
		return
	}
	fs, err = repo.GetFacilities(ctx)
	if err != nil {
		respond.Failed(w, r, http.StatusInternalServerError, err)
		return
	}
	pCats, err = repo.GetProductCategoriesInit(ctx)
	if err != nil {
		respond.Failed(w, r, http.StatusInternalServerError, err)
		return
	}
	currency, err = repo.GetCurrenciesInit(ctx)
	if err != nil {
		respond.Failed(w, r, http.StatusInternalServerError, err)
		return
	}
	segments, err = repo.GetSegments(ctx)
	respond.Success(w, r, http.StatusOK, &ResponseProduct{
		ResponseProductAndCategoryAndCurrency: service.ResponseProductAndCategoryAndCurrency{
			Facilities: fs,
			Categories: pCats,
			Currencies: currency,
			Segments:   segments,
		},
		Products: ps,
		DepositoRenewal: []RenewalInfo{
			{
				Code:  "N",
				Title: "Tidak diperpanjang otomatis",
				Desc:  "Saat jatuh tempo, deposito tidak akan diperpanjang. Nilai Pokok dan Nilai Bunga akan dicairkan ke rekening Anda.",
			},
			{
				Code:  "A",
				Title: "Perpanjangan dan Bunga Ditambahkan ke Pokok",
				Desc:  "Saat jatuh tempo, deposito akan otomatis diperpanjang menggunakan Nilai Pokok dan Nilai Bunga, bersama dengan Suku Bunga yang berlaku.",
			},
			{
				Code:  "Y",
				Title: "Perpanjangan otomatis",
				Desc:  "Saat jatuh tempo, deposito akan otmatis diperpanjang menggunakan Nilai Pokok saja.",
			},
		},
	})
}
