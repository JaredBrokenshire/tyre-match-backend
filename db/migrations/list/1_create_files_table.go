package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type CreateFilesTable struct{}

func (m *CreateFilesTable) GetName() string {
	return "CreateFilesTable"
}

func (m *CreateFilesTable) Up(con *sqlx.DB) {
	table := builder.NewTable("files", con)
	table.Column("id").Type("int unsigned").NotNull().Autoincrement()
	table.PrimaryKey("id")

	table.String("model", 100).NotNull()
	table.Column("model_id").Type("int unsigned")
	table.Column("file_type").Type("ENUM('original','normalised','enhanced','binary','clean','skeleton')").Default("original")

	table.String("name", 1000)
	table.String("location", 1000)

	table.WithTimestamps()
	table.MustExec()
}

func (m *CreateFilesTable) Down(con *sqlx.DB) {
	builder.DropTable("files", con).MustExec()
}
