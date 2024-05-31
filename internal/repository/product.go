package repository

import (
	"bytes"
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/model"
	"github.com/nozzlium/belimang/internal/util"
)

type ProductRepository struct {
	db *pgx.Conn
}

func NewProductRepository(
	db *pgx.Conn,
) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

func (r *ProductRepository) Insert(
	ctx context.Context,
	product model.Product,
) error {
	log.Println(product)
	query := `
    insert into
    products (
      id,
      user_id,
      merchant_id,
      name,
      price,
      product_category,
      image_url,
      created_at
    ) values (
      $1, $2, $3, $4, $5, $6, $7, $8
    )
  `
	_, err := r.db.Exec(ctx, query,
		product.ID,
		product.UserID,
		product.MerchantID,
		product.Name,
		product.Price,
		product.ProductCategory,
		product.ImageURL,
		product.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Println(err.Error())
			if pgErr.Code == "23503" {
				return constant.ErrNotFound
			}
		}
		return err
	}

	return nil
}

func (r *ProductRepository) FindAll(
	ctx context.Context,
	queries model.ProductQueries,
) ([]model.Product, int, error) {
	var queryItems bytes.Buffer
	queryItems.WriteString(`
    select
      id,
      name,
      product_category,
      price,
      image_url,
      created_at
    from products
    where 1 = 1
    `)
	queryItemsString, queryItemsParams := util.BuildQueryStringAndParams(
		&queryItems,
		queries.BuildWhereClauses,
		queries.BuildPagination,
		queries.BuildOrderByClause,
		false,
	)

	var queryTotal bytes.Buffer
	queryTotal.WriteString(`
    select
      count(id)
    from products
    where 1 = 1
    `)
	queryTotalString, queryTotalParams := util.BuildQueryStringAndParamsWithoutLimit(
		&queryTotal,
		queries.BuildWhereClauses,
		nil,
	)

	batch := &pgx.Batch{}
	batch.Queue(
		queryItemsString,
		queryItemsParams...)
	batch.Queue(
		queryTotalString,
		queryTotalParams...)

	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

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

	products := make(
		[]model.Product,
		0,
		queries.Limit,
	)
	for rows.Next() {
		var product model.Product
		rows.Scan(
			&product.ID,
			&product.Name,
			&product.ProductCategory,
			&product.Price,
			&product.ImageURL,
			&product.CreatedAt,
		)
		products = append(
			products,
			product,
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

	return products, total, nil
}
