package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type UpdateTyreModelsAddROIAndPPI struct{}

func (m *UpdateTyreModelsAddROIAndPPI) GetName() string {
	return "UpdateTyreModelsAddROIAndPPI"
}

func (m *UpdateTyreModelsAddROIAndPPI) Up(con *sqlx.DB) {
	table := builder.ChangeTable("tyre_models", con)
	table.Column("pixels_per_inch").Type("FLOAT").NotNull()
	table.Column("roi_top").Type("INT UNSIGNED").Nullable()
	table.Column("roi_left").Type("INT UNSIGNED").Nullable()
	table.Column("roi_right").Type("INT UNSIGNED").Nullable()
	table.Column("roi_bottom").Type("INT UNSIGNED").Nullable()
	table.MustExec()
}

func (m *UpdateTyreModelsAddROIAndPPI) Down(con *sqlx.DB) {
	table := builder.ChangeTable("tyre_models", con)
	table.DropColumn("pixels_per_inch")
	table.DropColumn("roi_top")
	table.DropColumn("roi_left")
	table.DropColumn("roi_right")
	table.DropColumn("roi_bottom")
	table.MustExec()
}
