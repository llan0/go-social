package db

import (
	"context"
	"database/sql"
	"time"
)

func New(addr string, maxOpenConn, maxIdelConn int, maxIdelTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConn)
	db.SetMaxIdleConns(maxIdelConn)

	//parse maxIdelConn str -> time.Duration
	// A clearner way to do this is to create a GetDuration() in env.go
	//
	//	duration, err := time.ParseDuration(maxIdelTime)
	//	if err != nil {
	//		return nil, err
	// }

	db.SetConnMaxIdleTime(maxIdelTime) // cleanrer implementaiton useing GetDuration()

	// ping the db to see if its alive
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
