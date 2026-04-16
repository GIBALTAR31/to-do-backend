package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

//Connection pooling (too maintain 10 - 20 diffrent connection inside of a pool)
//it will make the api call much faster if API is being called simultaneosly

func Connect(databaseURL string) (*pgxpool.Pool, error){
	var ctx context.Context = context.Background()

	var config *pgxpool.Config
	var err error

	//Parse database URL
	config, err = pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Printf("Unable to parse Database_URL: %v", err)
		return nil, err
	}

	//Create connection pool
	var pool *pgxpool.Pool
	pool, err = pgxpool.NewWithConfig(ctx, config)

	if err != nil{
		log.Printf("Unable to create connection pool: %v", err)
		return nil, err
	}

	// Ping database
	err = pool.Ping(ctx)
	if err != nil{
		log.Printf("Unable to ping database: %v", err)
		pool.Close()
		return nil, err
	}

	log.Println("Succesfully connected to PostgreSQL database")
	return pool, nil
}