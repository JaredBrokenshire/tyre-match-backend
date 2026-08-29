package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type UpdateFilesFileTypeAddNormalised struct{}

func (m *UpdateFilesFileTypeAddNormalised) GetName() string {
	return "UpdateFilesFileTypeAddNormalised"
}

func (m *UpdateFilesFileTypeAddNormalised) Up(con *sqlx.DB) {
	table := builder.ChangeTable("files", con)
	table.Column("file_type").Change().Type("ENUM('original', 'normalised','enhanced', 'binary')").Default("original")
	table.MustExec()
}

func (m *UpdateFilesFileTypeAddNormalised) Down(con *sqlx.DB) {
	table := builder.ChangeTable("files", con)
	table.Column("file_type").Change().Type("ENUM('original','enhanced','binary')").Default("original")
	table.MustExec()
}
