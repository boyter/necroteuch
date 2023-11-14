package data

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"necroteuch/common"
)

func ConnectDB(config common.Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite", config.SqliteName)
	if err != nil {
		return nil, err
	}

	return db, nil
}
