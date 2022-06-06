package product

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/RoseRocket/xerrs"
	"github.com/jmoiron/sqlx"

	"github.com/ndv6/rm-api/api/respond"
	"github.com/ndv6/rm-api/proto/model"
)

const (
	SqlProduct          = "products"
	SqlProductCategory  = "product_categories"
	SqlCurrency         = "currencies"
	SqlProductFacility  = "product_facilities"
	SqlSegments         = "segments"
	SqlTransactionBonds = "transaction_bonds"
)

type SqlStorage struct {
	dbr *sqlx.DB
}

func NewProductSqlStorage(dbr *sqlx.DB) Storage {
	return &SqlStorage{
		dbr: dbr,
	}
}

func (prd *SqlStorage) reader() sqlx.QueryerContext {
	return prd.dbr
}

func (prd *SqlStorage) GetFacilities(ctx context.Context) ([]*model.Facility, error) {
	var (
		q                                  string
		err                                error
		id                                 int64
		includes, excludes                 FacilityFeatures
		productId, customerType, segmentId uint64
		result                             []*model.Facility
	)
	q = fmt.Sprintf(`
	SELECT id, customer_type, segment_id, product_id, includes, excludes
	FROM %s`, SqlProductFacility)

	rows, err := prd.reader().QueryxContext(ctx, q)
	if err != nil {
		err = xerrs.Mask(err, respond.ErrQueryRead)
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		err = rows.Scan(&id, &customerType, &segmentId, &productId, &includes, &excludes)
		if err != nil {
			err = xerrs.Mask(err, respond.ErrQueryRead)
			return nil, err
		}
		result = append(result, &model.Facility{
			CustomerType: customerType,
			SegmentId:    segmentId,
			ProductId:    productId,
			Includes:     includes,
			Excludes:     excludes,
		})
	}
	return result, nil
}

func (prd *SqlStorage) GetCurrenciesInit(ctx context.Context) ([]*model.Currency, error) {
	var (
		q, title string
		err      error
		result   []*model.Currency
		symb     sql.NullString
		id       uint64
	)
	q = fmt.Sprintf(`
	SELECT id, code, title
	FROM %s where disabled = false`, SqlCurrency)

	rows, err := prd.reader().QueryxContext(ctx, q)
	if err != nil {
		err = xerrs.Mask(err, respond.ErrQueryRead)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		err = rows.Scan(&id, &symb, &title)
		if err != nil {
			err = xerrs.Mask(err, respond.ErrQueryRead)
			return nil, err
		}
		result = append(result, &model.Currency{
			Id:     id,
			Symbol: symb.String,
			Title:  title,
		})
	}

	return result, nil
}

func (prd *SqlStorage) GetProductsInit(ctx context.Context) ([]*Product, error) {
	var (
		q, title          string
		err               error
		id, catID, currID uint64
		multiCurr         bool
		result            []*Product
	)

	gU := getClaims(ctx)

	if gU.Private.GroupName == "TELESALES" || gU.Private.GroupName == "MSA" {
		q = fmt.Sprintf(`
	SELECT 
		p.id, p.title, p.category_id, p.currency_id,
		p.multi_currency, p.deposit_map
	FROM %s p WHERE disabled = false AND category_id != 3 ORDER BY p.id
	`, SqlProduct) //MSA cannot access Tanda Junior
	} else {
		q = fmt.Sprintf(`
	SELECT 
		p.id, p.title, p.category_id, p.currency_id,
		p.multi_currency, p.deposit_map
	FROM %s p WHERE disabled = false ORDER BY p.id
	`, SqlProduct)
	}

	rows, err := prd.reader().QueryxContext(ctx, q)
	if err != nil {
		err = xerrs.Mask(err, respond.ErrQueryRead)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var rate RateInfo
		err = rows.Scan(&id, &title, &catID, &currID, &multiCurr, &rate)
		if err != nil {
			err = xerrs.Mask(err, respond.ErrQueryRead)
			return nil, err
		}
		result = append(result, &Product{
			Product: model.Product{
				Id:            id,
				Title:         title,
				CategoryId:    catID,
				CurrencyId:    currID,
				MultiCurrency: multiCurr,
			},
			Tenor: struct {
				Month int    `json:"month"`
				Title string `json:"title"`
			}{
				rate.Tenor[0],
				strconv.Itoa(rate.Tenor[0]) + " bulan",
			},
		})
	}
	return result, nil
}

func (prd *SqlStorage) GetProductCategoriesInit(ctx context.Context) ([]*model.ProductCategory, error) {
	var (
		q        string
		err      error
		id       uint64
		title    string
		parentID uint64
		result   []*model.ProductCategory
	)
	gU := getClaims(ctx)

	if gU.Private.GroupName == "TELESALES" || gU.Private.GroupName == "MSA" {
		q = fmt.Sprintf(`
	SELECT id, title, parent_id
	FROM %s pc WHERE id != 3
	`, SqlProductCategory) //MSA cannot access Tanda Junior
	} else {
		q = fmt.Sprintf(`
	SELECT id, title, parent_id
	FROM %s pc
	`, SqlProductCategory)
	}

	rows, err := prd.reader().QueryxContext(ctx, q)
	if err != nil {
		err = xerrs.Mask(err, respond.ErrQueryRead)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		err = rows.Scan(&id, &title, &parentID)
		if err != nil {
			err = xerrs.Mask(err, respond.ErrQueryRead)
			return nil, err
		}
		result = append(result, &model.ProductCategory{
			Id:       id,
			Title:    title,
			ParentId: parentID,
		})
	}
	return result, nil
}

func (prd *SqlStorage) GetSegments(ctx context.Context) ([]*model.Segment, error) {
	var (
		rows *sqlx.Rows
		err  error
	)
	query := fmt.Sprintf(`select id, code, title, product_ids from %s`, SqlSegments)
	rows, err = prd.reader().QueryxContext(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return nil, nil
		}
		err = xerrs.Mask(err, respond.ErrQueryRead)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	var (
		items       []*model.Segment
		id          uint64
		code, title string
		productIds  Ids
	)
	for rows.Next() {
		err = rows.Scan(&id, &code, &title, &productIds)
		if err != nil {
			err = xerrs.Mask(err, respond.ErrQueryRead)
			return nil, err
		}
		items = append(items, &model.Segment{
			Id:         id,
			Code:       code,
			Title:      title,
			Content:    "",
			ProductIds: productIds,
		})
	}
	return items, nil
}

func (prd *SqlStorage) GetProductsByIds(ctx context.Context, ids []uint64) (map[uint64]*model.Product, error) {
	var (
		q, title          string
		err               error
		id, catID, currID uint64
		multiCurr         bool
	)

	gU := getClaims(ctx)

	if gU.Private.GroupName == "TELESALES" || gU.Private.GroupName == "MSA" {
		q = fmt.Sprintf(`
	SELECT 
		p.id, p.title, p.category_id, p.currency_id,
		p.multi_currency
	FROM %s p WHERE disabled = false AND category_id != 3
	`, SqlProduct) //MSA cannot access Tanda Junior
	} else {
		q = fmt.Sprintf(`
	SELECT 
		p.id, p.title, p.category_id, p.currency_id,
		p.multi_currency
	FROM %s p WHERE disabled = false
	`, SqlProduct)
	}

	var idstrings []string
	if len(ids) > 0 {
		for _, v := range ids {
			idstrings = append(idstrings, strconv.Itoa(int(v)))
		}
		q += fmt.Sprintf(` AND id IN (%s)`, strings.Join(idstrings, ","))
	}
	rows, err := prd.reader().QueryxContext(ctx, q+" ORDER BY id")
	if err != nil {
		err = xerrs.Mask(err, respond.ErrQueryRead)
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	result := map[uint64]*model.Product{}
	for rows.Next() {
		err = rows.Scan(&id, &title, &catID, &currID, &multiCurr)
		if err != nil {
			err = xerrs.Mask(err, respond.ErrQueryRead)
			return nil, err
		}
		result[id] = &model.Product{
			Id:            id,
			Title:         title,
			CategoryId:    catID,
			CurrencyId:    currID,
			MultiCurrency: multiCurr,
		}
	}
	return result, nil
}

func (prd *SqlStorage) GetProductById(ctx context.Context, u uint64) (ProductDetail, error) {
	dP := ProductDetail{}
	q := fmt.Sprintf("SELECT title, code FROM %s products WHERE disabled = false AND id = $1=?", SqlProduct)
	err := prd.dbr.QueryRowxContext(ctx, q, u).Scan(&dP.Title, &dP.Code)
	if err == nil {
		return dP, err
	}

	return dP, err
}

func (prd *SqlStorage) GetProductObligasiById(ctx context.Context, u uint64) (ObligasiItems, string, error) {
	oI := ObligasiItems{}
	var userBranch string
	q := fmt.Sprintf("SELECT security_no, deal_id, jenis_transaksi_bank, product_name, user_branch FROM %s WHERE application_id = $1", SqlTransactionBonds)

	err := prd.dbr.QueryRowxContext(ctx, q, u).Scan(&oI.SecurityNo, &oI.DealId, &oI.TransactionType, &oI.ProductName, &userBranch)
	if err == nil {
		return oI, userBranch, err
	}

	return oI, userBranch, err
}
