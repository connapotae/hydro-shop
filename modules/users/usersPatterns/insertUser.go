package usersPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/connapotae/hydro-shop/modules/users"
	"github.com/jmoiron/sqlx"
)

type IInsertUser interface {
	Customer() (IInsertUser, error)
	Admin() (IInsertUser, error)
	Result() (*users.UserPassport, error)
}

type userReq struct {
	id  string
	req *users.UserRegisterReq
	db  *sqlx.DB
}

type customer struct {
	*userReq
}

type admin struct {
	*userReq
}

func InsertUser(db *sqlx.DB, req *users.UserRegisterReq, isAdmin bool) IInsertUser {
	if isAdmin {
		return newAdmin(db, req)
	}
	return newCustomer(db, req)
}

func newCustomer(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser {
	return &customer{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func newAdmin(db *sqlx.DB, req *users.UserRegisterReq) IInsertUser {
	return &admin{
		userReq: &userReq{
			req: req,
			db:  db,
		},
	}
}

func (u *userReq) Customer() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2 ,$3, 1)
	RETURNING "id";`

	if err := u.db.QueryRowContext(
		ctx,
		query,
		u.req.Email,
		u.req.Password,
		u.req.Username,
	).Scan(&u.id); err != nil {
		return nil, fmt.Errorf("inser user failed: %v", err)
	}

	return u, nil
}

func (u *userReq) Admin() (IInsertUser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
	INSERT INTO "users" (
		"email",
		"password",
		"username",
		"role_id"
	)
	VALUES
		($1, $2 ,$3, 2)
	RETURNING "id";`

	if err := u.db.QueryRowContext(
		ctx,
		query,
		u.req.Email,
		u.req.Password,
		u.req.Username,
	).Scan(&u.id); err != nil {
		return nil, fmt.Errorf("inser user failed: %v", err)
	}

	return u, nil
}

func (u *userReq) Result() (*users.UserPassport, error) {
	query := `
	SELECT
		json_build_object(
			'user', "t",
			'token', NULL
		)
	FROM (
		SELECT
			"u"."id",
			"u"."email",
			"u"."username",
			"u"."role_id"
		FROM "users" "u"
		WHERE "u"."id" = $1
	) as "t"`

	data := make([]byte, 0)
	if err := u.db.Get(&data, query, u.id); err != nil {
		return nil, fmt.Errorf("get user failed: %v", err)
	}

	user := new(users.UserPassport)
	if err := json.Unmarshal(data, &user); err != nil {
		return nil, fmt.Errorf("unmarshal user failed: %v", err)
	}

	return user, nil
}
