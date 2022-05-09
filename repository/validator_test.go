package repository

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"log"
	"os"
	"testing"
)

func TestValidatorFindByPk(t *testing.T) {

	pk := "e782c9a2c62f085f4d1bedf307de525b13226c20c597e66b0cf246a061f31b2d"

	err := godotenv.Load("../.env")
	if err != nil {
		log.Println(".env file not found")
	}

	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(fmt.Sprintf("%s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))),
		pgdriver.WithTLSConfig(nil),
		pgdriver.WithUser(os.Getenv("DB_USER")),
		pgdriver.WithPassword(os.Getenv("DB_PASSWORD")),
		pgdriver.WithDatabase(os.Getenv("DB_NAME")),
		pgdriver.WithApplicationName("Explorer SDK Test"),
	)
	sqldb := sql.OpenDB(pgconn)
	err = sqldb.Ping()
	if err != nil {
		t.Error(err)
	}

	r := NewValidatorRepository(sqldb, pgdialect.New())

	v, err := r.FindByPk(pk)

	if v.PublicKey == "" {
		t.Error("empty result")
	}
}
