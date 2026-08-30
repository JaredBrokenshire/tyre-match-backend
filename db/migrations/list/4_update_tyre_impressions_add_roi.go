package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type UpdateTyreImpressionsAddROI struct{}

func (m *UpdateTyreImpressionsAddROI) GetName() string {
	return "UpdateTyreImpressionsAddROI"
}

func (m *UpdateTyreImpressionsAddROI) Up(con *sqlx.DB) {
	table := builder.ChangeTable("tyre_impressions", con)
	table.Column("roi_top").Type("INT UNSIGNED").Nullable()
	table.Column("roi_left").Type("INT UNSIGNED").Nullable()
	table.Column("roi_right").Type("INT UNSIGNED").Nullable()
	table.Column("roi_bottom").Type("INT UNSIGNED").Nullable()
	table.MustExec()
}

func (m *UpdateTyreImpressionsAddROI) Down(con *sqlx.DB) {
	table := builder.ChangeTable("tyre_impressions", con)
	table.DropColumn("roi_top")
	table.DropColumn("roi_left")
	table.DropColumn("roi_right")
	table.DropColumn("roi_bottom")
	table.MustExec()
}
