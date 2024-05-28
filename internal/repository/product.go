package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/model"
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
