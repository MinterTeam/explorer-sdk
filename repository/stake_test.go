package repository

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"os"
	"testing"
)

func TestGetStakes(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println(".env file not found")
	}

	//dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
	//	os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
	//	os.Getenv("DB_HOST"), os.Getenv("DB_PORT"),
	//	os.Getenv("DB_NAME"), os.Getenv("DB_SSL_ENABLED"))
	//pgconn := pgdriver.NewConnector(pgdriver.WithDSN(dsn))

	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithUser(os.Getenv("DB_USER")),
		pgdriver.WithPassword(os.Getenv("DB_PASSWORD")),
		pgdriver.WithDatabase(os.Getenv("DB_NAME")),
		pgdriver.WithApplicationName("myapp"),
	)
	sqldb := sql.OpenDB(pgconn)

	r := NewStakeRepository(sqldb, pgdialect.New())

	count, err := r.GetDelegatorsCount()

	if count == 0 {
		t.Error("empty result")
	}
}
