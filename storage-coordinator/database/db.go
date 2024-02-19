package database

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lielalmog/SkyArchive/storage-coordinator/configs"
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
		dbURL, err := configs.GetEnv("DATABASE_URL")
		if err != nil {
			panic(err)
		}

		pool, err := pgxpool.New(context.Background(), dbURL)
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
