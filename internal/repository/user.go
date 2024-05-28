package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/nozzlium/belimang/internal/constant"
	"github.com/nozzlium/belimang/internal/model"
)

type UserRepository struct {
	db *pgx.Conn
}

func NewUserRepository(
	db *pgx.Conn,
) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateAdmin(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return user, err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	queryInsertUser := `
    insert into
    users (
      id,
      username
    ) values (
      $1, $2
    );
  `
	batch.Queue(
		queryInsertUser,
		user.ID,
		user.Username,
	)

	queryInsertUsername := `
    insert into
    admin_details (
      user_id,
      email,
      password
    ) values (
      $1, $2, $3
    );
  `
	batch.Queue(
		queryInsertUsername,
		user.ID,
		user.Email,
		user.Password,
	)

	batchRes := tx.SendBatch(ctx, batch)
	if err := batchRes.Close(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return user, constant.ErrConflict
			}
		}
		return user, err
	}

	if err := tx.Commit(ctx); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return user, err
	}
	defer tx.Rollback(ctx)

	batch := &pgx.Batch{}

	queryInsertUser := `
    insert into
    users (
      id,
      username
    ) values (
      $1, $2
    );
  `
	batch.Queue(
		queryInsertUser,
		user.ID,
		user.Username,
	)

	queryInsertUsername := `
    insert into
    user_details (
      user_id,
      email,
      password
    ) values (
      $1, $2, $3
    );
  `
	batch.Queue(
		queryInsertUsername,
		user.ID,
		user.Email,
		user.Password,
	)

	batchRes := tx.SendBatch(ctx, batch)
	if err := batchRes.Close(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return user, constant.ErrConflict
			}
		}
		return user, err
	}

	if err := tx.Commit(ctx); err != nil {
		return user, err
	}

	return user, nil
}

func (r *UserRepository) FindAdminByUsername(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	queryFindUser := `
    select
      u.id,
      u.username,
      ad.email,
      ad.password
    from 
      users u 
      inner join admin_details ad on u.id = ad.user_id 
    where u.username = $1
  `

	err := r.db.QueryRow(
		ctx,
		queryFindUser,
		user.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return user, constant.ErrNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *UserRepository) FindUserByUsername(
	ctx context.Context,
	user model.User,
) (model.User, error) {
	queryFindUser := `
    select
      u.id,
      u.username,
      ud.email,
      ud.password
    from 
      users u 
      inner join user_details ud on u.id = ud.user_id 
    where u.username = $1
  `

	err := r.db.QueryRow(
		ctx,
		queryFindUser,
		user.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return user, constant.ErrNotFound
		}
		return user, err
	}

	return user, nil
}
