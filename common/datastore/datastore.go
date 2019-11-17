package datastore

import "github.com/jackc/pgx"

type DataStore struct {
	*pgx.ConnPool
}

func NewDataStore(user, password, database, host string) (*DataStore, error) {
	pgxconfig := pgx.ConnConfig{
		User:     user,
		Password: password,
		Database: database,
	}
	if len(host) > 0 {
		pgxconfig.Host = host
	}

	conn, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxconfig,
		MaxConnections: 2,
	})
	if err != nil {
		return nil, err
	}
	return &DataStore{
		ConnPool: conn}, nil
}
