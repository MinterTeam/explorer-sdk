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

func TestGetLastBlock(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println(".env file not found")
	}

	var tsl *tls.Config
	if os.Getenv("DB_SSL_ENABLED") == "true" {
		tsl = &tls.Config{InsecureSkipVerify: true}
	}
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))),
		pgdriver.WithTLSConfig(tsl),
		pgdriver.WithUser(os.Getenv("DB_USER")),
		pgdriver.WithPassword(os.Getenv("DB_PASSWORD")),
		pgdriver.WithDatabase(os.Getenv("DB_NAME")),
		pgdriver.WithApplicationName("myapp"),
	)
	sqlDB := sql.OpenDB(pgconn)
	r := NewBlockRepository(sqlDB, pgdialect.New())

	block, err := r.GetLastFromDB()
	if err != nil {
		t.Error(err.Error())
	}
	if block == nil {
		t.Error("empty result")
	}
}
