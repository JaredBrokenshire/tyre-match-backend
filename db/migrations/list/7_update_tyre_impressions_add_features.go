package list

import (
	"github.com/jmoiron/sqlx"
	"tyre-match-backend/pkg/go-migrations/builder"
)

type UpdateTyreImpressionsAddFeatures struct{}

func (m *UpdateTyreImpressionsAddFeatures) GetName() string {
	return "UpdateTyreImpressionsAddFeatures"
}

func (m *UpdateTyreImpressionsAddFeatures) Up(con *sqlx.DB) {
	table := builder.ChangeTable("tyre_impressions", con)
	table.Column("feature_json").Type("LONGTEXT").Nullable()
	table.Column("matches_json").Type("LONGTEXT").Nullable()
	table.MustExec()
}

func (m *UpdateTyreImpressionsAddFeatures) Down(con *sqlx.DB) {
	table := builder.ChangeTable("tyre_impressions", con)
	table.DropColumn("feature_json")
	table.DropColumn("matches_json")
	table.MustExec()
}
