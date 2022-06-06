package product

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/ndv6/rm-api/proto/model"
	"github.com/ndv6/rm-api/proto/service"
)

type Images []string

func (pi *Images) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var slc []string
	err := json.Unmarshal(b, &slc)
	if err != nil {
		return err
	}
	*pi = slc
	return nil
}

func (pi Images) Value() (driver.Value, error) {
	return json.Marshal(pi)
}

type Ids []uint64

func (i *Ids) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var slc []uint64
	err := json.Unmarshal(b, &slc)
	if err != nil {
		return err
	}
	*i = slc
	return nil
}

func (i Ids) Value() (driver.Value, error) {
	return json.Marshal(i)
}

type FacilityFeatures []model.FacilityFeature

func (ff *FacilityFeatures) Scan(v interface{}) error {
	b, ok := v.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	var slc []model.FacilityFeature
	err := json.Unmarshal(b, &slc)
	if err != nil {
		return err
	}
	*ff = slc
	return nil
}

func (ff FacilityFeatures) Value() (driver.Value, error) {
	return json.Marshal(ff)
}

type RateInfo struct {
	MinTarget uint64    `json:"min_target"`
	Mod       uint64    `json:"mod"` // target dana harus merupakan kelipatan mod
	Tenor     [2]int    `json:"tenor"`
	Amount    [2]uint64 `json:"amount"`
}

// Value() used by sql to store this data into jsonb. See driver.Valuer interface.
func (ri RateInfo) Value() (driver.Value, error) {
	return json.Marshal(ri)
}

// Scan() used by sql to convert jsonb data into this type. See sql.Scanner interface.
func (ri *RateInfo) Scan(src interface{}) error {
	v, ok := src.([]byte)
	if !ok {
		return errors.New("unknown data type")
	}

	return json.Unmarshal(v, ri)
}

type Product struct {
	model.Product
	Tenor struct {
		Month int    `json:"month"`
		Title string `json:"title"`
	} `json:"tenor,omitempty"`
}

type ResponseProduct struct {
	service.ResponseProductAndCategoryAndCurrency
	Products        []*Product    `json:"products"`
	DepositoRenewal []RenewalInfo `json:"deposito_renewal"`
}

type RenewalInfo struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type ProductDetail struct {
	ID      uint64 `json:"product_id"`
	Title   string `json:"product_title,omitempty"`
	Code    string `json:"product_code"`
	Type    int    `json:"product_type"`
	Context context.Context
	Error   error `json:"error"`
}

type ObligasiItems struct {
	SecurityNo      string `json:"security_no"`
	DealId          string `json:"deal_id"`
	ProductName     string `json:"product_name"`
	TransactionType string `json:"transaction_type"`
}

func (p *ProductDetail) GetProductDetail() *ProductDetail {
	productStorage := getProductStorage(p.Context)
	product, err := productStorage.GetProductById(p.Context, p.ID)
	if err != nil {
		return &ProductDetail{Error: err}
	}

	return &product
}

func (p *ProductDetail) GetProductName() string {
	return p.Title
}
