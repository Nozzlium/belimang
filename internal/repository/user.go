package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
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
    usernames (
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
    usernames (
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
      inner join usernames u on u.user_id = a.id 
    where u.username = $1
  `

	err := r.db.QueryRow(
		ctx,
		queryFindUser,
		user.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
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
      inner join usernames un on un.user_id = u.id 
    where un.username = $1
  `

	err := r.db.QueryRow(
		ctx,
		queryFindUser,
		user.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}
