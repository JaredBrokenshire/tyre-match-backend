package builder

import (
	"github.com/jmoiron/sqlx"
	"os"
	"tyre-match-backend/pkg/go-migrations/builder/contract"
	mysql "tyre-match-backend/pkg/go-migrations/builder/mysql/table"
	postgres "tyre-match-backend/pkg/go-migrations/builder/postgress/table"
)

func NewTable(name string, con *sqlx.DB) contract.Table {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return mysql.NewTable(name, con)
	case "postgres":
		return postgres.NewTable(name, con)
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func DropTable(name string, con *sqlx.DB) contract.Table {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return mysql.DropTable(name, con)
	case "postgres":
		return postgres.DropTable(name, con)
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func ChangeTable(name string, con *sqlx.DB) contract.Table {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return mysql.ChangeTable(name, con)
	case "postgres":
		return postgres.ChangeTable(name, con)
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}

func RenameTable(oldName, newName string, con *sqlx.DB) contract.Table {
	switch os.Getenv("DB_DRIVER") {
	case "mysql":
		return mysql.RenameTable(oldName, newName, con)
	case "postgres":
		return postgres.RenameTable(oldName, newName, con)
	default:
		panic("Not supported DB driver: " + os.Getenv("DB_DRIVER"))
	}
}
