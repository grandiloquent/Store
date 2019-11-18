package datastore

import "github.com/jackc/pgx"

type DataStore struct {
	*pgx.ConnPool
}

func (s *DataStore) Fetch(sql string, args ...interface{}) ([]interface{}, error) {
	rows, err := s.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	var items []interface{}
	for rows.Next() {
		row, err := rows.Values()
		if err != nil {
			return nil, err
		}
		items = append(items, row)
	}
	return items, nil

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
