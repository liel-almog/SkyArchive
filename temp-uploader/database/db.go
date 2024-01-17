package database

import (
	"context"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLpgx struct {
	Pool *pgxpool.Pool
}

var (
	db         *PostgreSQLpgx
	initDBOnce sync.Once
)

func newDB() {
	initDBOnce.Do(func() {
		pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			panic(err)
		}

		db = &PostgreSQLpgx{
			Pool: pool,
		}
	})
}

func (p *PostgreSQLpgx) Close() {
	p.Pool.Close()
}

func GetDB() *PostgreSQLpgx {
	if db == nil {
		newDB()
	}

	return db
}
