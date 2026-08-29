package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type UpdateFilesFileTypeAddBinary struct{}

func (m *UpdateFilesFileTypeAddBinary) GetName() string {
	return "UpdateFilesFileTypeAddBinary"
}

func (m *UpdateFilesFileTypeAddBinary) Up(con *sqlx.DB) {
	table := builder.ChangeTable("files", con)
	table.Column("file_type").Change().Type("ENUM('original','enhanced', 'binary')").Default("original")
	table.MustExec()
}

func (m *UpdateFilesFileTypeAddBinary) Down(con *sqlx.DB) {
	table := builder.ChangeTable("files", con)
	table.Column("file_type").Change().Type("ENUM('original','normalised','enhanced','binary','clean','skeleton')").Default("original")
	table.MustExec()
}
