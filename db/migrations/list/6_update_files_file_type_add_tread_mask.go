package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type UpdateFilesFileTypeAddTreadMask struct{}

func (m *UpdateFilesFileTypeAddTreadMask) GetName() string {
	return "UpdateFilesFileTypeAddTreadMask"
}

func (m *UpdateFilesFileTypeAddTreadMask) Up(con *sqlx.DB) {
	table := builder.ChangeTable("files", con)
	table.Column("file_type").Change().Type("ENUM('original','normalised','enhanced','binary','tread_mask')").Default("original")
	table.MustExec()
}

func (m *UpdateFilesFileTypeAddTreadMask) Down(con *sqlx.DB) {
	table := builder.ChangeTable("files", con)
	table.Column("file_type").Change().Type("ENUM('original','normalised','enhanced','binary')").Default("original")
	table.MustExec()
}
