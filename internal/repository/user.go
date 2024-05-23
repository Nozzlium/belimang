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
    admins (
      id,
      email, 
      password
    ) values (
      $1, $2, $3
    );
  `
	batch.Queue(
		queryInsertUser,
		user.ID,
		user.Email,
		user.Password,
	)

	queryInsertUsername := `
    insert into
    admin_usernames (
      admin_id,
      username
    ) values (
      $1, $2
    );
  `
	batch.Queue(
		queryInsertUsername,
		user.ID,
		user.Username,
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
      email, 
      password
    ) values (
      $1, $2, $3
    );
  `
	batch.Queue(
		queryInsertUser,
		user.ID,
		user.Email,
		user.Password,
	)

	queryInsertUsername := `
    insert into
    user_usernames (
      user_id,
      username
    ) values (
      $1, $2
    );
  `
	batch.Queue(
		queryInsertUsername,
		user.ID,
		user.Username,
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
      a.id,
      u.username,
      a.email,
      a.password
    from 
      admins a 
      inner join admin_usernames u on u.admin_id = a.id 
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
      un.username,
      u.email,
      u.password
    from 
      users u 
      inner join user_usernames un on un.user_id = u.id 
    where un.username = $1
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

func (r *UserRepository) FindAdminUsername(
	ctx context.Context,
	username string,
) (string, error) {
	query := `
    select username
    from admin_usernames
    where username = $1
  `
	var savedUsername string
	err := r.db.QueryRow(ctx, query, username).
		Scan(&savedUsername)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return "", constant.ErrNotFound
		}
		return "", err
	}

	return savedUsername, nil
}

func (r *UserRepository) FindUserUsername(
	ctx context.Context,
	username string,
) (string, error) {
	query := `
    select username
    from user_usernames
    where username = $1
  `
	var savedUsername string
	err := r.db.QueryRow(ctx, query, username).
		Scan(&savedUsername)
	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return "", constant.ErrNotFound
		}
		return "", err
	}

	return savedUsername, nil
}
