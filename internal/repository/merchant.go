package repository

import (
	"bytes"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/model"
	"github.com/nozzlium/belimang/internal/util"
)

type MerchantRepository struct {
	db *pgx.Conn
}

func NewMerchantRepository(
	db *pgx.Conn,
) *MerchantRepository {
	return &MerchantRepository{db: db}
}

func (r *MerchantRepository) Insert(
	ctx context.Context,
	merchant model.Merchant,
) (model.Merchant, error) {
	query := `
    insert into
    merchants (
      id,
      user_id,
      name,
      merchant_category,
      image_url,
      latitude,
      longitude,
      created_at
    ) values (
      $1, $2, $3, $4, $5, $6, $7, $8
    );
  `
	_, err := r.db.Exec(ctx, query,
		merchant.ID,
		merchant.UserID,
		merchant.Name,
		merchant.MerchantCategory,
		merchant.ImageURL,
		merchant.Latitude,
		merchant.Longitude,
		merchant.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return merchant, constant.ErrNotFound
			}
		}
		return merchant, err
	}

	return merchant, nil
}

func (r *MerchantRepository) FindAll(
	ctx context.Context,
	merchantQueries model.MerchantQueries,
) ([]model.Merchant, int, error) {
	var query bytes.Buffer
	query.WriteString(`
    select
      id,
      name, 
      merchant_category,
      image_url,
      latitude,
      longitude,
      created_at
    from merchants
    where 1 = 1
    `)
	queries, params := util.BuildQueryStringAndParams(
		&query,
		merchantQueries.BuildWhereClauses,
		merchantQueries.BuildPagination,
		merchantQueries.BuildOrderByClause,
		false,
	)

	var queryTotal bytes.Buffer
	queryTotal.WriteString(`
    select 
    count(id)
    from merchants
    where 1 = 1
    `)
	queryTotalString, paramsTotal := util.BuildQueryStringAndParamsWithoutLimit(
		&queryTotal,
		merchantQueries.BuildWhereClauses,
		nil,
	)

	batch := &pgx.Batch{}
	batch.Queue(queries, params...)
	batch.Queue(
		queryTotalString,
		paramsTotal...)

	br := r.db.SendBatch(ctx, batch)

	rows, err := br.Query()
	if err != nil {
		if !errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil, 0, err
		}
	}
	defer rows.Close()

	merchants := make(
		[]model.Merchant,
		0,
		merchantQueries.Limit,
	)
	for rows.Next() {
		var merchant model.Merchant
		rows.Scan(
			&merchant.ID,
			&merchant.Name,
			&merchant.MerchantCategory,
			&merchant.ImageURL,
			&merchant.Latitude,
			&merchant.Longitude,
			&merchant.CreatedAt,
		)
		merchants = append(
			merchants,
			merchant,
		)
	}

	var total int
	err = br.QueryRow().Scan(&total)
	if err != nil {
		if !errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil, 0, err
		}
	}
	defer br.Close()

	return merchants, total, nil
}
