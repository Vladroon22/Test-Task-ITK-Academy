package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect() (*pgxpool.Pool, error) {
	dbURL := os.Getenv("DB")
	if dbURL == "" {
		log.Fatalln("dbURL doesn't set")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}
	if err := ping(ctx, pool); err != nil {
		return nil, err
	}
	return pool, nil
}

func ping(c context.Context, cn *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(c, time.Second*5)
	defer cancel()
	var err error
	for i := 0; i < 5; i++ {
		if err = cn.Ping(ctx); err == nil {
			return nil
		}
		time.Sleep(time.Millisecond * 500)
	}
	return err
}

/*
type DataBase struct {
	sql *pgx.ConnPool
}

func NewDB() *DataBase {
	return &DataBase{}
}

func (d *DataBase) Connect() (*sql.DB, error) {
	dbURL := os.Getenv("DB")
	if dbURL == "" {
		log.Fatalln("dbURL doesn't set")
	}

	str := fmt.Sprintf("postgresql://%s", dbURL)
	db, err := sql.Open("postgres", str)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err := RetryPing(db); err != nil {
		log.Println(err)
		return nil, err
	}
	d.sql = db

	log.Println("Database connection is valid")
	return db, nil
}

func RetryPing(db *sql.DB) error {
	var err error
	for i := 0; i < 5; i++ {
		if err = db.Ping(); err == nil {
			return nil
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (d *DataBase) CloseDB() error {
	return d.sql.Close()
}
*/
