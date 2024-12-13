package database

import (
	"fmt"
	"os"
	"strings"

	"quiz_platform/internal/misc/config"

	"database/sql"

	_ "github.com/lib/pq"
)

// Initializes connected postgres database by the parameters in the config.
func initializePostgresDatabase(db *sql.DB) {
	script, err := os.ReadFile(config.GlobalConfig.Database.InitScriptPath)
	if err != nil {
		panic(fmt.Errorf("fatal error reading initialize postgres sql script: %v", err))
	}

	statements := strings.Split(string(script), ";")

	for _, statement := range statements {
		_, err := db.Exec(statement)
		if err != nil {
			panic(fmt.Errorf("fatal error executing initialize postgres sql script statement: %v", err))
		}
	}
}
