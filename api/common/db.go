package common

import (
	"context"
	"os"

	pg "github.com/go-pg/pg/v10"
)

var sto *pg.DB

func ConnectToDB() (*pg.DB, error) {

	if sto != nil {
		return sto, nil
	}

	opt, err := pg.ParseURL(os.Getenv("DB_CONNSTRING"))
	if err != nil {
		return nil, err
	}

	// required for testing
	opt.TLSConfig = nil

	db := pg.Connect(opt)

	err = db.Ping(context.Background())

	if err != nil {
		return nil, err
	}

	sto = db

	return sto, nil
}
