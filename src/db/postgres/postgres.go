package postgres

import (
	"context"
	"fmt"
	"strings"

	"highload-sn-backend/config"
	"highload-sn-backend/types"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	host     string
	username string
	password string
	db       string

	conn *pgxpool.Pool
}

var pgClient *Client

func InitClient() error {
	host, err := config.Get(config.PostgresHost)
	if err != nil {
		return err
	}
	username, err := config.Get(config.PostgresUsername)
	if err != nil {
		return err
	}
	password, err := config.Get(config.PostgresPassword)
	if err != nil {
		return err
	}
	db, err := config.Get(config.PostgresDB)
	if err != nil {
		return err
	}

	client := Client{
		host:     host,
		username: username,
		password: password,
		db:       db,

		conn: nil,
	}

	err = client.init()
	if err != nil {
		return err
	}

	pgClient = &client

	return nil
}

func (db *Client) init() error {
	err := db.initConn()
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(context.Background(), initQuery)
	return err
}

func (db *Client) initConn() error {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		db.username, db.password, db.host, db.db,
	)
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return err
	}

	db.conn = conn

	return nil
}

func (db *Client) exec(ctx context.Context, query types.QueryItem) error {
	_, err := db.conn.Exec(ctx, query.Query, query.Params...)
	return err
}

func (db *Client) execTx(ctx context.Context, transQueries []types.QueryItem) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, t := range transQueries {
		_, err = tx.Exec(ctx, t.Query, t.Params...)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (db *Client) query(ctx context.Context, query types.QueryItem, scanCb func(pgx.Rows) (any, error)) ([]any, error) {
	rows, err := db.conn.Query(ctx, query.Query, query.Params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var (
		item   any
		result []any
	)
	for rows.Next() {
		item, err = scanCb(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func AddSession(userId, token string) error {
	query := types.QueryItem{
		Query: fmt.Sprintf(
			insertQuery,
			sessionsTableName,
			strings.Join(fieldsOf[sessionsTableName], ","),
			fmt.Sprintf("(%s)", strings.Join(sqlParamsPH(1, 2), ",")),
		),
		Params: []any{userId, token},
	}

	return pgClient.exec(context.Background(), query)
}

func GetSession(userId string) (string, error) {
	query := types.QueryItem{
		Query: fmt.Sprintf(
			"%s %s",
			fmt.Sprintf(
				selectQuery,
				strings.Join(fieldsOf[sessionsTableName], ","),
				sessionsTableName,
			),
			fmt.Sprintf(
				whereQuerySec,
				fmt.Sprintf("user_id = $1")),
		),
		Params: []any{userId},
	}

	results, err := pgClient.query(context.Background(), query, scanSession)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", types.ErrNoUser
	}

	return results[0].(types.Session).Token, nil
}

func AddUsers(users []types.User) error {
	var (
		batchParams   []any
		batchParamsPH []string
	)

	for _, user := range users {
		params := addUserSqlParams(user)
		paramsPH := fmt.Sprintf("(%s)", strings.Join(sqlParamsPH(len(batchParams)+1, len(params)), ","))

		batchParamsPH = append(batchParamsPH, paramsPH)
		batchParams = append(batchParams, params...)
	}

	query := types.QueryItem{
		Query: fmt.Sprintf(
			insertQuery,
			usersTableName,
			strings.Join(fieldsOf[usersTableName], ","),
			strings.Join(batchParamsPH, ","),
		),
		Params: batchParams,
	}

	return pgClient.exec(context.Background(), query)
}

func GetUsers(ids []string) ([]types.User, error) {
	var (
		batchParams   []any
		batchParamsPH = sqlParamsPH(1, len(ids))
	)

	for _, id := range ids {
		batchParams = append(batchParams, id)
	}

	query := types.QueryItem{
		Query: fmt.Sprintf(
			"%s %s",
			fmt.Sprintf(
				selectQuery,
				strings.Join(fieldsOf[usersTableName], ","),
				usersTableName,
			),
			fmt.Sprintf(
				whereQuerySec,
				fmt.Sprintf("id IN (%s)", strings.Join(batchParamsPH, ","))),
		),
		Params: batchParams,
	}

	results, err := pgClient.query(context.Background(), query, scanUser)
	if err != nil {
		return nil, err
	}

	var users []types.User
	for _, result := range results {
		users = append(users, result.(types.User))
	}

	return users, nil
}
