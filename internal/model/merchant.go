package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/util"
)

type (
	MerchantCategory string
	OrderBy          string
)

const (
	Asc  OrderBy = "asc"
	Desc OrderBy = "desc"
)

func (o OrderBy) IsValid() bool {
	switch o {
	case Asc, Desc:
		return true
	default:
		return false
	}
}

const (
	SmallRestaurant       MerchantCategory = "SmallRestaurant"
	MediumRestaurant      MerchantCategory = "MediumRestaurant"
	LargeRestaurant       MerchantCategory = "LargeRestaurant"
	MerchandiseRestaurant MerchantCategory = "MerchandiseRestaurant"
	BoothKiosk            MerchantCategory = "BoothKiosk"
	ConvenienceStore      MerchantCategory = "ConvenienceStore"
)

type Merchant struct {
	ID               uuid.UUID
	UserID           uuid.UUID
	Name             string
	MerchantCategory MerchantCategory
	ImageURL         string
	Latitude         float64
	Longitude        float64
	CreatedAt        time.Time
}

type MerchantRequestBody struct {
	Name             string                      `json:"name"`
	MerchantCategory MerchantCategory            `json:"merchantCategory"`
	ImageURL         string                      `json:"imageUrl"`
	Location         MerchantLocationRequestBody `json:"location"`
}

type MerchantLocationRequestBody struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func (body MerchantRequestBody) IsValid() (Merchant, error) {
	var merchant Merchant
	if nameLen := len(body.Name); nameLen < 2 ||
		nameLen > 30 {
		return merchant, constant.ErrBadInput
	}
	merchant.Name = body.Name

	switch body.MerchantCategory {
	case SmallRestaurant,
		MediumRestaurant,
		LargeRestaurant,
		MerchandiseRestaurant,
		BoothKiosk,
		ConvenienceStore:
		merchant.MerchantCategory = body.MerchantCategory
	default:
		return merchant, constant.ErrBadInput
	}

	err := util.ValidateURL(
		body.ImageURL,
	)
	if err != nil {
		return Merchant{}, err
	}
	merchant.ImageURL = body.ImageURL

	merchant.Latitude = body.Location.Lat
	merchant.Longitude = body.Location.Long

	return merchant, nil
}

type MerchantQueries struct {
	MerchantID       string           `query:"merchantId"`
	Name             string           `query:"name"`
	MerchantCategory MerchantCategory `query:"merchantCategory"`
	Limit            int
	Offset           int
	CreatedAt        string
}

func (q *MerchantQueries) BuildWhereClauses() ([]string, []interface{}) {
	clauses := make([]string, 0, 4)
	params := make([]interface{}, 0, 4)

	merchantId, err := uuid.Parse(
		q.MerchantID,
	)
	if err == nil {
		clauses = append(
			clauses,
			"id = $%d",
		)
		params = append(
			params,
			merchantId,
		)
	}

	if q.Name != "" {
		clauses = append(
			clauses,
			"name ilike '%%' || $%d || '%%'",
		)
		params = append(params, q.Name)
	}

	switch q.MerchantCategory {
	case SmallRestaurant,
		MediumRestaurant,
		LargeRestaurant,
		MerchandiseRestaurant,
		BoothKiosk,
		ConvenienceStore:
		clauses = append(
			clauses,
			"merchant_category = $%d",
		)
		params = append(
			params,
			q.MerchantCategory,
		)
	}

	return clauses, params
}

func (q *MerchantQueries) BuildPagination() (string, []interface{}) {
	var params []interface{}

	limit := 5
	offset := 0
	if q.Limit > 0 {
		limit = q.Limit
	}
	if q.Offset > 0 {
		offset = q.Offset
	}
	params = append(
		params,
		limit,
		offset,
	)

	return "limit $%d offset $%d", params
}

func (q *MerchantQueries) BuildOrderByClause() []string {
	var sqlClause []string

	if q.CreatedAt != "" ||
		OrderBy(
			q.CreatedAt,
		).IsValid() {
		sqlClause = append(
			sqlClause,
			fmt.Sprintf(
				"created_at %s",
				q.CreatedAt,
			),
		)
	} else {
		sqlClause = append(
			sqlClause,
			"created_at desc",
		)
	}

	return sqlClause
}

type MerchantResponaeBody struct {
	MerchantID       string               `json:"merchantId"`
	Name             string               `json:"name"`
	MerchantCategory string               `json:"merchantCategory"`
	ImageURL         string               `json:"imageUrl"`
	Location         LocationResponseBody `json:"location"`
	CreatedAt        string               `json:"createdat"`
}

type LocationResponseBody struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"Long"`
}

func (m *Merchant) ToResponseBody() MerchantResponaeBody {
	return MerchantResponaeBody{
		MerchantID: m.ID.String(),
		Name:       m.Name,
		MerchantCategory: string(
			m.MerchantCategory,
		),
		ImageURL: m.ImageURL,
		Location: LocationResponseBody{
			Lat:  m.Latitude,
			Long: m.Longitude,
		},
		CreatedAt: util.ToISO8601(
			m.CreatedAt,
		),
	}
}
