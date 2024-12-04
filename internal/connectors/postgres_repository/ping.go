package postgres_repository

import (
	"context"
	"database/sql"
)

type PingRepo struct {
	connection *sql.DB
}

func NewPingRepo(connection *sql.DB) *PingRepo {
	return &PingRepo{connection: connection}
}

func (r *PingRepo) PingDB(_ context.Context) (string, error) {
	return "ok", nil
}
