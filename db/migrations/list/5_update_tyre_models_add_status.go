package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type UpdateTyreModelsAddStatus struct{}

func (m *UpdateTyreModelsAddStatus) GetName() string {
	return "UpdateTyreModelsAddStatus"
}

func (m *UpdateTyreModelsAddStatus) Up(con *sqlx.DB) {
	table := builder.ChangeTable("tyre_models", con)
	table.Column("status").Type("ENUM('Uploaded','Processing','Processed','Matched','Failed')").NotNull().Default("Uploaded")
	table.MustExec()
}

func (m *UpdateTyreModelsAddStatus) Down(con *sqlx.DB) {
	table := builder.ChangeTable("tyre_models", con)
	table.DropColumn("status")
	table.MustExec()
}
