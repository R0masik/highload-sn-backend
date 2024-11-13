package postgres

import (
	"fmt"

	"highload-sn-backend/types"

	"github.com/jackc/pgx/v5"
)

const (
	usersTableName    = "users"
	sessionsTableName = "sessions"
)

var fieldsOf = map[string][]string{
	usersTableName: {
		"id",
		"first_name",
		"last_name",
		"birth_date",
		"sex",
		"biography",
		"city",
		"password_hash",
	},
	sessionsTableName: {
		"user_id",
		"token",
	},
}

func sqlParamsPH(start, count int) []string {
	var params []string
	for i := start; i < start+count; i++ {
		params = append(params, fmt.Sprintf("$%d", i))
	}

	return params
}

func addUserSqlParams(item types.User) []any {
	return []any{
		item.Id,
		item.FirstName,
		item.LastName,
		item.BirthDate,
		item.Sex,
		item.Biography,
		item.City,
		item.PasswordHash,
	}
}

func scanUser(rows pgx.Rows) (any, error) {
	var item types.User

	err := rows.Scan(
		&item.Id,
		&item.FirstName,
		&item.LastName,
		&item.BirthDate,
		&item.Sex,
		&item.Biography,
		&item.City,
		&item.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func scanSession(rows pgx.Rows) (any, error) {
	var item types.Session

	err := rows.Scan(
		&item.UserId,
		&item.Token,
	)
	if err != nil {
		return nil, err
	}

	return item, nil
}
